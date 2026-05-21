package plugin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CatalogItem represents a plugin available in the remote official catalog.
type CatalogItem struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	URL         string `json:"url"`
	ComposeURL  string `json:"composeUrl,omitempty"`
}

// CatalogItemStatus wraps CatalogItem with its local installation state.
type CatalogItemStatus struct {
	CatalogItem
	Installed bool   `json:"installed"`
	Status    string `json:"status"` // "enabled", "disabled", "not_installed"
}

// Catalog is the static index of official M3TAL plugins available for download.
var Catalog = []CatalogItem{
	{
		Name:        "m3tal",
		Kind:        "Stack",
		Description: "M3TAL Core dashboard and system services stack",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/stacks/m3tal.yml",
		ComposeURL:  "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/m3tal-compose.yml",
	},
	{
		Name:        "maintenance",
		Kind:        "Stack",
		Description: "System maintenance and monitoring tools",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/stacks/maintenance.yml",
		ComposeURL:  "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/maintenance-compose.yml",
	},
	{
		Name:        "network",
		Kind:        "Stack",
		Description: "Core network services and reverse proxy setup",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/stacks/network.yml",
		ComposeURL:  "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/network-compose.yml",
	},
	{
		Name:        "routing",
		Kind:        "Stack",
		Description: "Traefik reverse proxy and SSL termination",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/stacks/routing.yml",
		ComposeURL:  "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/routing-compose.yml",
	},
	{
		Name:        "basicauth",
		Kind:        "Middleware",
		Description: "Basic authentication preset for routes",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/middleware/basicauth.yml",
	},
	{
		Name:        "ratelimit",
		Kind:        "Middleware",
		Description: "Rate limiting configuration for endpoints",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/middleware/ratelimit.yml",
	},
	{
		Name:        "security-headers",
		Kind:        "Middleware",
		Description: "Hardened security headers for HTTP services",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/middleware/security-headers.yml",
	},
	{
		Name:        "example-route",
		Kind:        "Route",
		Description: "Example route plugin — customize for your service",
		Version:     "1.0.0",
		Author:      "M3TAL Team",
		URL:         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/plugins/routes/example-route.yml",
	},
}

// GetPluginBaseName extracts the clean name of the plugin from its filepath.
func GetPluginBaseName(path string) string {
	if path == "" {
		return ""
	}
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, ".disabled")
	base = strings.TrimSuffix(base, ".yml")
	base = strings.TrimSuffix(base, ".yaml")
	return base
}

// MatchesPluginName checks if a plugin matches the target catalog name.
func MatchesPluginName(sourcePath, metadataName, targetName string) bool {
	if strings.EqualFold(metadataName, targetName) {
		return true
	}
	base := GetPluginBaseName(sourcePath)
	return base != "" && strings.EqualFold(base, targetName)
}

// ListCatalog correlates remote catalog items with local loaded plugins.
func ListCatalog(reg *Registry) []CatalogItemStatus {
	var list []CatalogItemStatus

	for _, item := range Catalog {
		status := CatalogItemStatus{
			CatalogItem: item,
			Installed:   false,
			Status:      "not_installed",
		}

		switch item.Kind {
		case KindRoute:
			for _, r := range reg.Routes {
				if MatchesPluginName(r.SourcePath, r.Metadata.Name, item.Name) {
					status.Installed = true
					if r.Enabled {
						status.Status = "enabled"
					} else {
						status.Status = "disabled"
					}
					break
				}
			}
		case KindStack:
			for _, s := range reg.Stacks {
				if MatchesPluginName(s.SourcePath, s.Metadata.Name, item.Name) {
					status.Installed = true
					if s.Enabled {
						status.Status = "enabled"
					} else {
						status.Status = "disabled"
					}
					break
				}
			}
		case KindMiddleware:
			for _, m := range reg.Middlewares {
				if MatchesPluginName(m.SourcePath, m.Metadata.Name, item.Name) {
					status.Installed = true
					if m.Enabled {
						status.Status = "enabled"
					} else {
						status.Status = "disabled"
					}
					break
				}
			}
		}

		list = append(list, status)
	}

	return list
}

// downloadFile fetches a URL and writes it locally.
func downloadFile(url, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(dest), err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}

