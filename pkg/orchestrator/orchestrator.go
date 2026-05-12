package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
)

// Stack represents a collection of docker-compose files
type Stack struct {
	Files []string
}

// NewStack returns a stack with standard M3TAL compose paths
func NewStack() *Stack {
	paths := []string{
		"source/m3tal-stack/network-compose.yml",
		"source/m3tal-stack/routing-compose.yml",
		"source/m3tal-stack/m3tal-compose.yml",
	}

	// Check if we are running as a system installation
	if _, err := os.Stat(paths[0]); os.IsNotExist(err) {
		systemPath := "/usr/share/m3tal/stack/"
		if _, err := os.Stat(systemPath + "network-compose.yml"); err == nil {
			paths = []string{
				systemPath + "network-compose.yml",
				systemPath + "routing-compose.yml",
				systemPath + "m3tal-compose.yml",
			}
		}
	}

	return &Stack{
		Files: paths,
	}
}

// Run executes a docker compose command across all stack files
func (s *Stack) Run(action string, args ...string) error {
	for _, file := range s.Files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("🚀 Running docker compose %s on %s...\n", action, file)
		cmdArgs := append([]string{"compose", "-f", file, action}, args...)
		cmd := exec.Command("docker", cmdArgs...)

		// Pass current environment + .env file if it exists
		cmd.Env = os.Environ()
		envFile := ".env"
		if _, err := os.Stat("/etc/m3tal/config.yaml"); err == nil {
			envFile = "/etc/m3tal/config.yaml"
		}

		if _, err := os.Stat(envFile); err == nil {
			cmdArgs = append([]string{"compose", "--env-file", envFile, "-f", file, action}, args...)
			cmd = exec.Command("docker", cmdArgs...)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %s on %s: %w", action, file, err)
		}
	}
	return nil
}
