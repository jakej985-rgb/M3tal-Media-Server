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

	envVars := loadAIEnvFunc()

	// 1. Determine model based on mode
	var model string
	switch strings.ToLower(req.Mode) {
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
		aiResp, lastErr = callOllamaGenerate(ollamaHost, m, req.Prompt)
		if lastErr == nil {
			break
		}
		log.Printf("⚠️ AI API fallback: model %s failed, attempting next if available. Error: %v", m, lastErr)
	}

	if lastErr != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("AI execution failed: %v", lastErr))
		return
	}

	writeJSON(w, http.StatusOK, AIResponse{
		Model:    chosenModel,
		Response: aiResp,
		Status:   "success",
	})
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
