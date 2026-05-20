package api

import (
	"net/http"

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
