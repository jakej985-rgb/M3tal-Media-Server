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
	return &Stack{
		Files: []string{
			"source/m3tal-stack/network-compose.yml",
			"source/m3tal-stack/routing-compose.yml",
			"source/m3tal-stack/m3tal-compose.yml",
		},
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
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %s on %s: %w", action, file, err)
		}
	}
	return nil
}
