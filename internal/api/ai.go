package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jakej985-rgb/m3tal-core/internal/plugin"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// AIRequest represents the input to the AI endpoint
type AIRequest struct {
	Prompt string `json:"prompt"`
	Mode   string `json:"mode,omitempty"`
}

// AIResponse represents the output of the AI endpoint
type AIResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Status   string `json:"status"`
}

// AIJobStatus represents the current status of an AI job
type AIJobStatus struct {
	ID     string `json:"id"`
	Prompt string `json:"prompt"`
	Mode   string `json:"mode"`
	Status string `json:"status"` // "pending", "running", "completed", "failed"
}

// AIJob represents a queued AI generation task
type AIJob struct {
	ID         string
	Prompt     string
	Mode       string
	ResultChan chan AIJobResult
}

// AIJobResult represents the result of a processed AI job
type AIJobResult struct {
	Model    string
	Response string
	Error    error
}

var (
	jobQueue        chan *AIJob
	startWorkerOnce sync.Once
	jobsTracker     []*AIJobStatus
	jobsTrackerMu   sync.Mutex
	jobCounter      uint64
	jobCounterMu    sync.Mutex
)

func generateJobID() string {
	jobCounterMu.Lock()
	defer jobCounterMu.Unlock()
	jobCounter++
	return fmt.Sprintf("job-%d", jobCounter)
}

// startWorker initializes the job queue and starts the worker loop
func startWorker() {
	jobQueue = make(chan *AIJob, 100)
	go workerLoop()
}

// workerLoop processes queued AI jobs sequentially (MAX=1 concurrency)
func workerLoop() {
	for job := range jobQueue {
		jobsTrackerMu.Lock()
		for _, js := range jobsTracker {
			if js.ID == job.ID {
				js.Status = "running"
				GlobalEventBus.Publish("ai.job.started", map[string]string{
					"id":     job.ID,
					"prompt": job.Prompt,
					"mode":   job.Mode,
				})
				break
			}
		}
		jobsTrackerMu.Unlock()

		envVars := loadAIEnvFunc()

		// 1. Determine model based on mode
		var model string
		switch strings.ToLower(job.Mode) {
		case "code":
			if m, ok := envVars["AI_MODEL_CODE"]; ok && m != "" {
				model = m
			}
		case "chat":
			if m, ok := envVars["AI_MODEL_CHAT"]; ok && m != "" {
				model = m
			}
		}
		if model == "" {
			if m, ok := envVars["AI_MODEL"]; ok && m != "" {
				model = m
			} else {
				model = "qwen3-coder-next:cloud"
			}
		}

		// 2. Build list of fallback models
		var modelsToTry []string
		modelsToTry = append(modelsToTry, model)
		for i := 1; i <= 10; i++ {
			fallbackKey := fmt.Sprintf("AI_FALLBACK_%d", i)
			if fallbackVal, ok := envVars[fallbackKey]; ok && fallbackVal != "" {
				modelsToTry = append(modelsToTry, fallbackVal)
			}
		}

		// 3. Resolve Ollama Host address
		ollamaHost := "http://localhost:11434"
		if h, ok := envVars["OLLAMA_HOST"]; ok && h != "" {
			ollamaHost = h
		}

		var lastErr error
		var chosenModel string
		var aiResp string

		// 4. Fallback execution retry loop
		for _, m := range modelsToTry {
			chosenModel = m
			aiResp, lastErr = callOllamaGenerate(ollamaHost, m, job.Prompt)
			if lastErr == nil {
				break
			}
			log.Printf("⚠️ AI API fallback: model %s failed, attempting next if available. Error: %v", m, lastErr)
		}

		job.ResultChan <- AIJobResult{
			Model:    chosenModel,
			Response: aiResp,
			Error:    lastErr,
		}

		jobsTrackerMu.Lock()
		for _, js := range jobsTracker {
			if js.ID == job.ID {
				if lastErr != nil {
					js.Status = "failed"
					GlobalEventBus.Publish("ai.job.completed", map[string]any{
						"id":     job.ID,
						"status": "failed",
						"error":  lastErr.Error(),
					})
				} else {
					js.Status = "completed"
					GlobalEventBus.Publish("ai.job.completed", map[string]any{
						"id":     job.ID,
						"status": "completed",
						"model":  chosenModel,
					})
				}
				break
			}
		}
		jobsTrackerMu.Unlock()
	}
}

// isAIAddonActive checks if the AI stack plugin is loaded and enabled
func isAIAddonActive() bool {
	dirs := system.GetPluginDirs()
	reg, err := plugin.LoadAll(dirs...)
	if err != nil {
		return false
	}
	stack := reg.GetStack("ai")
	return stack != nil && stack.Enabled
}

