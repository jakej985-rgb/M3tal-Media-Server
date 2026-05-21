package plugin

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Supported API versions and kinds.
const (
	APIVersion = "m3tal/v1"

	KindRoute      = "Route"
	KindStack      = "Stack"
	KindMiddleware = "Middleware"
	KindTraefik    = "Traefik"
)

// Plugin represents a loaded plugin manifest.
type Plugin struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   PluginMetadata `yaml:"metadata"`
	Spec       map[string]any `yaml:"spec"`

	// SourcePath is set by the loader to track where this plugin was loaded from.
	SourcePath string `yaml:"-"`
	// Enabled indicates whether the plugin is active (no .disabled suffix).
	Enabled bool `yaml:"-"`
}

// PluginMetadata contains descriptive information about a plugin.
type PluginMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version,omitempty"`
	Author      string `yaml:"author,omitempty"`
}

// RoutePlugin is the typed representation of a Route plugin spec.
type RoutePlugin struct {
	Metadata    PluginMetadata
	Service     string   `yaml:"service"`
	Domain      string   `yaml:"domain"`
	Port        int      `yaml:"port"`
	Entrypoints string   `yaml:"entrypoints,omitempty"`
	Network     string   `yaml:"network,omitempty"`
	Middlewares []string `yaml:"middlewares,omitempty"`
	SourcePath  string   `yaml:"-"`
	Enabled     bool     `yaml:"-"`
}

// StackPlugin is the typed representation of a Stack plugin spec.
type StackPlugin struct {
	Metadata    PluginMetadata
	ComposePath string `yaml:"composePath"`
	EnvTemplate string `yaml:"envTemplate,omitempty"`
	Priority    int    `yaml:"priority,omitempty"`
	Category    string `yaml:"category,omitempty"`
	SourcePath  string `yaml:"-"`
	Enabled     bool   `yaml:"-"`
}

// MiddlewarePlugin is the typed representation of a Middleware plugin spec.
type MiddlewarePlugin struct {
	Metadata   PluginMetadata
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`
	Config     map[string]string `yaml:"config,omitempty"`
	SourcePath string            `yaml:"-"`
	Enabled    bool              `yaml:"-"`
}

// TraefikPlugin is the typed representation of a Traefik plugin spec.
type TraefikPlugin struct {
	Metadata    PluginMetadata
	Routes      []RouteSpec      `yaml:"routes,omitempty"`
	Middlewares []MiddlewareSpec `yaml:"middlewares,omitempty"`
	SourcePath  string           `yaml:"-"`
	Enabled     bool             `yaml:"-"`
}

// RouteSpec represents a route configuration inside a Traefik plugin.
type RouteSpec struct {
	Name        string   `yaml:"name"`
	Service     string   `yaml:"service"`
	Domain      string   `yaml:"domain"`
	Port        int      `yaml:"port"`
	Entrypoints string   `yaml:"entrypoints,omitempty"`
	Network     string   `yaml:"network,omitempty"`
	Middlewares []string `yaml:"middlewares,omitempty"`
}

// MiddlewareSpec represents a middleware configuration inside a Traefik plugin.
type MiddlewareSpec struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config,omitempty"`
}

// Validate checks a raw Plugin for structural correctness.
func (p *Plugin) Validate() error {
	if p.APIVersion != APIVersion {
		return fmt.Errorf("unsupported apiVersion %q (expected %q)", p.APIVersion, APIVersion)
	}

	switch p.Kind {
	case KindRoute, KindStack, KindMiddleware, KindTraefik:
		// valid
	default:
		return fmt.Errorf("unsupported kind %q (expected one of: %s)",
			p.Kind, strings.Join([]string{KindRoute, KindStack, KindMiddleware, KindTraefik}, ", "))
	}

	if p.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}

	return p.validateSpec()
}

