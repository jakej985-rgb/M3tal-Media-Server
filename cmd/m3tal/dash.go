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
		Short: "Pull config, image, and start dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			pullConfig()
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
	composeFile := filepath.Join(stackDir, "m3tal-compose.yml")
	url := "https://raw.githubusercontent.com/jakej985-rgb/m3tal-core/main/deploy/stack/m3tal-compose.yml"

	fmt.Printf("📥 Pulling latest dashboard config from GitHub...\n")

	// Ensure directory exists
	if err := os.MkdirAll(stackDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create stack directory: %v", err)
	}

	// Use curl for simplicity (handles redirects, etc.)
	cmd := exec.Command("curl", "-fsSL", "-o", composeFile, url)
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to download compose file: %v", err)
	}
	fmt.Printf("✅ Saved to %s\n", composeFile)
}

func runDashCompose(action string, args ...string) {
	stackDir := system.GetStackDir()
	composeFile := filepath.Join(stackDir, "m3tal-compose.yml")

	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		log.Fatalf("❌ Dashboard compose file not found at %s", composeFile)
	}

	fmt.Printf("🚀 Dashboard: %s...\n", action)

	envFile := system.GetConfigPath()
	cmdArgs := []string{"compose"}
	
	if _, err := os.Stat(envFile); err == nil {
		cmdArgs = append(cmdArgs, "--env-file", envFile)
	}
	
	cmdArgs = append(cmdArgs, "-f", composeFile, action, "m3tal-dashboard")
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Dashboard command failed: %v", err)
	}
}
