package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents M3TAL client/CLI configuration.
type Config struct {
	APIURL   string `json:"api_url"`
	APIToken string `json:"api_token"`
}

var (
	globalConfig Config
)

func init() {
	// Initialize default values
	globalConfig = Config{
		APIURL:   "http://localhost:5050",
		APIToken: os.Getenv("API_TOKEN"),
	}
	// Try to load ~/.m3tal/config
	_ = Load()
}

// Load loads the configuration from ~/.m3tal/config.
func Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, ".m3tal", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}

	if loaded.APIURL != "" {
		globalConfig.APIURL = loaded.APIURL
	}
	if loaded.APIToken != "" {
		globalConfig.APIToken = loaded.APIToken
	}

	return nil
}

// GetAPIURL returns the configured API base URL.
func GetAPIURL() string {
	return globalConfig.APIURL
}

// GetAPIToken returns the configured API authorization token.
func GetAPIToken() string {
	if globalConfig.APIToken != "" {
		return globalConfig.APIToken
	}
	return os.Getenv("API_TOKEN")
}
