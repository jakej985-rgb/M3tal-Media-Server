package plugin

import (
	"fmt"
	"sort"
)

// Registry holds all loaded and typed plugins, indexed by kind.
type Registry struct {
	Routes      []RoutePlugin
	Stacks      []StackPlugin
	Middlewares []MiddlewarePlugin
}

// BuildRegistry takes a list of validated plugins and constructs a typed Registry.
func BuildRegistry(plugins []Plugin) (*Registry, error) {
	r := &Registry{}

	for i := range plugins {
		p := &plugins[i]

		switch p.Kind {
		case KindRoute:
			rp, err := p.AsRoute()
			if err != nil {
				return nil, fmt.Errorf("plugin %q: %w", p.Metadata.Name, err)
			}
			r.Routes = append(r.Routes, *rp)

		case KindStack:
			sp, err := p.AsStack()
			if err != nil {
				return nil, fmt.Errorf("plugin %q: %w", p.Metadata.Name, err)
			}
			r.Stacks = append(r.Stacks, *sp)

		case KindMiddleware:
			mp, err := p.AsMiddleware()
			if err != nil {
				return nil, fmt.Errorf("plugin %q: %w", p.Metadata.Name, err)
			}
			r.Middlewares = append(r.Middlewares, *mp)
		}
	}

	// Sort stacks by priority (lower = earlier)
	sort.Slice(r.Stacks, func(i, j int) bool {
		return r.Stacks[i].Priority < r.Stacks[j].Priority
	})

	// Sort routes and middlewares by name for deterministic output
	sort.Slice(r.Routes, func(i, j int) bool {
		return r.Routes[i].Metadata.Name < r.Routes[j].Metadata.Name
	})
	sort.Slice(r.Middlewares, func(i, j int) bool {
		return r.Middlewares[i].Metadata.Name < r.Middlewares[j].Metadata.Name
	})

	return r, nil
}

// GetRoute returns a route plugin by metadata name, or nil if not found.
func (r *Registry) GetRoute(name string) *RoutePlugin {
	for i := range r.Routes {
		if r.Routes[i].Metadata.Name == name {
			return &r.Routes[i]
		}
	}
	return nil
}

// GetStack returns a stack plugin by metadata name, or nil if not found.
func (r *Registry) GetStack(name string) *StackPlugin {
	for i := range r.Stacks {
		if r.Stacks[i].Metadata.Name == name {
			return &r.Stacks[i]
		}
	}
	return nil
}

// GetMiddleware returns a middleware plugin by metadata name, or nil if not found.
func (r *Registry) GetMiddleware(name string) *MiddlewarePlugin {
	for i := range r.Middlewares {
		if r.Middlewares[i].Metadata.Name == name {
			return &r.Middlewares[i]
		}
	}
	return nil
}

// ListRoutes returns all loaded route plugins.
func (r *Registry) ListRoutes() []RoutePlugin {
	if r.Routes == nil {
		return []RoutePlugin{}
	}
	return r.Routes
}

// ListStacks returns all loaded stack plugins, sorted by priority.
func (r *Registry) ListStacks() []StackPlugin {
	if r.Stacks == nil {
		return []StackPlugin{}
	}
	return r.Stacks
}

// ListMiddlewares returns all loaded middleware plugins.
func (r *Registry) ListMiddlewares() []MiddlewarePlugin {
	if r.Middlewares == nil {
		return []MiddlewarePlugin{}
	}
	return r.Middlewares
}

// Count returns the total number of plugins in the registry.
func (r *Registry) Count() int {
	return len(r.Routes) + len(r.Stacks) + len(r.Middlewares)
}

// Summary returns a human-readable summary of loaded plugins.
func (r *Registry) Summary() string {
	return fmt.Sprintf("%d plugins loaded (%d routes, %d stacks, %d middleware)",
		r.Count(), len(r.Routes), len(r.Stacks), len(r.Middlewares))
}
