package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestAIRun_NotActive(t *testing.T) {
	origActive := isAIActiveCheck
	defer func() { isAIActiveCheck = origActive }()

	isAIActiveCheck = func() bool { return false }

	srv := NewServer("test-token")
	reqBody, _ := json.Marshal(AIRequest{Prompt: "hello"})
	req := httptest.NewRequest("POST", "/ai/run", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	srv.AIRun(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
	if resp.Status != "error" {
		t.Errorf("expected status 'error', got %q", resp.Status)
	}
	if resp.Error != "AI addon is not active or enabled" {
		t.Errorf("expected error message, got %q", resp.Error)
	}
}

func TestAIRun_Success(t *testing.T) {
	origActive := isAIActiveCheck
	origEnv := loadAIEnvFunc
	defer func() {
		isAIActiveCheck = origActive
		loadAIEnvFunc = origEnv
	}()

	isAIActiveCheck = func() bool { return true }

	ollamaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]any
		json.NewDecoder(r.Body).Decode(&payload)
		if payload["model"] != "mock-model" {
			t.Errorf("expected model 'mock-model', got %q", payload["model"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "hello from mock-model"}`))
	}))
	defer ollamaServer.Close()

	loadAIEnvFunc = func() map[string]string {
		return map[string]string{
			"AI_MODEL":    "mock-model",
			"OLLAMA_HOST": ollamaServer.URL,
		}
	}

	srv := NewServer("test-token")
	reqBody, _ := json.Marshal(AIRequest{Prompt: "say hello", Mode: "chat"})
	req := httptest.NewRequest("POST", "/ai/run", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	srv.AIRun(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var apiResponse struct {
		Status string     `json:"status"`
		Data   AIResponse `json:"data"`
		Error  any        `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &apiResponse); err != nil {
		t.Fatalf("failed to parse json response: %v, body: %s", err, w.Body.String())
	}
	if apiResponse.Status != "success" {
		t.Fatalf("expected success status, got %q (error: %v)", apiResponse.Status, apiResponse.Error)
	}
	resp := apiResponse.Data
	if resp.Model != "mock-model" {
		t.Errorf("expected model 'mock-model', got %q", resp.Model)
	}
	if resp.Response != "hello from mock-model" {
		t.Errorf("expected response 'hello from mock-model', got %q", resp.Response)
	}
	if resp.Status != "success" {
		t.Errorf("expected status 'success', got %q", resp.Status)
	}
}

func TestAIRun_Fallback(t *testing.T) {
	origActive := isAIActiveCheck
	origEnv := loadAIEnvFunc
	defer func() {
		isAIActiveCheck = origActive
		loadAIEnvFunc = origEnv
	}()

	isAIActiveCheck = func() bool { return true }

	attempts := 0
	ollamaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		json.NewDecoder(r.Body).Decode(&payload)
		attempts++

		if payload["model"] == "primary-model" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "model not found"}`))
			return
		}

		if payload["model"] == "fallback-model" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"response": "recovered using fallback"}`))
			return
		}
	}))
	defer ollamaServer.Close()

	loadAIEnvFunc = func() map[string]string {
		return map[string]string{
			"AI_MODEL":      "primary-model",
			"AI_FALLBACK_1": "fallback-model",
			"OLLAMA_HOST":   ollamaServer.URL,
		}
	}

	srv := NewServer("test-token")
	reqBody, _ := json.Marshal(AIRequest{Prompt: "say hello"})
	req := httptest.NewRequest("POST", "/ai/run", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	srv.AIRun(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var apiResponse struct {
		Status string     `json:"status"`
		Data   AIResponse `json:"data"`
		Error  any        `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &apiResponse); err != nil {
		t.Fatalf("failed to parse json response: %v, body: %s", err, w.Body.String())
	}
	if apiResponse.Status != "success" {
		t.Fatalf("expected success status, got %q (error: %v)", apiResponse.Status, apiResponse.Error)
	}
	resp := apiResponse.Data
	if resp.Model != "fallback-model" {
		t.Errorf("expected fallback model to be chosen, got %q", resp.Model)
	}
	if resp.Response != "recovered using fallback" {
		t.Errorf("expected response 'recovered using fallback', got %q", resp.Response)
	}
}

func TestAIRun_QueueConcurrency(t *testing.T) {
	origActive := isAIActiveCheck
	origEnv := loadAIEnvFunc
	defer func() {
		isAIActiveCheck = origActive
		loadAIEnvFunc = origEnv
	}()

	isAIActiveCheck = func() bool { return true }

	var mu sync.Mutex
	activeRequests := 0
	maxConcurrency := 0
	completedRequests := 0

	ollamaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		activeRequests++
		if activeRequests > maxConcurrency {
			maxConcurrency = activeRequests
		}
		mu.Unlock()

		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		activeRequests--
		completedRequests++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "done"}`))
	}))
	defer ollamaServer.Close()

	loadAIEnvFunc = func() map[string]string {
		return map[string]string{
			"AI_MODEL":    "mock-model",
			"OLLAMA_HOST": ollamaServer.URL,
		}
	}

	srv := NewServer("test-token")

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			reqBody, _ := json.Marshal(AIRequest{Prompt: fmt.Sprintf("prompt %d", id)})
			req := httptest.NewRequest("POST", "/ai/run", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()
			srv.AIRun(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
			}
		}(i)
	}

	wg.Wait()

	if maxConcurrency != 1 {
		t.Errorf("expected max concurrency 1, got %d", maxConcurrency)
	}
	if completedRequests != 3 {
		t.Errorf("expected 3 completed requests, got %d", completedRequests)
	}
}
