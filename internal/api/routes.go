package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jakej985-rgb/m3tal-core/internal/engine"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
)

// RouteHandlers provides CRUD operations for Traefik routes.
type RouteHandlers struct {
	Store *store.Store
}

// ListRoutes returns all stored routes.
// GET /api/v2/routes
func (h *RouteHandlers) ListRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.Store.ListRoutes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if routes == nil {
		routes = []store.RouteRecord{}
	}
	writeJSON(w, http.StatusOK, routes)
}

// CreateRoute validates and creates a new route.
// POST /api/v2/routes
func (h *RouteHandlers) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var input engine.RouteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Validate route input
	errs := engine.ValidateRoute(input)
	if len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":      "validation failed",
			"violations": errs,
		})
		return
	}

	// Check domain uniqueness
	existing, _ := h.Store.ListRoutes()
	var existingInputs []engine.RouteInput
	for _, r := range existing {
		existingInputs = append(existingInputs, engine.RouteInput{
			Service: r.Service,
			Domain:  r.Domain,
			Port:    r.Port,
		})
	}
	uniqueErrs := engine.ValidateDomainUnique(input.Domain, existingInputs)
	if len(uniqueErrs) > 0 {
		writeJSON(w, http.StatusConflict, map[string]any{
			"error":      "domain conflict",
			"violations": uniqueErrs,
		})
		return
	}

	// Generate labels
	labels := engine.GenerateLabels(input)

	// Store in DB
	id, err := h.Store.CreateRoute(input.Service, input.Domain, input.Port, input.Entrypoints, "", false, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":     id,
		"labels": labels,
		"route":  input,
	})
}

// DeleteRoute removes a route by ID.
// DELETE /api/v2/routes/{id}
func (h *RouteHandlers) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid route ID")
		return
	}

	if err := h.Store.DeleteRoute(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

// GetRouteLabels generates labels for a route without storing it.
// POST /api/v2/routes/preview
func (h *RouteHandlers) GetRouteLabels(w http.ResponseWriter, r *http.Request) {
	var input engine.RouteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	errs := engine.ValidateRoute(input)
	if len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":      "validation failed",
			"violations": errs,
		})
		return
	}

	labels := engine.GenerateLabels(input)
	writeJSON(w, http.StatusOK, labels)
}
