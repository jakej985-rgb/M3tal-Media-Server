package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jakej985-rgb/m3tal-core/internal/queue"
)

// GlobalQueue manages background priority-based worker tasks.
var GlobalQueue = queue.NewManager(1, 100) // 1 worker for sequential execution, 100 queue capacity

// ListQueue returns the status of all queued and finished tasks.
// GET /api/v2/queue
func ListQueue(w http.ResponseWriter, r *http.Request) {
	jobs := GlobalQueue.List()
	sendSuccess(w, http.StatusOK, jobs, nil)
}

// GetQueueJob returns the status and result of a single task by ID.
// GET /api/v2/queue/{id}
func GetQueueJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing job ID")
		return
	}

	record, found := GlobalQueue.Get(id)
	if !found {
		writeError(w, http.StatusNotFound, "job not found: "+id)
		return
	}

	sendSuccess(w, http.StatusOK, record, nil)
}

// CancelQueueJob aborts a pending or active job.
// POST /api/v2/queue/cancel
func CancelQueueJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}

	// Try decoding body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		// Fallback to query param
		req.ID = r.URL.Query().Get("id")
	}

	if req.ID == "" {
		writeError(w, http.StatusBadRequest, "missing job ID")
		return
	}

	cancelled := GlobalQueue.Cancel(req.ID)
	if !cancelled {
		writeError(w, http.StatusConflict, "failed to cancel job (either not found, or already completed/failed)")
		return
	}

	sendSuccess(w, http.StatusOK, map[string]any{
		"cancelled": true,
		"id":        req.ID,
	}, nil)
}
