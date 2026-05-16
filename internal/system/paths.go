package system

import (
	"os"
)

const (
	// ConfigPath is the main environment configuration file
	ConfigPath = "/etc/m3tal/.env"

	// StackPath is where the Docker Compose files are stored
	StackPath = "/opt/m3tal/stack"

	// DataPath is the root for persistent application data
	DataPath = "/var/lib/m3tal"

	// UserfacingStackPath is the symlink used for user interaction
	UserfacingStackPath = "/docker"

	// DashPath is where the dashboard files are stored
	DashPath = "/opt/m3tal/dash"

	// StacksDir is where individual stack env files live
	StacksDir = "/etc/m3tal/stacks"
)

// GetConfigPath returns the active configuration path, supporting overrides
func GetConfigPath() string {
	if p := os.Getenv("M3TAL_CONFIG"); p != "" {
		return p
	}
	// Check system path first
	if _, err := os.Stat(ConfigPath); err == nil {
		return ConfigPath
	}
	// Fallback to local .env
	return ".env"
}

// GetStackDir returns the directory containing stack compose files
func GetStackDir() string {
	if p := os.Getenv("M3TAL_STACK_DIR"); p != "" {
		return p
	}
	// Check user-facing path first (as per UX rule)
	if _, err := os.Stat(UserfacingStackPath); err == nil {
		return UserfacingStackPath
	}
	// Then system path
	if _, err := os.Stat(StackPath); err == nil {
		return StackPath
	}
	// Fallback to local deploy/stack
	if _, err := os.Stat("deploy/stack"); err == nil {
		return "deploy/stack"
	}
	return "."
}

// IsDashInstalled checks if the dashboard is present on the system
func IsDashInstalled() bool {
	// The plan says check in /opt/m3tal/dash/dash-compose.yml
	// But user corrected: "the compose files will be in /docker"
	// So we check /docker/dash/dash-compose.yml which is /opt/m3tal/stack/dash/dash-compose.yml
	path := "/docker/dash/dash-compose.yml"
	if _, err := os.Stat(path); err == nil {
		return true
	}
	// Also check the real path just in case symlink isn't there yet
	realPath := StackPath + "/dash/dash-compose.yml"
	if _, err := os.Stat(realPath); err == nil {
		return true
	}
	return false
}
