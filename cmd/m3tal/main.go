package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jakej985-rgb/m3tal-core/pkg/auth"
	"github.com/jakej985-rgb/m3tal-core/pkg/containers"
	"github.com/jakej985-rgb/m3tal-core/pkg/orchestrator"
	"github.com/jakej985-rgb/m3tal-core/pkg/system"
	"github.com/spf13/cobra"
)

func main() {
	// First-run check for Linux system installations
	if runtime.GOOS == "linux" && os.Geteuid() != 0 {
		if _, err := os.Stat("/etc/m3tal"); os.IsNotExist(err) {
			// Probably a local/dev run, skip system check
		} else if _, err := os.Stat("/etc/m3tal/config.yaml"); os.IsNotExist(err) {
			fmt.Println("⚠️  System configuration not found at /etc/m3tal/config.yaml")
			fmt.Println("👉 Run: sudo m3tal init")
			fmt.Println("")
		}
	}

	var rootCmd = &cobra.Command{Use: "m3tal"}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			mgr, err := containers.GetProvider()
			if err != nil {
				log.Fatal(err)
			}
			list, err := mgr.ListContainers()
			if err != nil {
				log.Fatal(err)
			}
			printJSON(list)
		},
	}

	var startCmd = &cobra.Command{
		Use:   "start [name]",
		Short: "Start a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mgr, err := containers.GetProvider()
			if err != nil {
				log.Fatal(err)
			}
			if err := mgr.StartContainer(args[0]); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Started %s\n", args[0])
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop [name]",
		Short: "Stop a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mgr, err := containers.GetProvider()
			if err != nil {
				log.Fatal(err)
			}
			if err := mgr.StopContainer(args[0]); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Stopped %s\n", args[0])
		},
	}

	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Show system stats",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := system.GetStats()
			if err != nil {
				log.Fatal(err)
			}
			printJSON(s)
		},
	}

	var daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run M3TAL background agents",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("🚀 M3TAL Core Daemon starting...")

			// Background Metrics Collection
			go func() {
				for {
					stats, err := system.GetStats()
					if err == nil {
						data, _ := json.Marshal(stats)
						stateDir := os.Getenv("STATE_DIR")
						if stateDir != "" {
							_ = os.WriteFile(filepath.Join(stateDir, "metrics.json"), data, 0644)
						}
					}
					time.Sleep(10 * time.Second)
				}
			}()

			// Keep alive
			select {}
		},
	}

	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Initialize and start the M3TAL environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚀 Initializing M3TAL Orchestrator...")
			stack := orchestrator.NewStack()
			if err := stack.Run("up", "-d"); err != nil {
				log.Fatal(err)
			}
			fmt.Println("\n✅ M3TAL Stack is UP!")
			fmt.Println("--------------------------------------------------")
			fmt.Println("Dashboard: http://localhost:8082")
			fmt.Println("API:       http://localhost:5050")
			fmt.Println("--------------------------------------------------")
			fmt.Println("Use './m3tal logs' or 'docker compose ps' to monitor.")
		},
	}

	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Stop all M3TAL stacks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🛑 Stopping M3TAL Stacks...")
			stack := orchestrator.NewStack()
			if err := stack.Run("down"); err != nil {
				log.Fatal(err)
			}
			fmt.Println("✅ All services stopped.")
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull latest images",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📥 Pulling latest service images...")
			stack := orchestrator.NewStack()
			if err := stack.Run("pull"); err != nil {
				log.Fatal(err)
			}
			fmt.Println("✅ Images updated.")
		},
	}

	var dashpassCmd = &cobra.Command{
		Use:   "dashpass [username] [password]",
		Short: "Manage dashboard users",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			password := ""
			if len(args) > 1 {
				password = args[1]
			} else {
				fmt.Printf("Password for %s: ", username)
				fmt.Scanln(&password)
			}

			fmt.Printf("✅ Updating user %s...\n", username)
			usersFile := "source/m3tal-stack/users.json"
			if err := auth.UpdateUser(usersFile, username, password); err != nil {
				log.Fatal(err)
			}
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize environment and generate secrets",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat(".env"); err == nil {
				fmt.Println("⚠️  .env already exists. Use 'm3tal config wizard' to update.")
				return
			}
			runWizard(false)
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage M3TAL environment variables",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				runWizard(true)
				return
			}
		},
	}

	var configWizardCmd = &cobra.Command{
		Use:   "wizard",
		Short: "Run the interactive configuration wizard",
		Run: func(cmd *cobra.Command, args []string) {
			runWizard(true)
		},
	}

	var configListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all configuration variables",
		Run: func(cmd *cobra.Command, args []string) {
			content, err := os.ReadFile(".env")
			if err != nil {
				log.Fatal("❌ .env not found. Run 'init' first.")
			}
			fmt.Println(string(content))
		},
	}

	var configSetCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration variable",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			val := args[1]

			content, err := os.ReadFile(".env")
			if err != nil {
				log.Fatal("❌ .env not found. Run 'init' first.")
			}

			newContent := replaceSecret(string(content), key+"=", val)
			// If not found, append
			if !strings.Contains(newContent, key+"=") {
				newContent += fmt.Sprintf("\n%s=%s", key, val)
			}

			if err := os.WriteFile(".env", []byte(newContent), 0600); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("✅ Config updated: %s=%s\n", key, val)
		},
	}

	var configGetCmd = &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration variable",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			content, err := os.ReadFile(".env")
			if err != nil {
				log.Fatal("❌ .env not found. Run 'init' first.")
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, key+"=") {
					fmt.Println(strings.TrimPrefix(line, key+"="))
					return
				}
			}
			fmt.Printf("❌ Key '%s' not found.\n", key)
		},
	}

	configCmd.AddCommand(configListCmd, configSetCmd, configGetCmd, configWizardCmd)
	rootCmd.AddCommand(listCmd, startCmd, stopCmd, statsCmd, daemonCmd, upCmd, downCmd, pullCmd, dashpassCmd, initCmd, configCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runWizard(update bool) {
	fmt.Println("🛠️  M3TAL Configuration Wizard")

	targetFile := ".env"
	isSystem := false
	if runtime.GOOS == "linux" && os.Geteuid() == 0 {
		isSystem = true
		targetFile = "/etc/m3tal/config.yaml"
		_ = os.MkdirAll("/etc/m3tal", 0755)
		_ = os.MkdirAll("/var/lib/m3tal", 0755)
		fmt.Println("✅ System directories initialized (/etc/m3tal, /var/lib/m3tal)")
	}

	_ = os.MkdirAll("./data", 0755)

	sourceFile := ".env.example"
	if update {
		if _, err := os.Stat(targetFile); err == nil {
			sourceFile = targetFile
		} else if _, err := os.Stat(".env"); err == nil {
			sourceFile = ".env"
		}
	}

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		log.Fatalf("❌ Missing %s file.", sourceFile)
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			key := parts[0]
			val := parts[1]

			if (key == "DASHBOARD_SECRET" || key == "API_TOKEN") && !update {
				newSecret := generateSecret()
				fmt.Printf("[Auto] %s generated: %s\n", key, newSecret)
				lines[i] = key + "=" + newSecret
				continue
			}

			fmt.Printf("%s [%s]: ", key, val)
			var input string
			fmt.Scanln(&input)
			if input != "" {
				lines[i] = key + "=" + input
			}
		}
	}

	content := strings.Join(lines, "\n")
	if err := os.WriteFile(targetFile, []byte(content), 0600); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n✅ Configuration saved to %s\n", targetFile)
	if isSystem {
		fmt.Println("👉 You may also want to symlink this to .env for local stack commands.")
	}
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func generateSecret() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func replaceSecret(content, key, secret string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, key) {
			lines[i] = key + secret
		}
	}
	return strings.Join(lines, "\n")
}
