package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	dashboardExists := false
	cmd := exec.Command("docker", "ps", "-a", "-q", "-f", "name=m3tal-dashboard")
	if out, err := cmd.Output(); err == nil && len(strings.TrimSpace(string(out))) > 0 {
		dashboardExists = true
	}

	// Scan system paths from helper (prioritizes /docker or /opt/m3tal/stack)
	stackDir := system.GetStackDir()
	matches, _ := filepath.Glob(filepath.Join(stackDir, "*-compose.yml"))
	for _, match := range matches {
		// Ignore dashboard compose so it is strictly managed by 'm3tal dash' commands
		// UNLESS it has already been started and the container exists.
		if !dashboardExists && filepath.Base(match) == "m3tal-compose.yml" {
			continue
		}
		files = append(files, match)
	}

	// Also check for files in subdirectories (legacy/extra stacks)
	subMatches, _ := filepath.Glob(filepath.Join(stackDir, "*", "*-compose.yml"))
	for _, match := range subMatches {
		if !dashboardExists && filepath.Base(match) == "m3tal-compose.yml" {
			continue
		}
		files = append(files, match)
	}

	// Deduplicate and sort for deterministic order
	files = uniqueSorted(files)

	return files
}

// uniqueSorted deduplicates and sorts file paths with priority for infrastructure stacks.
func uniqueSorted(files []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, f := range files {
		if !seen[f] {
			seen[f] = true
			unique = append(unique, f)
		}
	}

	// Priority weights (lower is earlier)
	priority := map[string]int{
		"network-compose.yml":     1,
		"routing-compose.yml":     2,
		"m3tal-compose.yml":       3,
		"maintenance-compose.yml": 4,
	}

	sort.Slice(unique, func(i, j int) bool {
		nameI := filepath.Base(unique[i])
		nameJ := filepath.Base(unique[j])

		pI, okI := priority[nameI]
		pJ, okJ := priority[nameJ]

		if okI && okJ {
			return pI < pJ
		}
		if okI {
			return true
		}
		if okJ {
			return false
		}
		return nameI < nameJ
	})

	return unique
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
			cmdArgs = []string{"compose", "-p", stackName, "--env-file", envFile, "-f", file}
		} else {
			cmdArgs = []string{"compose", "-p", stackName, "-f", file}
		}

		// Inject override configs for the dashboard
		if stackName == "m3tal" {
			mode := "local" // Default
			if data, err := os.ReadFile(envFile); err == nil {
				// Inline getEnvValue implementation since it's defined in main.go, not here
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if strings.HasPrefix(trimmed, "DASHBOARD_EXPOSE_MODE=") && !strings.HasPrefix(trimmed, "#") {
						parts := strings.SplitN(trimmed, "=", 2)
						if len(parts) == 2 {
							mode = parts[1]
							break
						}
					}
				}
			}

			if mode == "local" {
				localOverride := filepath.Join(stackDir, "m3tal-compose.local.yml")
				if _, err := os.Stat(localOverride); err == nil {
					cmdArgs = append(cmdArgs, "-f", localOverride)
				}
			} else {
				traefikOverride := filepath.Join(stackDir, "m3tal-compose.traefik.yml")
				if _, err := os.Stat(traefikOverride); err == nil {
					cmdArgs = append(cmdArgs, "-f", traefikOverride)
				}
			}
		}

		cmdArgs = append(cmdArgs, action)
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

	// NEW: Fix strict interpolation checks using Regex
	// This catches ${VAR:?Error message} and converts to ${VAR:-}
	re := regexp.MustCompile(`:\?([^}]+)}`)
	if re.MatchString(sContent) {
		sContent = re.ReplaceAllString(sContent, ":-}")
		modified = true
	}

	// NEW: Fix legacy Traefik volume mounts
	if strings.Contains(sContent, "./dynamic/api.yml") {
		// Remove the individual file mount
		sContent = strings.ReplaceAll(sContent, "      - ./dynamic/api.yml:/etc/traefik/dynamic/api.yml:ro", "")
		// Update the directory mount to the simplified local version
		sContent = strings.ReplaceAll(sContent, "      - ${CONFIG_PATH:-/mnt/config}/traefik/dynamic:/etc/traefik/dynamic:ro", "      - ./dynamic:/etc/traefik/dynamic:ro")
		modified = true
	}

	if modified {
		return os.WriteFile(filePath, []byte(sContent), 0664)
	}
	return nil
}

// RunRaw executes an arbitrary command and pipes its output to stdout/stderr.
func RunRaw(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
