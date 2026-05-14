package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Stack represents a collection of docker-compose files
type Stack struct {
	Files []string
}

// NewStack returns a stack with discovered M3TAL compose files.
// It scans for files matching patterns:
//   - ./docker/*-compose.yml
//   - ./*-stack/*-compose.yml
//
// If running as a system installation, it also scans /usr/share/m3tal/stack/.
func NewStack() *Stack {
	paths := discoverComposeFiles()
	return &Stack{
		Files: paths,
	}
}

// discoverComposeFiles finds all matching compose files across known locations.
func discoverComposeFiles() []string {
	var files []string

	// Scan local project root patterns
	localDirs := []string{
		".",
		"source",
	}

	for _, dir := range localDirs {
		matches := findComposeFiles(dir)
		files = append(files, matches...)
	}

	// If no local files found, check system paths
	if len(files) == 0 {
		if _, err := os.Stat("/usr/share/m3tal/stack"); err == nil {
			matches := findComposeFiles("/usr/share/m3tal/stack")
			files = append(files, matches...)
		}
	}

	// Deduplicate and sort for deterministic order
	files = uniqueSorted(files)

	return files
}

// findComposeFiles scans a directory for files matching the compose patterns.
func findComposeFiles(root string) []string {
	var matches []string

	// Pattern 1: ./docker/*-compose.yml
	pat1 := filepath.Join(root, "docker", "*-compose.yml")
	m1, _ := filepath.Glob(pat1)
	matches = append(matches, m1...)

	// Pattern 2: ./*-stack/*-compose.yml
	pat2 := filepath.Join(root, "*-stack", "*-compose.yml")
	m2, _ := filepath.Glob(pat2)
	matches = append(matches, m2...)

	return matches
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
func (s *Stack) Run(action string, args ...string) error {
	if len(s.Files) == 0 {
		fmt.Println("⚠️  No compose files found. Nothing to do.")
		return nil
	}

	for _, file := range s.Files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("🚀 Running docker compose %s on %s...\n", action, file)
		cmdArgs := append([]string{"compose", "-f", file, action}, args...)
		cmd := exec.Command("docker", cmdArgs...)

		// Pass current environment + env file if it exists
		cmd.Env = os.Environ()
		envFile := ".env"
		if _, err := os.Stat("/etc/m3tal/config.yaml"); err == nil {
			envFile = "/etc/m3tal/config.yaml"
		}

		if _, err := os.Stat(envFile); err == nil {
			// Filter out non-exportable keys from config.yaml
			if strings.HasSuffix(envFile, ".yaml") {
				cmdArgs = []string{"compose", "-f", file, action}
				if len(args) > 0 {
					cmdArgs = append(cmdArgs, args...)
				}
				cmd = exec.Command("docker", cmdArgs...)
				cmd.Env = append(os.Environ(), "M3TAL_CONFIG="+envFile)
			} else {
				cmdArgs = append([]string{"compose", "--env-file", envFile, "-f", file, action}, args...)
				cmd = exec.Command("docker", cmdArgs...)
			}
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %s on %s: %w", action, file, err)
		}
	}
	return nil
}
