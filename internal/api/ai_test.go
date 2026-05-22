package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
	if resp["error"] != "AI addon is not active or enabled" {
		t.Errorf("expected error message, got %q", resp["error"])
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

	var resp AIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
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

	var resp AIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
	if resp.Model != "fallback-model" {
		t.Errorf("expected fallback model to be chosen, got %q", resp.Model)
	}
	if resp.Response != "recovered using fallback" {
		t.Errorf("expected response 'recovered using fallback', got %q", resp.Response)
	}
}
