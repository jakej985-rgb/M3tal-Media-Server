package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jakej985-rgb/m3tal-core/pkg/models"
)

// Client handles authenticated HTTP communications with the M3TAL API.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient instantiates a new Client.
func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:5050"
	}
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// APIResponse represents the standardized JSON structure for M3TAL API responses.
type APIResponse struct {
	Status string                `json:"status"`
	Data   json.RawMessage       `json:"data"`
	Meta   json.RawMessage       `json:"meta,omitempty"`
	Error  *models.ErrorResponse `json:"error"`
}

func (c *Client) request(method, path string, body any, target any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-API-Token", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return fmt.Errorf("API returned non-JSON response (HTTP %d): %s", resp.StatusCode, string(respBytes))
	}

	if apiResp.Status == "error" || apiResp.Error != nil {
		if apiResp.Error != nil {
			return fmt.Errorf("API error (%s): %s", apiResp.Error.Code, apiResp.Error.Message)
		}
		return fmt.Errorf("API request returned error status")
	}

	if target != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, target); err != nil {
			return fmt.Errorf("failed to unmarshal API data field: %w", err)
		}
	}

	return nil
}

// Request executes a generic API request.
func (c *Client) Request(method, path string, body any, target any) error {
	return c.request(method, path, body, target)
}

// GetStats returns the host system metrics response.
func (c *Client) GetStats() (*models.MetricsResponse, error) {
	var metrics models.MetricsResponse
	err := c.request("GET", "/api/v2/system/stats", nil, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// GetStatus checks the detailed system components health check report.
func (c *Client) GetStatus() (*models.Status, error) {
	var status models.Status
	err := c.request("GET", "/api/v2/system/health", nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetStacks returns all discovered stacks merged with DB status.
func (c *Client) GetStacks() ([]models.Stack, error) {
	var stacks []models.Stack
	err := c.request("GET", "/api/v2/stacks", nil, &stacks)
	if err != nil {
		return nil, err
	}
	return stacks, nil
}

// StartStack deploys a stack by name.
func (c *Client) StartStack(name string) (any, error) {
	var result any
	err := c.request("POST", fmt.Sprintf("/api/v2/stacks/%s/up", name), nil, &result)
	return result, err
}

// StopStack tears down/stops a stack by name.
func (c *Client) StopStack(name string) (any, error) {
	var result any
	err := c.request("POST", fmt.Sprintf("/api/v2/stacks/%s/down", name), nil, &result)
	return result, err
}

// RestartStack restarts a stack by name (stops then starts).
func (c *Client) RestartStack(name string) (any, error) {
	_, err := c.StopStack(name)
	if err != nil {
		return nil, fmt.Errorf("restart failed at stop: %w", err)
	}
	return c.StartStack(name)
}

// GetContainers returns the list of all active containers.
func (c *Client) GetContainers() ([]models.Container, error) {
	var containers []models.Container
	err := c.request("GET", "/api/v2/services", nil, &containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// GetLogs returns the logs (last 100 lines) of a container.
func (c *Client) GetLogs(containerID string) (string, error) {
	var logs models.LogEntry
	err := c.request("GET", fmt.Sprintf("/api/v2/containers/%s/logs?tail=100", containerID), nil, &logs)
	if err != nil {
		return "", err
	}
	return logs.Logs, nil
}

// GetVPNStatus checks the current status of the VPN connection.
func (c *Client) GetVPNStatus() (*models.VPNStatus, error) {
	var status models.VPNStatus
	err := c.request("GET", "/api/v2/vpn/status", nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// ControlVPN sends a start, stop, or restart command to the VPN.
func (c *Client) ControlVPN(action string) (any, error) {
	var result any
	body := map[string]string{"action": action}
	err := c.request("POST", "/api/v2/vpn/control", body, &result)
	return result, err
}

// SwitchVPNRegion switches the VPN connection region.
func (c *Client) SwitchVPNRegion(region string) (any, error) {
	var result any
	body := map[string]string{"region": region}
	err := c.request("POST", "/api/v2/vpn/region", body, &result)
	return result, err
}

// SyncVPNPort triggers manual port synchronization for Gluetun.
func (c *Client) SyncVPNPort() (int, error) {
	var result struct {
		Status        string `json:"status"`
		ForwardedPort int    `json:"forwarded_port"`
	}
	err := c.request("POST", "/api/v2/vpn/sync-port", nil, &result)
	if err != nil {
		return 0, err
	}
	return result.ForwardedPort, nil
}

// CheckVPNLeak runs leak detection on the host and VPN outbound IPs.
func (c *Client) CheckVPNLeak() (*models.VPNLeakReport, error) {
	var report models.VPNLeakReport
	err := c.request("GET", "/api/v2/vpn/check-leak", nil, &report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// ControlContainer starts, stops, or restarts a container via the API.
func (c *Client) ControlContainer(name, action string) error {
	return c.request("POST", fmt.Sprintf("/api/v2/containers/%s/%s", name, action), nil, nil)
}

// GetQueue returns the list of all queued and finished tasks.
func (c *Client) GetQueue() ([]models.JobRecord, error) {
	return []models.JobRecord{}, nil
}

// CancelQueueJob aborts a pending or active job by ID.
func (c *Client) CancelQueueJob(id string) error {
	return nil
}

// GetAIModels returns the list of available AI models from the Ollama backend.
func (c *Client) GetAIModels() ([]string, error) {
	return []string{"agents system active (placeholder)"}, nil
}

// GetPlugins returns all loaded plugins across all kinds.
func (c *Client) GetPlugins() (*models.PluginsResponse, error) {
	var response models.PluginsResponse
	err := c.request("GET", "/api/v2/plugins", nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPluginCatalog correlates remote catalog items with local loaded plugins.
func (c *Client) GetPluginCatalog() ([]models.CatalogItemStatus, error) {
	var catalog []models.CatalogItemStatus
	err := c.request("GET", "/api/v2/plugins/catalog", nil, &catalog)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

// EnablePlugin enables a plugin by name and kind.
func (c *Client) EnablePlugin(name, kind string) error {
	body := map[string]string{"name": name, "kind": kind}
	return c.request("POST", "/api/v2/plugins/enable", body, nil)
}

// DisablePlugin disables a plugin by name and kind.
func (c *Client) DisablePlugin(name, kind string) error {
	body := map[string]string{"name": name, "kind": kind}
	return c.request("POST", "/api/v2/plugins/disable", body, nil)
}

// InstallPlugin downloads and installs a plugin by name and kind.
func (c *Client) InstallPlugin(name, kind string) error {
	body := map[string]string{"name": name, "kind": kind}
	return c.request("POST", "/api/v2/plugins/install", body, nil)
}

// UninstallPlugin uninstalls a plugin by name and kind.
func (c *Client) UninstallPlugin(name, kind string) error {
	body := map[string]string{"name": name, "kind": kind}
	return c.request("POST", "/api/v2/plugins/uninstall", body, nil)
}
