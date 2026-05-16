package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

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

	// Auto-create required external networks
	ensureNetworkExists("proxy")
	ensureNetworkExists("m3tal")

	envFile := system.GetConfigPath()
	hasEnv := false
	if _, err := os.Stat(envFile); err == nil {
		hasEnv = true
	}

	for _, file := range s.Files {
		// 1. Determine stack name and ensure <stack>.env exists
		stackBase := filepath.Base(file)
		stackName := stackBase[:len(stackBase)-len("-compose.yml")]
		stackDir := filepath.Dir(file)
		localEnv := filepath.Join(stackDir, stackName+".env")

		// Automatically symlink global env to <stack>.env if it doesn't exist
		if _, err := os.Lstat(localEnv); os.IsNotExist(err) && hasEnv {
			_ = os.Symlink(envFile, localEnv)
		}

		// Self-healing: Fix legacy .env paths in the manifest
		_ = fixLegacyManifest(file, stackName)

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

// fixLegacyManifest scans a compose file for old relative .env paths and repairs them.
func fixLegacyManifest(filePath, stackName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// We look for common legacy patterns and replace with the local ./<stack>.env
	brokenPaths := []string{"../../.env", "../m3tal-core/.env"}
	newPath := fmt.Sprintf("./%s.env", stackName)
	
	modified := false
	sContent := string(content)
	for _, old := range brokenPaths {
		if strings.Contains(sContent, old) {
			sContent = strings.ReplaceAll(sContent, old, newPath)
			modified = true
		}
	}

	// Remove legacy docker-proxy service block if it exists
	if strings.Contains(sContent, "docker-proxy:") {
		// Simple block removal for docker-proxy
		startIdx := strings.Index(sContent, "  docker-proxy:")
		if startIdx != -1 {
			// Find the start of the next top-level block (networks:)
			endIdx := strings.Index(sContent[startIdx:], "\nnetworks:")
			if endIdx != -1 {
				sContent = sContent[:startIdx] + sContent[startIdx+endIdx+1:]
				modified = true
			}
		}
	}

	if modified {
		return os.WriteFile(filePath, []byte(sContent), 0664)
	}
	return nil
}

// ensureNetworkExists checks if a docker network exists and creates it if not.
func ensureNetworkExists(name string) {
	// check if exists
	checkCmd := exec.Command("docker", "network", "inspect", name)
	if err := checkCmd.Run(); err != nil {
		// Doesn't exist, create it
		fmt.Printf("🌐 Creating required external network: %s\n", name)
		createCmd := exec.Command("docker", "network", "create", name)
		_ = createCmd.Run()
	}
}
