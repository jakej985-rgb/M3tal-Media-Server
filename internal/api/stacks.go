package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jakej985-rgb/m3tal-core/internal/engine"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// StackHandlers provides operations for stack management.
type StackHandlers struct {
	Store *store.Store
}

// ListStacks scans the stack directory and returns all discovered stacks
// merged with stored state from the database.
// GET /api/v2/stacks
func (h *StackHandlers) ListStacks(w http.ResponseWriter, r *http.Request) {
	stackDir := system.GetStackDir()
	matches, err := filepath.Glob(filepath.Join(stackDir, "*-compose.yml"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot scan stack directory")
		return
	}

	type stackInfo struct {
		Name        string   `json:"name"`
		ComposePath string   `json:"compose_path"`
		Services    []string `json:"services,omitempty"`
		Status      string   `json:"status"`
	}

	var stacks []stackInfo
	for _, match := range matches {
		base := filepath.Base(match)
		name := strings.TrimSuffix(base, "-compose.yml")

		info := stackInfo{
			Name:        name,
			ComposePath: match,
			Status:      "discovered",
		}

		// Try to parse for service names
		if cf, err := engine.ParseCompose(match); err == nil {
			info.Services = cf.ServiceNames()
		}

		// Check DB for stored status
		if h.Store != nil {
			_ = h.Store.UpsertStack(name, match)
			dbStacks, _ := h.Store.ListStacks()
			for _, ds := range dbStacks {
				if ds.Name == name && ds.Status != "" {
					info.Status = ds.Status
				}
			}
		}

		stacks = append(stacks, info)
	}

	if stacks == nil {
		stacks = []stackInfo{}
	}
	writeJSON(w, http.StatusOK, stacks)
}

// LoadStack parses a compose file and returns its services, ports, and labels.
// POST /api/v2/stacks/load
func (h *StackHandlers) LoadStack(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	// Security: only allow paths under the stack directory
	stackDir := system.GetStackDir()
	absPath, _ := filepath.Abs(req.Path)
	absStackDir, _ := filepath.Abs(stackDir)
	if !strings.HasPrefix(absPath, absStackDir) {
		// Also allow /docker paths
		if !strings.HasPrefix(absPath, "/docker") && !strings.HasPrefix(absPath, system.StackPath) {
			writeError(w, http.StatusForbidden, "path must be within the stack directory")
			return
		}
	}

	cf, err := engine.ParseCompose(req.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Build service details
	type serviceDetail struct {
		Name          string            `json:"name"`
		Image         string            `json:"image"`
		Ports         []string          `json:"ports,omitempty"`
		ContainerPort int               `json:"container_port,omitempty"`
		Labels        map[string]string `json:"labels,omitempty"`
		Networks      []string          `json:"networks,omitempty"`
	}

	var services []serviceDetail
	for name, svc := range cf.Services {
		sd := serviceDetail{
			Name:          name,
			Image:         svc.Image,
			Ports:         svc.Ports,
			ContainerPort: cf.GetContainerPort(name),
			Labels:        cf.GetTraefikLabels(name),
			Networks:      svc.Networks.Values,
		}
		services = append(services, sd)
	}

	// Validate
	validationErrors := engine.ValidateCompose(cf)

	writeJSON(w, http.StatusOK, map[string]any{
		"path":              req.Path,
		"services":          services,
		"service_count":     len(services),
		"validation_errors": validationErrors,
	})
}

// DeployStack validates and deploys a stack.
// POST /api/v2/stacks/deploy
func (h *StackHandlers) DeployStack(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	// Verify file exists
	if _, err := os.Stat(req.Path); err != nil {
		writeError(w, http.StatusNotFound, "compose file not found: "+req.Path)
		return
	}

	// Parse to get stack name
	base := filepath.Base(req.Path)
	stackName := strings.TrimSuffix(base, "-compose.yml")

	// Deploy
	result, err := engine.DeployStack(req.Path, 0)

	// Update DB status regardless of outcome
	if h.Store != nil {
		_ = h.Store.UpsertStack(stackName, req.Path)
		if err != nil {
			_ = h.Store.UpdateStackStatus(stackName, "failed")
		} else {
			_ = h.Store.UpdateStackStatus(stackName, "running")
		}
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error":  err.Error(),
			"result": result,
		})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// MiddlewareHandlers provides CRUD for Traefik middleware templates.
type MiddlewareHandlers struct {
	Store *store.Store
}

// ListMiddleware returns all stored middleware definitions.
// GET /api/v2/middleware
func (h *MiddlewareHandlers) ListMiddleware(w http.ResponseWriter, r *http.Request) {
	mws, err := h.Store.ListMiddleware()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if mws == nil {
		mws = []store.MiddlewareRecord{}
	}
	writeJSON(w, http.StatusOK, mws)
}

// CreateMiddleware creates a new middleware definition and generates its labels.
// POST /api/v2/middleware
func (h *MiddlewareHandlers) CreateMiddleware(w http.ResponseWriter, r *http.Request) {
	var input engine.MiddlewareInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if input.Type == "" {
		writeError(w, http.StatusBadRequest, "type is required")
		return
	}

	// Generate labels
	labels := engine.GenerateMiddlewareLabels(input)

	// Store
	id, err := h.Store.CreateMiddleware(input.Name, input.Type, input.Config)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":     id,
		"labels": labels,
	})
}
