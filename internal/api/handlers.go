package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/health"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// Server handles API requests
type Server struct {
	APIToken string
}

// NewServer creates a new API server
func NewServer(token string) *Server {
	return &Server{APIToken: token}
}

// AuthMiddleware validates the API token
func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-Token")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token != s.APIToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// GetHealth returns system health status
func (s *Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	reg := health.UpdateAndSaveHealthRegistry()
	sendSuccess(w, http.StatusOK, reg, nil)
}

// GetServices returns the list of managed containers
func (s *Server) GetServices(w http.ResponseWriter, r *http.Request) {
	mgr, err := containers.GetProvider()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	list, err := mgr.ListContainers()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccess(w, http.StatusOK, list, nil)
}

// GetStack returns information about the compose stack
func (s *Server) GetStack(w http.ResponseWriter, r *http.Request) {
	stackDir := system.GetStackDir()
	sendSuccess(w, http.StatusOK, map[string]string{
		"path": stackDir,
	}, nil)
}

// GetConfig returns the current configuration (sanitized)
func (s *Server) GetConfig(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			// Only include M3TAL related env vars
			if strings.HasPrefix(key, "M3TAL_") || key == "BASE_STORAGE_PATH" {
				// Sanitize sensitive info
				if strings.Contains(key, "TOKEN") || strings.Contains(key, "SECRET") || strings.Contains(key, "PASSWORD") {
					config[key] = "********"
				} else {
					config[key] = pair[1]
				}
			}
		}
	}
	sendSuccess(w, http.StatusOK, config, nil)
}

// GetStats returns system metrics
func (s *Server) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := system.GetStats()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccess(w, http.StatusOK, stats, nil)
}

// HandleContainerAction processes start/stop/restart
func (s *Server) HandleContainerAction(w http.ResponseWriter, r *http.Request, action func(string) error) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := action(req.Name); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccess(w, http.StatusOK, map[string]bool{"ok": true}, nil)
}