// validateSpec performs kind-specific spec checks.
func (p *Plugin) validateSpec() error {
	if p.Spec == nil {
		return fmt.Errorf("spec is required")
	}

	switch p.Kind {
	case KindRoute:
		for _, required := range []string{"service", "domain", "port"} {
			if _, ok := p.Spec[required]; !ok {
				return fmt.Errorf("Route plugin %q: spec.%s is required", p.Metadata.Name, required)
			}
		}
	case KindStack:
		if _, ok := p.Spec["composePath"]; !ok {
			return fmt.Errorf("Stack plugin %q: spec.composePath is required", p.Metadata.Name)
		}
	case KindMiddleware:
		for _, required := range []string{"name", "type"} {
			if _, ok := p.Spec[required]; !ok {
				return fmt.Errorf("Middleware plugin %q: spec.%s is required", p.Metadata.Name, required)
			}
		}
	case KindTraefik:
		if routes, ok := p.Spec["routes"]; ok {
			routesList, ok := routes.([]any)
			if !ok {
				return fmt.Errorf("Traefik plugin %q: spec.routes must be a list", p.Metadata.Name)
			}
			for i, r := range routesList {
				rMap, ok := r.(map[string]any)
				if !ok {
					// Also support map[any]any in case yaml unmarshals it that way
					if rMapAny, ok := r.(map[any]any); ok {
						rMap = make(map[string]any)
						for k, v := range rMapAny {
							rMap[fmt.Sprintf("%v", k)] = v
						}
					} else {
						return fmt.Errorf("Traefik plugin %q: spec.routes[%d] must be a map", p.Metadata.Name, i)
					}
				}
				for _, required := range []string{"service", "domain", "port"} {
					if _, ok := rMap[required]; !ok {
						return fmt.Errorf("Traefik plugin %q routes[%d]: %s is required", p.Metadata.Name, i, required)
					}
				}
			}
		}
		if middlewares, ok := p.Spec["middlewares"]; ok {
			middlewaresList, ok := middlewares.([]any)
			if !ok {
				return fmt.Errorf("Traefik plugin %q: spec.middlewares must be a list", p.Metadata.Name)
			}
			for i, m := range middlewaresList {
				mMap, ok := m.(map[string]any)
				if !ok {
					if mMapAny, ok := m.(map[any]any); ok {
						mMap = make(map[string]any)
						for k, v := range mMapAny {
							mMap[fmt.Sprintf("%v", k)] = v
						}
					} else {
						return fmt.Errorf("Traefik plugin %q: spec.middlewares[%d] must be a map", p.Metadata.Name, i)
					}
				}
				for _, required := range []string{"name", "type"} {
					if _, ok := mMap[required]; !ok {
						return fmt.Errorf("Traefik plugin %q middlewares[%d]: %s is required", p.Metadata.Name, i, required)
					}
				}
			}
		}
	}

	return nil
}

// AsRoute converts a validated Route plugin to its typed struct.
func (p *Plugin) AsRoute() (*RoutePlugin, error) {
	if p.Kind != KindRoute {
		return nil, fmt.Errorf("cannot convert %s plugin to RoutePlugin", p.Kind)
	}

	// Re-marshal spec to YAML then decode into typed struct.
	data, err := yaml.Marshal(p.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec: %w", err)
	}

	var rp RoutePlugin
	if err := yaml.Unmarshal(data, &rp); err != nil {
		return nil, fmt.Errorf("failed to decode Route spec: %w", err)
	}
	rp.Metadata = p.Metadata
	rp.SourcePath = p.SourcePath
	rp.Enabled = p.Enabled
	return &rp, nil
}

// AsStack converts a validated Stack plugin to its typed struct.
func (p *Plugin) AsStack() (*StackPlugin, error) {
	if p.Kind != KindStack {
		return nil, fmt.Errorf("cannot convert %s plugin to StackPlugin", p.Kind)
	}

	data, err := yaml.Marshal(p.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec: %w", err)
	}

	var sp StackPlugin
	if err := yaml.Unmarshal(data, &sp); err != nil {
		return nil, fmt.Errorf("failed to decode Stack spec: %w", err)
	}
	sp.Metadata = p.Metadata
	sp.SourcePath = p.SourcePath
	sp.Enabled = p.Enabled
	return &sp, nil
}

// AsMiddleware converts a validated Middleware plugin to its typed struct.
func (p *Plugin) AsMiddleware() (*MiddlewarePlugin, error) {
	if p.Kind != KindMiddleware {
		return nil, fmt.Errorf("cannot convert %s plugin to MiddlewarePlugin", p.Kind)
	}

	data, err := yaml.Marshal(p.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec: %w", err)
	}

	var mp MiddlewarePlugin
	if err := yaml.Unmarshal(data, &mp); err != nil {
		return nil, fmt.Errorf("failed to decode Middleware spec: %w", err)
	}
	mp.Metadata = p.Metadata
	mp.SourcePath = p.SourcePath
	mp.Enabled = p.Enabled
	return &mp, nil
}

// AsTraefik converts a validated Traefik plugin to its typed struct.
func (p *Plugin) AsTraefik() (*TraefikPlugin, error) {
	if p.Kind != KindTraefik {
		return nil, fmt.Errorf("cannot convert %s plugin to TraefikPlugin", p.Kind)
	}

	data, err := yaml.Marshal(p.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec: %w", err)
	}

	var tp TraefikPlugin
	if err := yaml.Unmarshal(data, &tp); err != nil {
		return nil, fmt.Errorf("failed to decode Traefik spec: %w", err)
	}
	tp.Metadata = p.Metadata
	tp.SourcePath = p.SourcePath
	tp.Enabled = p.Enabled
	return &tp, nil
}

// ParsePlugin parses raw YAML bytes into a Plugin.
func ParsePlugin(data []byte) (*Plugin, error) {
	var p Plugin
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("invalid plugin YAML: %w", err)
	}
	return &p, nil
}
