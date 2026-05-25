package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jakej985-rgb/m3tal-core/internal/engine"
	"github.com/jakej985-rgb/m3tal-core/internal/proxy"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
)

// ProxyHandlers wraps the SQLite store to manage reverse proxy endpoints.
type ProxyHandlers struct {
	Store *store.Store
}

// NewProxyHandlers creates a new set of proxy handlers.
func NewProxyHandlers(db *store.Store) *ProxyHandlers {
	return &ProxyHandlers{Store: db}
}

// DiscoverServices auto-discovers containers.
// GET /api/v2/proxy/discover
func (h *ProxyHandlers) DiscoverServices(w http.ResponseWriter, r *http.Request) {
	mgr := proxy.NewManager(h.Store)
	services, err := mgr.DiscoverServices()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccess(w, http.StatusOK, services, nil)
}

// ExposeService registers/exposes a service on a domain.
// POST /api/v2/proxy/expose
func (h *ProxyHandlers) ExposeService(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Service     string   `json:"service"`
		Domain      string   `json:"domain"`
		Port        int      `json:"port"`
		SSL         bool     `json:"ssl"`
		Middlewares []string `json:"middlewares"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Route validation
	routeInput := engine.RouteInput{
		Service: input.Service,
		Domain:  input.Domain,
		Port:    input.Port,
	}
	errs := engine.ValidateRoute(routeInput)
	if len(errs) > 0 {
		writeJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "validation failed",
			Data:   map[string]any{"violations": errs},
		})
		return
	}

	mgr := proxy.NewManager(h.Store)
	err := mgr.ExposeService(input.Service, input.Domain, input.Port, input.SSL, input.Middlewares)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(w, http.StatusOK, map[string]any{
		"ok":      true,
		"exposed": true,
		"route":   input,
	}, nil)
}

// UnexposeService removes route mapping.
// POST /api/v2/proxy/unexpose
func (h *ProxyHandlers) UnexposeService(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if input.Domain == "" {
		writeError(w, http.StatusBadRequest, "domain is required")
		return
	}

	mgr := proxy.NewManager(h.Store)
	err := mgr.UnexposeService(input.Domain)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(w, http.StatusOK, map[string]any{
		"ok":        true,
		"unexposed": true,
	}, nil)
}

// ConfigureSSL triggers SSL automated Let's Encrypt certificates setup.
// POST /api/v2/proxy/secure
func (h *ProxyHandlers) ConfigureSSL(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	email := strings.TrimSpace(input.Email)
	if email == "" {
		writeError(w, http.StatusBadRequest, "valid email is required for SSL configuration")
		return
	}

	mgr := proxy.NewManager(h.Store)
	err := mgr.ConfigureSSL(email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(w, http.StatusOK, map[string]any{
		"ok":      true,
		"secured": true,
		"message": "Let's Encrypt SSL configured and routing stack redeployed.",
	}, nil)
}
