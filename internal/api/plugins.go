package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/jakej985-rgb/m3tal-core/internal/plugin"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// PluginHandlers provides endpoints for viewing loaded plugins.
type PluginHandlers struct {
	registry *plugin.Registry
}

// NewPluginHandlers creates a handler set with a lazily-loaded plugin registry.
func NewPluginHandlers() *PluginHandlers {
	return &PluginHandlers{}
}

// ensureLoaded lazily loads the plugin registry on first access.
func (h *PluginHandlers) ensureLoaded() *plugin.Registry {
	if h.registry != nil {
		return h.registry
	}

	dirs := system.GetPluginDirs()
	reg, err := plugin.LoadAll(dirs...)
	if err != nil {
		// Return empty registry on error
		h.registry = &plugin.Registry{}
		return h.registry
	}
	h.registry = reg
	return h.registry
}

// ListPlugins returns all loaded plugins across all kinds.
// GET /api/v2/plugins
func (h *PluginHandlers) ListPlugins(w http.ResponseWriter, r *http.Request) {
	reg := h.ensureLoaded()

	writeJSON(w, http.StatusOK, map[string]any{
		"summary":    reg.Summary(),
		"routes":     reg.ListRoutes(),
		"stacks":     reg.ListStacks(),
		"middleware": reg.ListMiddlewares(),
	})
}

// ListPluginsByKind returns plugins filtered by kind.
// GET /api/v2/plugins/{kind}
func (h *PluginHandlers) ListPluginsByKind(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")
	reg := h.ensureLoaded()

	switch kind {
	case "routes":
		writeJSON(w, http.StatusOK, reg.ListRoutes())
	case "stacks":
		writeJSON(w, http.StatusOK, reg.ListStacks())
	case "middleware":
		writeJSON(w, http.StatusOK, reg.ListMiddlewares())
	default:
		writeError(w, http.StatusBadRequest, "invalid plugin kind: "+kind+" (expected: routes, stacks, middleware)")
	}
}

// Reload forces a reload of the plugin registry from disk.
// POST /api/v2/plugins/reload
func (h *PluginHandlers) Reload(w http.ResponseWriter, r *http.Request) {
	h.registry = nil // Force reload on next access
	reg := h.ensureLoaded()
	writeJSON(w, http.StatusOK, map[string]any{
		"reloaded": true,
		"summary":  reg.Summary(),
	})
}

// Enable renames a plugin file ending in `.disabled` by removing the suffix.
// POST /api/v2/plugins/enable
func (h *PluginHandlers) Enable(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	reg := h.ensureLoaded()
	var path string
	switch req.Kind {
	case "Route":
		p := reg.GetRoute(req.Name)
		if p != nil {
			path = p.SourcePath
		}
	case "Stack":
		p := reg.GetStack(req.Name)
		if p != nil {
			path = p.SourcePath
		}
	case "Middleware":
		p := reg.GetMiddleware(req.Name)
		if p != nil {
			path = p.SourcePath
		}
	default:
		if p := reg.GetRoute(req.Name); p != nil {
			path = p.SourcePath
		} else if p := reg.GetStack(req.Name); p != nil {
			path = p.SourcePath
		} else if p := reg.GetMiddleware(req.Name); p != nil {
			path = p.SourcePath
		}
	}

	if path == "" {
		writeError(w, http.StatusNotFound, "plugin not found: "+req.Name)
		return
	}

	newPath, err := plugin.EnablePlugin(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to enable: "+err.Error())
		return
	}

	h.registry = nil // force reload
	writeJSON(w, http.StatusOK, map[string]any{"enabled": true, "path": newPath})
}

// Disable renames a plugin file to append `.disabled`.
// POST /api/v2/plugins/disable
func (h *PluginHandlers) Disable(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	reg := h.ensureLoaded()
	var path string
	switch req.Kind {
	case "Route":
		p := reg.GetRoute(req.Name)
		if p != nil {
			path = p.SourcePath
		}
	case "Stack":
		p := reg.GetStack(req.Name)
		if p != nil {
			path = p.SourcePath
		}
	case "Middleware":
		p := reg.GetMiddleware(req.Name)
		if p != nil {
			path = p.SourcePath
		}
	default:
		if p := reg.GetRoute(req.Name); p != nil {
			path = p.SourcePath
		} else if p := reg.GetStack(req.Name); p != nil {
			path = p.SourcePath
		} else if p := reg.GetMiddleware(req.Name); p != nil {
			path = p.SourcePath
		}
	}

	if path == "" {
		writeError(w, http.StatusNotFound, "plugin not found: "+req.Name)
		return
	}

	newPath, err := plugin.DisablePlugin(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to disable: "+err.Error())
		return
	}

	h.registry = nil // force reload
	writeJSON(w, http.StatusOK, map[string]any{"disabled": true, "path": newPath})
}

// Sync writes the dynamic Traefik configuration file
// POST /api/v2/plugins/sync
func (h *PluginHandlers) Sync(w http.ResponseWriter, r *http.Request) {
	reg := h.ensureLoaded()
	configData, err := reg.GenerateTraefikConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate Traefik config: "+err.Error())
		return
	}

	stackDir := system.GetStackDir()
	dynamicDir := filepath.Join(stackDir, "dynamic")
	if err := os.MkdirAll(dynamicDir, 0755); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create dynamic directory: "+err.Error())
		return
	}

	outputPath := filepath.Join(dynamicDir, "m3tal-plugins.yml")
	if err := os.WriteFile(outputPath, configData, 0644); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to write config file: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"synced": true,
		"path":   outputPath,
	})
}

// ListCatalog returns all official plugins and their local installation state.
// GET /api/v2/plugins/catalog
func (h *PluginHandlers) ListCatalog(w http.ResponseWriter, r *http.Request) {
	reg := h.ensureLoaded()
	catalogStatus := plugin.ListCatalog(reg)
	writeJSON(w, http.StatusOK, catalogStatus)
}

// Install downloads and installs a plugin from the remote catalog.
// POST /api/v2/plugins/install
func (h *PluginHandlers) Install(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Name == "" || req.Kind == "" {
		writeError(w, http.StatusBadRequest, "missing name or kind")
		return
	}

	userDir := system.UserPluginsDir
	if _, err := os.Stat("deploy/plugins"); err == nil {
		userDir = "deploy/plugins"
	}

	err := plugin.InstallPlugin(req.Name, req.Kind, userDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to install plugin: "+err.Error())
		return
	}

	h.registry = nil // force reload
	reg := h.ensureLoaded()

	writeJSON(w, http.StatusOK, map[string]any{
		"installed": true,
		"summary":   reg.Summary(),
	})
}

// Uninstall deletes a user-installed plugin.
// POST /api/v2/plugins/uninstall
func (h *PluginHandlers) Uninstall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Name == "" || req.Kind == "" {
		writeError(w, http.StatusBadRequest, "missing name or kind")
		return
	}

	reg := h.ensureLoaded()

	userDir := system.UserPluginsDir
	if _, err := os.Stat("deploy/plugins"); err == nil {
		userDir = "deploy/plugins"
	}

	err := plugin.UninstallPlugin(req.Name, req.Kind, userDir, reg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to uninstall plugin: "+err.Error())
		return
	}

	h.registry = nil // force reload
	newReg := h.ensureLoaded()

	writeJSON(w, http.StatusOK, map[string]any{
		"uninstalled": true,
		"summary":     newReg.Summary(),
	})
}

