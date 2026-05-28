package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/health"
	"github.com/jakej985-rgb/m3tal-core/core/state/system"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
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
			sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
			return
		}
		next(w, r)
	}
}

// GetHealth returns system health status
func (s *Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	reg := health.UpdateAndSaveHealthRegistry()

	components := map[string]string{
		"system": reg.System.Status,
		"docker": reg.Docker.Status,
		"agents": reg.Agents.Status,
		"disk":   reg.Disk.Status,
	}
	details := map[string]string{
		"last_seen_healthy": reg.System.LastSeenHealthy,
		"last_failure":      reg.System.LastFailure,
	}

	sendSuccess(w, http.StatusOK, models.Status{
		Status:     reg.System.Status,
		Components: components,
		Details:    details,
	}, nil)
}

// GetServices returns the list of managed containers
func (s *Server) GetServices(w http.ResponseWriter, r *http.Request) {
	mgr, err := containers.GetProvider()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "DOCKER_UNAVAILABLE", err.Error(), nil)
		return
	}
	list, err := mgr.ListContainers()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "DOCKER_ERROR", err.Error(), nil)
		return
	}

	typedList := make([]models.Container, len(list))
	for i, c := range list {
		ports := make([]models.PortInfo, len(c.Ports))
		for j, p := range c.Ports {
			ports[j] = models.PortInfo{
				IP:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			}
		}
		typedList[i] = models.Container{
			ID:       c.ID,
			Names:    c.Names,
			Image:    c.Image,
			Status:   c.Status,
			State:    c.State,
			CPU:      c.CPU,
			Memory:   c.Memory,
			Labels:   c.Labels,
			Ports:    ports,
			Networks: c.Networks,
		}
	}
	sendSuccess(w, http.StatusOK, typedList, nil)
}

// GetStack returns information about the compose stack
func (s *Server) GetStack(w http.ResponseWriter, r *http.Request) {
	stackDir := system.GetStackDir()
	sendSuccess(w, http.StatusOK, models.Stack{
		Name:        "default",
		ComposePath: stackDir,
		Status:      "active",
	}, nil)
}

// GetConfig returns the current configuration (sanitized)
func (s *Server) GetConfig(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			if strings.HasPrefix(key, "M3TAL_") || key == "BASE_STORAGE_PATH" {
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
		sendError(w, http.StatusInternalServerError, "SYSTEM_METRICS_ERROR", err.Error(), nil)
		return
	}
	sendSuccess(w, http.StatusOK, models.MetricsResponse{
		CPUUsage:    stats.CPUUsage,
		MemoryUsage: stats.MemoryUsage,
		DiskUsage:   stats.DiskUsage,
		Uptime:      stats.Uptime,
		Hostname:    stats.Hostname,
	}, nil)
}

// HandleContainerAction processes start/stop/restart
func (s *Server) HandleContainerAction(w http.ResponseWriter, r *http.Request, action func(string) error) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_JSON", err.Error(), nil)
		return
	}
	if req.Name == "" {
		sendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "container name is required", nil)
		return
	}
	if err := action(req.Name); err != nil {
		sendError(w, http.StatusInternalServerError, "CONTAINER_ACTION_FAILED", err.Error(), nil)
		return
	}
	sendSuccess(w, http.StatusOK, map[string]bool{"ok": true}, nil)
}

// GetContainerLogs returns logs for a container
func (s *Server) GetContainerLogs(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		sendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "container name is required", nil)
		return
	}
	tail := r.URL.Query().Get("tail")
	if tail == "" {
		tail = "all"
	}

	mgr, err := containers.GetProvider()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "DOCKER_UNAVAILABLE", err.Error(), nil)
		return
	}

	logs, err := mgr.Logs(name, tail)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "DOCKER_ERROR", err.Error(), nil)
		return
	}

	sendSuccess(w, http.StatusOK, map[string]string{"logs": logs}, nil)
}