// loadAIEnv reads the active environment configuration profile on-demand
func loadAIEnv() map[string]string {
	env := make(map[string]string)
	paths := []string{
		"deploy/plugins/ai/ai.env",
		"../../deploy/plugins/ai/ai.env",
		"/etc/m3tal/plugins/ai/ai.env",
		"/opt/m3tal/plugins/ai/ai.env",
	}
	var data []byte
	var err error
	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if len(data) == 0 {
		return env
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			env[key] = val
		}
	}
	return env
}

// isAIActiveCheck is a package-level variable function that can be overridden in tests
var isAIActiveCheck = isAIAddonActive

// loadAIEnvFunc is a package-level variable function that can be overridden in tests
var loadAIEnvFunc = loadAIEnv

// AIRun handles request routing, model selection, and fallback execution
func (s *Server) AIRun(w http.ResponseWriter, r *http.Request) {
	if !isAIActiveCheck() {
		writeError(w, http.StatusServiceUnavailable, "AI addon is not active or enabled")
		return
	}

	var req AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Prompt == "" {
		writeError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	startWorkerOnce.Do(startWorker)

	jobID := generateJobID()
	status := &AIJobStatus{
		ID:     jobID,
		Prompt: req.Prompt,
		Mode:   req.Mode,
		Status: "pending",
	}

	jobsTrackerMu.Lock()
	jobsTracker = append(jobsTracker, status)
	if len(jobsTracker) > 50 {
		jobsTracker = jobsTracker[len(jobsTracker)-50:]
	}
	jobsTrackerMu.Unlock()

	resultChan := make(chan AIJobResult, 1)
	job := &AIJob{
		ID:         jobID,
		Prompt:     req.Prompt,
		Mode:       req.Mode,
		ResultChan: resultChan,
	}

	select {
	case jobQueue <- job:
		// Queued successfully
	case <-r.Context().Done():
		jobsTrackerMu.Lock()
		for _, js := range jobsTracker {
			if js.ID == jobID {
				js.Status = "failed"
				break
			}
		}
		jobsTrackerMu.Unlock()
		writeError(w, http.StatusRequestTimeout, "request cancelled or timed out before queueing")
		return
	}

	select {
	case result := <-resultChan:
		if result.Error != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("AI execution failed: %v", result.Error))
			return
		}
		sendSuccess(w, http.StatusOK, AIResponse{
			Model:    result.Model,
			Response: result.Response,
			Status:   "success",
		}, nil)
	case <-r.Context().Done():
		jobsTrackerMu.Lock()
		for _, js := range jobsTracker {
			if js.ID == jobID {
				js.Status = "failed"
				break
			}
		}
		jobsTrackerMu.Unlock()
		writeError(w, http.StatusRequestTimeout, "request cancelled or timed out during execution")
	}
}

// callOllamaGenerate triggers standard Ollama /api/generate REST API request
func callOllamaGenerate(host string, model string, prompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", strings.TrimSuffix(host, "/"))
	payload := map[string]any{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var ollamaResp struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(respBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w (raw response: %s)", err, string(respBytes))
	}

	return ollamaResp.Response, nil
}

// GetAIQueue returns the list of queued/running AI jobs
func (s *Server) GetAIQueue(w http.ResponseWriter, r *http.Request) {
	jobsTrackerMu.Lock()
	defer jobsTrackerMu.Unlock()

	list := make([]AIJobStatus, len(jobsTracker))
	for i, j := range jobsTracker {
		list[i] = *j
	}
	sendSuccess(w, http.StatusOK, list, nil)
}

// GetAIModels returns the list of active/fallback models from Ollama or env config
func (s *Server) GetAIModels(w http.ResponseWriter, r *http.Request) {
	envVars := loadAIEnvFunc()
	ollamaHost := "http://localhost:11434"
	if h, ok := envVars["OLLAMA_HOST"]; ok && h != "" {
		ollamaHost = h
	}

	url := fmt.Sprintf("%s/api/tags", strings.TrimSuffix(ollamaHost, "/"))
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var ollamaResp struct {
			Models []struct {
				Name string `json:"name"`
			} `json:"models"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err == nil {
			var models []string
			for _, m := range ollamaResp.Models {
				models = append(models, m.Name)
			}
			sendSuccess(w, http.StatusOK, models, nil)
			return
		}
	}

	// Fallback to configured models in env
	var models []string
	if m := envVars["AI_MODEL"]; m != "" {
		models = append(models, m)
	}
	if m := envVars["AI_MODEL_CHAT"]; m != "" {
		models = append(models, m)
	}
	if m := envVars["AI_MODEL_CODE"]; m != "" {
		models = append(models, m)
	}
	for i := 1; i <= 10; i++ {
		if m := envVars[fmt.Sprintf("AI_FALLBACK_%d", i)]; m != "" {
			models = append(models, m)
		}
	}
	if len(models) == 0 {
		models = []string{"qwen3-coder-next:cloud", "llama3:latest"}
	}
	sendSuccess(w, http.StatusOK, models, nil)
}
