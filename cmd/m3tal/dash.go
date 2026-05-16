package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jakej985-rgb/m3tal-core/internal/system"
	"github.com/spf13/cobra"
)

func initDashCmd() *cobra.Command {
	var dashCmd = &cobra.Command{
		Use:   "dash",
		Short: "Manage the M3TAL Dashboard",
		Long:  `Start, stop, and monitor the M3TAL dashboard container.`,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the dashboard service",
		Run: func(cmd *cobra.Command, args []string) {
			runDashCompose("up", "-d")
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop the dashboard service",
		Run: func(cmd *cobra.Command, args []string) {
			runDashCompose("stop")
		},
	}

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart the dashboard service",
		Run: func(cmd *cobra.Command, args []string) {
			runDashCompose("restart")
		},
	}

	var logsCmd = &cobra.Command{
		Use:   "logs",
		Short: "View dashboard logs",
		Run: func(cmd *cobra.Command, args []string) {
			runDashCompose("logs", "-f", "--tail=100")
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check dashboard status",
		Run: func(cmd *cobra.Command, args []string) {
			runDashCompose("ps")
		},
	}

	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Pull config, start API, and start dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			pullConfig()
			
			fmt.Println("🚀 Starting M3TAL API service...")
			if err := exec.Command("sudo", "systemctl", "start", "m3tal-api.service").Run(); err != nil {
				fmt.Println("⚠️  Could not start m3tal-api.service. It may already be running or systemd is not available.")
			}
			
			runDashCompose("pull")
			runDashCompose("up", "-d")
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull dashboard image and config",
		Run: func(cmd *cobra.Command, args []string) {
			pullConfig()
			runDashCompose("pull")
		},
	}

	var pullConfigCmd = &cobra.Command{
		Use:   "pull-config",
		Short: "Download the latest dashboard compose file from GitHub",
		Run: func(cmd *cobra.Command, args []string) {
			pullConfig()
		},
	}

	dashCmd.AddCommand(upCmd, pullCmd, pullConfigCmd, startCmd, stopCmd, restartCmd, logsCmd, statusCmd)
	return dashCmd
}

func pullConfig() {
	stackDir := system.GetStackDir()
	
	urls := map[string]string{
		"m3tal-compose.yml":         "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/m3tal-compose.yml",
		"m3tal-compose.local.yml":   "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/m3tal-compose.local.yml",
		"m3tal-compose.traefik.yml": "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/m3tal-compose.traefik.yml",
	}

	fmt.Printf("📥 Pulling latest dashboard config from GitHub...\n")

	// Pre-flight check: Can we write to this directory?
	if info, err := os.Stat(stackDir); err == nil {
		if info.Mode().Perm()&0200 == 0 {
			fmt.Printf("⚠️  Insufficient permissions to write to %s\n", stackDir)
			fmt.Println("👉 Try running: sudo m3tal dash up")
			return
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(stackDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create stack directory: %v\n", err)
		fmt.Println("👉 Try running with sudo")
		os.Exit(1)
	}

	for filename, url := range urls {
		targetFile := filepath.Join(stackDir, filename)
		cmd := exec.Command("curl", "-fsSL", "-o", targetFile, url)
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to download %s: %v\n", filename, err)
			if _, err := os.Stat(targetFile); err != nil {
				log.Fatalf("❌ No local %s found and download failed.", filename)
			}
		} else {
			fmt.Printf("✅ Saved to %s\n", targetFile)
		}
	}
}

func runDashCompose(action string, args ...string) {
	stackDir := system.GetStackDir()
	composeFile := filepath.Join(stackDir, "m3tal-compose.yml")

	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		log.Fatalf("❌ Dashboard compose file not found at %s", composeFile)
	}

	fmt.Printf("🚀 Dashboard: %s...\n", action)

	envFile := system.GetConfigPath()
	cmdArgs := []string{"compose", "-p", "m3tal"}
	
	if _, err := os.Stat(envFile); err == nil {
		cmdArgs = append(cmdArgs, "--env-file", envFile)
	}
	
	cmdArgs = append(cmdArgs, "-f", composeFile)

	mode := "local" // Default
	if data, err := os.ReadFile(envFile); err == nil {
		if val := getEnvValue(string(data), "DASHBOARD_EXPOSE_MODE"); val != "" {
			mode = val
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

	cmdArgs = append(cmdArgs, action)
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Dashboard command failed: %v", err)
	}
}
