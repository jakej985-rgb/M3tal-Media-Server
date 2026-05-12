package api

import (
	"encoding/json"
	"net/http"

	"github.com/jakej985-rgb/m3tal-core/pkg/containers"
	"github.com/jakej985-rgb/m3tal-core/pkg/system"
)

// Server handles API requests
type Server struct {
	Manager  *containers.Manager
	APIToken string
}

// NewServer creates a new API server
func NewServer(mgr *containers.Manager, token string) *Server {
	return &Server{Manager: mgr, APIToken: token}
}

// AuthMiddleware validates the API token
func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-Token")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token != s.APIToken {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		next(w, r)
	}
}

// ListContainers returns all containers
func (s *Server) ListContainers(w http.ResponseWriter, r *http.Request) {
	list, err := s.Manager.ListContainers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GetStats returns system metrics
func (s *Server) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := system.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// HandleContainerAction processes start/stop/restart
func (s *Server) HandleContainerAction(w http.ResponseWriter, r *http.Request, action func(string) error) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := action(req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
