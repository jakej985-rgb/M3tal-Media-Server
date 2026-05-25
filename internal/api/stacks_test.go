package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestListStacks(t *testing.T) {
	// Create temporary stack directory
	stackDir, err := os.MkdirTemp("", "m3tal-test-api-stacks")
	if err != nil {
		t.Fatalf("failed to create temp stack dir: %v", err)
	}
	defer os.RemoveAll(stackDir)

	t.Setenv("M3TAL_STACK_DIR", stackDir)

	// Create dummy compose files
	filesToCreate := []string{
		"network-compose.yml",
		"routing-compose.yml",
		"custom-compose.yml",
	}
	for _, f := range filesToCreate {
		path := filepath.Join(stackDir, f)
		if err := os.WriteFile(path, []byte("version: '3'\nservices:\n  dummy:\n    image: nginx"), 0644); err != nil {
			t.Fatalf("failed to write dummy compose %s: %v", f, err)
		}
	}

	// Create temp deploy/plugins structure in current working directory of test (internal/api)
	err = os.MkdirAll("deploy/plugins/stacks", 0755)
	if err != nil {
		t.Fatalf("failed to create local deploy/plugins dir: %v", err)
	}
	defer os.RemoveAll("deploy")

	// Write mock plugin manifests
	networkPluginYAML := `
apiVersion: m3tal/v1
kind: Stack
metadata:
  name: network
spec:
  composePath: network-compose.yml
  priority: 1
`
	routingPluginYAML := `
apiVersion: m3tal/v1
kind: Stack
metadata:
  name: routing
spec:
  composePath: routing-compose.yml
  priority: 2
`

	err = os.WriteFile("deploy/plugins/stacks/network.yml", []byte(networkPluginYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write mock network plugin: %v", err)
	}
	err = os.WriteFile("deploy/plugins/stacks/routing.yml.disabled", []byte(routingPluginYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write mock routing plugin: %v", err)
	}

	// Make API request to ListStacks
	h := &StackHandlers{}
	req := httptest.NewRequest("GET", "/api/v2/stacks", nil)
	w := httptest.NewRecorder()

	h.ListStacks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	type stackInfo struct {
		Name        string   `json:"name"`
		ComposePath string   `json:"compose_path"`
		Services    []string `json:"services,omitempty"`
		Status      string   `json:"status"`
	}

	var apiResponse struct {
		Status string      `json:"status"`
		Data   []stackInfo `json:"data"`
		Error  any         `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body: %s", err, w.Body.String())
	}
	if apiResponse.Status != "success" {
		t.Fatalf("expected success status, got %q (error: %v)", apiResponse.Status, apiResponse.Error)
	}
	response := apiResponse.Data

	// We expect routing to be filtered out (because it is disabled),
	// and network to be first (priority 1), custom to be second (priority 100).
	if len(response) != 2 {
		t.Fatalf("expected 2 stacks, got %d: %v", len(response), response)
	}

	if response[0].Name != "network" {
		t.Errorf("expected first stack to be 'network', got %q", response[0].Name)
	}
	if response[1].Name != "custom" {
		t.Errorf("expected second stack to be 'custom', got %q", response[1].Name)
	}
}
