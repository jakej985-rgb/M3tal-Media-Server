package system

import (
	"os"
	"path/filepath"
)

const (
	// DefaultSystemConfigDir is the root for configuration
	DefaultSystemConfigDir = "/etc/m3tal"

	// DefaultSystemStackDir is where compose files live
	DefaultSystemStackDir = "/opt/m3tal/stack"

	// DefaultSystemDataDir is for persistent data
	DefaultSystemDataDir = "/var/lib/m3tal"

	// DefaultSystemLogDir is for logs
	DefaultSystemLogDir = "/var/log/m3tal"

	// DefaultSystemStacksDir is where individual stack env files live
	DefaultSystemStacksDir = "/etc/m3tal/stacks"

	// DefaultSystemGlobalEnvPath is the global env file
	DefaultSystemGlobalEnvPath = "/etc/m3tal/global.env"
)

// GetConfigPath returns the path to the main configuration file.
// It prioritizes M3TAL_CONFIG env var, then system path, then local .env.
func GetConfigPath() string {
	if p := os.Getenv("M3TAL_CONFIG"); p != "" {
		return p
	}
	systemPath := filepath.Join(DefaultSystemConfigDir, "config.yaml")
	if _, err := os.Stat(systemPath); err == nil {
		return systemPath
	}
	return ".env"
}

// GetStackDir returns the directory containing stack compose files.
func GetStackDir() string {
	if p := os.Getenv("M3TAL_STACK_DIR"); p != "" {
		return p
	}
	if _, err := os.Stat(DefaultSystemStackDir); err == nil {
		return DefaultSystemStackDir
	}
	// Fallback to local deploy/stack if it exists
	if _, err := os.Stat("deploy/stack"); err == nil {
		return "deploy/stack"
	}
	return "."
}

// GetDataDir returns the root data directory.
func GetDataDir() string {
	if p := os.Getenv("M3TAL_DATA_DIR"); p != "" {
		return p
	}
	return DefaultSystemDataDir
}
