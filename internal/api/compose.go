package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jakej985-rgb/m3tal-core/internal/compose"
	"github.com/jakej985-rgb/m3tal-core/internal/engine"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// ComposeHandlers provides HTTP handlers for Smart Compose operations.
type ComposeHandlers struct {
	Store *store.Store
}

// ValidateCompose handles validation and linting of a compose YAML string.
// POST /api/v2/compose/validate
func (h *ComposeHandlers) ValidateCompose(w http.ResponseWriter, r *http.Request) {
	var req struct {
		YAML string `json:"yaml"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cfg, err := compose.Parse([]byte(req.YAML))
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"valid":  false,
			"errors": []string{err.Error()},
			"issues": []compose.LintIssue{},
		})
		return
	}

	issues := compose.Lint(cfg)
	hasErrors := false
	var errs []string
	for _, issue := range issues {
		if issue.Severity == compose.SeverityError {
			hasErrors = true
			errs = append(errs, issue.Message)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"valid":  !hasErrors,
		"errors": errs,
		"issues": issues,
	})
}

// FixCompose handles auto-fixing a compose YAML string.
// POST /api/v2/compose/fix
func (h *ComposeHandlers) FixCompose(w http.ResponseWriter, r *http.Request) {
	var req struct {
		YAML string `json:"yaml"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	fixed, fixes, err := compose.AutoFix([]byte(req.YAML))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"yaml":  string(fixed),
		"fixes": fixes,
	})
}

// ListTemplates handles listing the pre-defined compose templates.
// GET /api/v2/compose/templates
func (h *ComposeHandlers) ListTemplates(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, compose.Templates)
}

// GenerateTemplate handles generating a compose YAML from a template.
// POST /api/v2/compose/generate
func (h *ComposeHandlers) GenerateTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Template   string            `json:"template"`
		Parameters map[string]string `json:"parameters"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	yamlData, err := compose.Generate(req.Template, req.Parameters)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"yaml": string(yamlData),
	})
}

// SaveCompose handles saving a compose file and deploying it.
// POST /api/v2/compose/save
func (h *ComposeHandlers) SaveCompose(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Stack string `json:"stack"`
		YAML  string `json:"yaml"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Stack == "" {
		writeError(w, http.StatusBadRequest, "stack name is required")
		return
	}

	// Clean/sanitize stack name to prevent path traversal
	stackName := strings.TrimSpace(req.Stack)
	stackName = filepath.Base(stackName)
	stackName = strings.TrimSuffix(stackName, "-compose.yml")

	// Target path
	stackDir := system.GetStackDir()
	filePath := filepath.Join(stackDir, fmt.Sprintf("%s-compose.yml", stackName))

	// Write file
	if err := os.WriteFile(filePath, []byte(req.YAML), 0644); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save file: "+err.Error())
		return
	}

	// Deploy stack
	result, err := engine.DeployStack(filePath, 0)

	// Update DB status
	if h.Store != nil {
		_ = h.Store.UpsertStack(stackName, filePath)
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

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "success",
		"result": result,
	})
}