// InstallPlugin downloads a plugin's definition and optional compose files.
func InstallPlugin(name, kind, userPluginsDir string) error {
	var targetItem *CatalogItem
	for i := range Catalog {
		if strings.EqualFold(Catalog[i].Name, name) && strings.EqualFold(Catalog[i].Kind, kind) {
			targetItem = &Catalog[i]
			break
		}
	}

	if targetItem == nil {
		return fmt.Errorf("plugin %q (kind: %s) not found in catalog", name, kind)
	}

	var subfolder string
	switch targetItem.Kind {
	case KindRoute:
		subfolder = "routes"
	case KindStack:
		subfolder = "stacks"
	case KindMiddleware:
		subfolder = "middleware"
	default:
		return fmt.Errorf("unsupported plugin kind: %s", targetItem.Kind)
	}

	pluginDest := filepath.Join(userPluginsDir, subfolder, targetItem.Name+".yml")
	err := downloadFile(targetItem.URL, pluginDest)
	if err != nil {
		return fmt.Errorf("failed to download plugin manifest: %w", err)
	}

	// Download compose file if it is a Stack plugin and has a ComposeURL
	if targetItem.Kind == KindStack && targetItem.ComposeURL != "" {
		composeFilename := filepath.Base(targetItem.ComposeURL)
		composeDest := filepath.Join(userPluginsDir, subfolder, composeFilename)
		err = downloadFile(targetItem.ComposeURL, composeDest)
		if err != nil {
			_ = os.Remove(pluginDest) // clean up manifest
			return fmt.Errorf("failed to download associated compose file %s: %w", composeFilename, err)
		}
	}

	return nil
}

// UninstallPlugin deletes a user-installed plugin's files.
func UninstallPlugin(name, kind, userPluginsDir string, reg *Registry) error {
	var sourcePath string
	var composePathRel string

	switch kind {
	case KindRoute:
		for i := range reg.Routes {
			if MatchesPluginName(reg.Routes[i].SourcePath, reg.Routes[i].Metadata.Name, name) {
				sourcePath = reg.Routes[i].SourcePath
				break
			}
		}
	case KindStack:
		for i := range reg.Stacks {
			if MatchesPluginName(reg.Stacks[i].SourcePath, reg.Stacks[i].Metadata.Name, name) {
				sourcePath = reg.Stacks[i].SourcePath
				composePathRel = reg.Stacks[i].ComposePath
				break
			}
		}
	case KindMiddleware:
		for i := range reg.Middlewares {
			if MatchesPluginName(reg.Middlewares[i].SourcePath, reg.Middlewares[i].Metadata.Name, name) {
				sourcePath = reg.Middlewares[i].SourcePath
				break
			}
		}
	}

	if sourcePath == "" {
		return fmt.Errorf("plugin %q (kind: %s) is not currently installed", name, kind)
	}

	// Security check: Only allow deleting files inside userPluginsDir
	absUserDir, err := filepath.Abs(userPluginsDir)
	if err != nil {
		return fmt.Errorf("failed to resolve user plugins directory path: %w", err)
	}

	absSourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to resolve plugin source path: %w", err)
	}

	// Allow development/local paths if we are running in local/dev directory
	// Normally SystemPluginsDir is "/opt/m3tal/plugins", and userPluginsDir is "/etc/m3tal/plugins"
	// But in development, they might use "deploy/plugins"
	isDeveloperPath := false
	if _, err := os.Stat("deploy/plugins"); err == nil {
		absDevDir, err := filepath.Abs("deploy/plugins")
		if err == nil && strings.HasPrefix(absSourcePath, absDevDir) {
			isDeveloperPath = true
		}
	}

	if !strings.HasPrefix(absSourcePath, absUserDir) && !isDeveloperPath {
		return fmt.Errorf("security violation: cannot uninstall system plugin at %s", sourcePath)
	}

	// Delete main manifest
	err = os.Remove(absSourcePath)
	if err != nil {
		return fmt.Errorf("failed to remove plugin manifest: %w", err)
	}

	// If Stack plugin, also clean up the compose file
	if kind == KindStack && composePathRel != "" && !filepath.IsAbs(composePathRel) {
		composeFile := filepath.Join(filepath.Dir(absSourcePath), composePathRel)
		if _, err := os.Stat(composeFile); err == nil {
			_ = os.Remove(composeFile)
		}
	}

	return nil
}
