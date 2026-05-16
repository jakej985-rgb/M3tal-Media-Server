package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/jakej985-rgb/m3tal-core/internal/system"
)

// StackManager represents a collection of docker-compose files
type StackManager struct {
	Files []string
}

// NewStackManager returns a manager with discovered M3TAL compose files.
func NewStackManager() *StackManager {
	paths := discoverComposeFiles()
	return &StackManager{
		Files: paths,
	}
}

// discoverComposeFiles finds all matching compose files across known locations.
func discoverComposeFiles() []string {
	var files []string

	// Scan system paths from helper (prioritizes /docker or /opt/m3tal/stack)
	stackDir := system.GetStackDir()
	matches, _ := filepath.Glob(filepath.Join(stackDir, "*-compose.yml"))
	files = append(files, matches...)

	// Also check for files in subdirectories (legacy/extra stacks)
	subMatches, _ := filepath.Glob(filepath.Join(stackDir, "*", "*-compose.yml"))
	files = append(files, subMatches...)

	// Deduplicate and sort for deterministic order
	files = uniqueSorted(files)

	return files
}

// uniqueSorted deduplicates and sorts file paths.
func uniqueSorted(files []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, f := range files {
		if !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	sort.Strings(out)
	return out
}

// Run executes a docker compose command across all stack files
func (s *StackManager) Run(action string, args ...string) error {
	if len(s.Files) == 0 {
		fmt.Println("⚠️  No compose files found. Nothing to do.")
		return nil
	}

	envFile := system.GetConfigPath()
	hasEnv := false
	if _, err := os.Stat(envFile); err == nil {
		hasEnv = true
	}

	for _, file := range s.Files {
		fmt.Printf("🚀 Running docker compose %s on %s...\n", action, file)
		
		var cmdArgs []string
		if hasEnv {
			// Use --env-file for the central .env config
			cmdArgs = []string{"compose", "--env-file", envFile, "-f", file, action}
		} else {
			cmdArgs = []string{"compose", "-f", file, action}
		}
		
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.Command("docker", cmdArgs...)

		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %s on %s: %w", action, file, err)
		}
	}
	return nil
}
