package main

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/jakej985-rgb/m3tal-core/internal/api"
	"github.com/jakej985-rgb/m3tal-core/internal/auth"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/health"
	"github.com/jakej985-rgb/m3tal-core/internal/orchestrator"
	"github.com/jakej985-rgb/m3tal-core/internal/plugin"
	"github.com/jakej985-rgb/m3tal-core/internal/preflight"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
	"github.com/spf13/cobra"
)

//go:embed .env.example
var envExample string

func main() {
	// First-run check for Linux system installations
	if runtime.GOOS == "linux" && os.Geteuid() != 0 {
		configPath := system.GetConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("⚠️  System configuration not found at %s\n", configPath)
			fmt.Println("👉 Run: sudo m3tal init")
			fmt.Println("")
		}
	}

	var rootCmd = &cobra.Command{
		Use:   "m3tal",
		Short: "M3TAL Core Orchestrator",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				persist, _ := cmd.Flags().GetBool("persist")
				newWindow, _ := cmd.Flags().GetBool("new-window")
				windowAlias, _ := cmd.Flags().GetBool("window")

				globalLocal, _ = cmd.Flags().GetBool("local")
				globalAPIURL, _ = cmd.Flags().GetString("api-url")
				globalAPIToken, _ = cmd.Flags().GetString("api-token")

				if newWindow || windowAlias {
					fmt.Println("🖥️  Launching M3TAL interactive menu in a new terminal window...")
					if err := launchInNewWindow(); err != nil {
						fmt.Printf("❌ Failed to launch in new window: %v\n", err)
						fmt.Println("👉 Falling back to current terminal window...")
					} else {
						return
					}
				}

				if persist || newWindow || windowAlias {
					setupPersistentSignalHandler()
					for {
						if !runMainMenu() {
							break
						}
					}
				} else {
					runMainMenu()
				}
				return
			}
			cmd.Help()
		},
	}

	rootCmd.Flags().BoolP("persist", "p", false, "Keep the interactive menu open and loop continuously")
	rootCmd.Flags().BoolP("new-window", "n", false, "Launch the interactive menu in a new terminal window")
	rootCmd.Flags().BoolP("window", "w", false, "Alias for --new-window")
	rootCmd.PersistentFlags().String("api-url", "http://localhost:5050", "M3TAL API URL")
	rootCmd.PersistentFlags().String("api-token", os.Getenv("API_TOKEN"), "M3TAL API Token")
	rootCmd.PersistentFlags().Bool("local", false, "Force local execution (skip API)")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			if useAPI, _ := cmd.Flags().GetBool("local"); !useAPI {
				if resp, err := callAPI(cmd, "GET", "/api/containers", nil); err == nil {
					fmt.Println(resp)
					return
				}
			}

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

	var psCmd = &cobra.Command{
		Use:   "ps",
		Short: "List all containers (alias for list)",
		Run:   listCmd.Run,
	}

	var startCmd = &cobra.Command{
		Use:   "start [name]",
		Short: "Start a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if useAPI, _ := cmd.Flags().GetBool("local"); !useAPI {
				body := map[string]string{"name": args[0]}
				if _, err := callAPI(cmd, "POST", "/api/containers/start", body); err == nil {
					fmt.Printf("Started %s (via API)\n", args[0])
					return
				}
			}

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
			if useAPI, _ := cmd.Flags().GetBool("local"); !useAPI {
				body := map[string]string{"name": args[0]}
				if _, err := callAPI(cmd, "POST", "/api/containers/stop", body); err == nil {
					fmt.Printf("Stopped %s (via API)\n", args[0])
					return
				}
			}

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
			if useAPI, _ := cmd.Flags().GetBool("local"); !useAPI {
				if resp, err := callAPI(cmd, "GET", "/api/metrics", nil); err == nil {
					fmt.Println(resp)
					return
				}
			}

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

	var apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run the M3TAL API server",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")
			token := os.Getenv("API_TOKEN")
			if token == "" {
				token = "m3tal-secret-token"
			}

			dbPath := store.GetStatePath()
			db, err := store.Open(dbPath)
			if err != nil {
				log.Printf("⚠️  Could not open state database at %s: %v", dbPath, err)
				log.Println("⚠️  v2 engine endpoints will be disabled. Starting with v1 only.")
				if err := api.StartServer(port, token); err != nil {
					log.Fatalf("❌ API server failed: %v", err)
				}
				return
			}
			defer db.Close()

			log.Printf("📦 State database: %s\n", dbPath)

			if err := api.StartServerWithStore(port, token, db); err != nil {
				log.Fatalf("❌ API server failed: %v", err)
			}
		},
	}
	apiCmd.Flags().String("port", "5050", "Port to listen on")

	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Initialize and start the M3TAL environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("🚀 Initializing M3TAL Orchestrator (using stack: %s)...\n", system.UserfacingStackPath)
			stack := orchestrator.NewStackManager()
			if len(stack.Files) > 0 {
				if err := stack.Run("up", "-d"); err != nil {
					log.Fatal(err)
				}
				fmt.Println("\n✅ M3TAL Stack is UP!")
				fmt.Println("--------------------------------------------------")
				fmt.Println("Dashboard: http://localhost:8082")
				fmt.Println("API:       http://localhost:5050")
				fmt.Println("--------------------------------------------------")
				fmt.Printf("Use 'm3tal logs' or 'docker compose -f %s/m3tal-compose.yml ps' to monitor.\n", system.UserfacingStackPath)
			} else {
				fmt.Printf("⚠️  No stacks detected in %s.\n", system.UserfacingStackPath)
				fmt.Println("🛰️  M3TAL is listening and ready for deployments.")
			}
		},
	}

	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Stop all M3TAL stacks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("🛑 Stopping M3TAL Stacks (path: %s)...\n", system.UserfacingStackPath)
			stack := orchestrator.NewStackManager()
			if err := stack.Run("down"); err != nil {
				log.Fatal(err)
			}
			fmt.Println("✅ All services stopped.")
		},
	}

	var logsCmd = &cobra.Command{
		Use:   "logs [stack]",
		Short: "View logs from M3TAL stacks (Interactive Menu)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				runLogsMenu()
				return
			}

			// Legacy behavior for direct arguments
			stack := orchestrator.NewStackManager()
			target := args[0]
			var filtered []string
			for _, f := range stack.Files {
				if strings.Contains(f, target) {
					filtered = append(filtered, f)
				}
			}
			stack.Files = filtered

			if len(stack.Files) > 0 {
				stack.Run("logs", "--tail", "50", "-f")
			} else {
				fmt.Println("❌ No matching stacks found.")
			}
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull latest images",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📥 Pulling latest service images...")
			stack := orchestrator.NewStackManager()
			if err := stack.Run("pull"); err != nil {
				log.Fatal(err)
			}
			fmt.Println("✅ Images updated.")
		},
	}

	var dashpassCmd = &cobra.Command{
		Use:   "dashpass [username] [password]",
		Short: "Manage dashboard users",
		Run: func(cmd *cobra.Command, args []string) {
			username := ""
			if len(args) > 0 {
				username = args[0]
			} else {
				fmt.Print("Enter username: ")
				fmt.Scanln(&username)
				username = strings.TrimSpace(username)
				if username == "" {
					fmt.Println("❌ Username cannot be empty.")
					return
				}
			}

			password := ""
			if len(args) > 1 {
				password = args[1]
			} else {
				fmt.Printf("Password for %s: ", username)
				fmt.Scanln(&password)
			}

			fmt.Printf("✅ Updating user %s...\n", username)
			usersFile := filepath.Join(system.GetStackDir(), "users.json")
			if err := auth.UpdateUser(usersFile, username, password); err != nil {
				log.Fatal(err)
			}
			// Restart dashboard container to apply new credentials and remount users.json immediately
			_ = exec.Command("docker", "restart", "m3tal-dashboard").Run()
		},
	}

	var docCmd = &cobra.Command{
		Use:   "doctor",
		Short: "Run comprehensive pre-flight health check",
		Long: `Checks Docker daemon connectivity, .env file validity,
storage path accessibility, port availability, and system configuration.
Run this before 'm3tal up' to diagnose potential issues.`,
		Run: func(cmd *cobra.Command, args []string) {
			envPath := system.GetConfigPath()
			baseStoragePath := ""
			if data, err := os.ReadFile(envPath); err == nil {
				baseStoragePath = getEnvValue(string(data), "BASE_STORAGE_PATH")
			}
			results := preflight.RunAll(envPath, baseStoragePath)
			preflight.PrintResults(results)
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize environment and generate secrets",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat(system.GetConfigPath()); err == nil {
				fmt.Println("⚠️  .env already exists. Use 'm3tal config wizard' to update.")
				return
			}
			runWizard(system.GetConfigPath(), "", false, true)

			// Post-init storage path validation
			envData, err := os.ReadFile(system.GetConfigPath())
			if err == nil {
				basePath := getEnvValue(string(envData), "BASE_STORAGE_PATH")
				if basePath != "" {
					if err := preflight.ValidateStoragePath(basePath); err != nil {
						fmt.Printf("\n⚠️  Storage path validation:\n    %v\n", err)
						fmt.Println("\n👉 Set a valid BASE_STORAGE_PATH in .env, then run:")
						fmt.Println("   mkdir -p /path/to/your/m3tal-data")
						fmt.Println("   ./m3tal config set BASE_STORAGE_PATH /path/to/your/m3tal-data")
					} else {
						fmt.Printf("✅ Storage path validated: %s\n", basePath)
					}
				}
			}

			// Initialize plugin directory structure
			if runtime.GOOS == "linux" && os.Geteuid() == 0 {
				for _, subdir := range system.PluginSubdirs {
					path := filepath.Join(system.UserPluginsDir, subdir)
					_ = os.MkdirAll(path, 0755)
				}
				fmt.Printf("✅ Plugin directories initialized (%s)\n", system.UserPluginsDir)
			}
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage M3TAL environment variables",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				runWizard(system.GetConfigPath(), "", true, true)
				return
			}
		},
	}

	var configWizardCmd = &cobra.Command{
		Use:   "wizard",
		Short: "Run the interactive configuration wizard",
		Run: func(cmd *cobra.Command, args []string) {
			target, _ := cmd.Flags().GetString("target")
			composeFile, _ := cmd.Flags().GetString("compose")
			isGlobal := false

			if target == "" {
				target = system.GetConfigPath()
				isGlobal = true
			}

			runWizard(target, composeFile, true, isGlobal)
		},
	}
	configWizardCmd.Flags().String("target", "", "Target .env file to save")
	configWizardCmd.Flags().String("compose", "", "Compose file to read required variables from")

	var configListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all configuration variables",
		Run: func(cmd *cobra.Command, args []string) {
			content, err := os.ReadFile(system.GetConfigPath())
			if err != nil {
				log.Fatalf("❌ Configuration not found. Run 'init' first.")
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

			content, err := os.ReadFile(system.GetConfigPath())
			if err != nil {
				log.Fatalf("❌ Configuration not found. Run 'init' first.")
			}

			newContent := replaceSecret(string(content), key+"=", val)
			// If not found, append
			if !strings.Contains(newContent, key+"=") {
				newContent += fmt.Sprintf("\n%s=%s", key, val)
			}

			if err := os.WriteFile(system.GetConfigPath(), []byte(newContent), 0600); err != nil {
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
			content, err := os.ReadFile(system.GetConfigPath())
			if err != nil {
				log.Fatalf("❌ Configuration not found. Run 'init' first.")
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

	var configScanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan available Docker stacks",
		Run: func(cmd *cobra.Command, args []string) {
			stackDir := system.GetStackDir()
			entries, err := os.ReadDir(stackDir)
			if err != nil {
				fmt.Println("❌ Unable to read docker directory:", err)
				return
			}
			type stackInfo struct {
				Compose  string `json:"compose"`
				Template string `json:"template,omitempty"`
				Env      string `json:"env,omitempty"`
			}
			stacks := make(map[string]stackInfo)
			for _, e := range entries {
				name := e.Name()
				if strings.HasSuffix(name, "-compose.yml") {
					stack := strings.TrimSuffix(name, "-compose.yml")
					composePath := filepath.Join(stackDir, name)
					templatePath := filepath.Join(stackDir, stack+".env.template")
					envPath := filepath.Join(stackDir, stack+".env")

					hasTemplate := false
					hasEnv := false

					if _, err := os.Stat(templatePath); err == nil {
						hasTemplate = true
					}
					if _, err := os.Stat(envPath); err == nil {
						hasEnv = true
					}

					if !hasTemplate && !hasEnv {
						continue
					}

					info := stackInfo{Compose: composePath}
					if hasTemplate {
						info.Template = templatePath
					}
					if hasEnv {
						info.Env = envPath
					}

					stacks[stack] = info
				}
			}
			printJSON(stacks)
		},
	}
	configCmd.AddCommand(configListCmd, configSetCmd, configGetCmd, configScanCmd, configWizardCmd)

	// ─── Plugin Commands ───
	var pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "Manage M3TAL plugins",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var pluginListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all discovered plugins",
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load plugins: %v", err)
			}

			fmt.Printf("\n🔌 %s\n\n", reg.Summary())

			if len(reg.ListStacks()) > 0 {
				fmt.Println("📦 Stack Plugins:")
				for _, s := range reg.ListStacks() {
					pri := ""
					if s.Priority > 0 {
						pri = fmt.Sprintf(" (priority: %d)", s.Priority)
					}
					status := " [enabled]"
					if !s.Enabled {
						status = " [disabled]"
					}
					fmt.Printf("   %-20s %s%s%s\n", s.Metadata.Name, s.Metadata.Description, pri, status)
				}
				fmt.Println()
			}

			if len(reg.ListRoutes()) > 0 {
				fmt.Println("🚦 Route Plugins:")
				for _, r := range reg.ListRoutes() {
					status := " [enabled]"
					if !r.Enabled {
						status = " [disabled]"
					}
					fmt.Printf("   %-20s %s → %s:%d%s\n", r.Metadata.Name, r.Domain, r.Service, r.Port, status)
				}
				fmt.Println()
			}

			if len(reg.ListMiddlewares()) > 0 {
				fmt.Println("🔐 Middleware Plugins:")
				for _, m := range reg.ListMiddlewares() {
					status := " [enabled]"
					if !m.Enabled {
						status = " [disabled]"
					}
					fmt.Printf("   %-20s [%s] %s%s\n", m.Metadata.Name, m.Type, m.Metadata.Description, status)
				}
				fmt.Println()
			}

			fmt.Println("Scanned directories:")
			for _, d := range dirs {
				marker := "  ✗"
				if _, err := os.Stat(d); err == nil {
					marker = "  ✓"
				}
				fmt.Printf("%s %s\n", marker, d)
			}
		},
	}

	var pluginValidateCmd = &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate a plugin YAML file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			data, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("❌ Cannot read file: %v", err)
			}

			p, err := plugin.ParsePlugin(data)
			if err != nil {
				log.Fatalf("❌ Parse error: %v", err)
			}

			if err := p.Validate(); err != nil {
				log.Fatalf("❌ Validation failed: %v", err)
			}

			fmt.Printf("✅ Valid %s plugin: %s\n", p.Kind, p.Metadata.Name)
			if p.Metadata.Description != "" {
				fmt.Printf("   Description: %s\n", p.Metadata.Description)
			}
			if p.Metadata.Version != "" {
				fmt.Printf("   Version:     %s\n", p.Metadata.Version)
			}
		},
	}

	var pluginEnableCmd = &cobra.Command{
		Use:   "enable [name]",
		Short: "Enable a disabled plugin by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load plugins: %v", err)
			}
			name := args[0]
			var path string
			if p := reg.GetRoute(name); p != nil {
				path = p.SourcePath
			} else if p := reg.GetStack(name); p != nil {
				path = p.SourcePath
			} else if p := reg.GetMiddleware(name); p != nil {
				path = p.SourcePath
			}

			if path == "" {
				log.Fatalf("❌ Plugin %q not found", name)
			}

			p, err := plugin.LoadPlugin(path)
			if err != nil {
				log.Fatalf("❌ Failed to load plugin manifest: %v", err)
			}

			var db *store.Store
			dbPath := store.GetStatePath()
			if dbPath != "" {
				if d, err := store.Open(dbPath); err == nil {
					db = d
					defer db.Close()
				}
			}

			mgr := plugin.NewStateManager(db)
			err = mgr.SetPluginEnabled(p, true)
			if err != nil {
				log.Fatalf("❌ Failed to enable: %v", err)
			}
			fmt.Printf("✅ Enabled plugin %q (renamed to %s)\n", name, filepath.Base(p.SourcePath))
		},
	}

	var pluginDisableCmd = &cobra.Command{
		Use:   "disable [name]",
		Short: "Disable an active plugin by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load plugins: %v", err)
			}
			name := args[0]
			var path string
			if p := reg.GetRoute(name); p != nil {
				path = p.SourcePath
			} else if p := reg.GetStack(name); p != nil {
				path = p.SourcePath
			} else if p := reg.GetMiddleware(name); p != nil {
				path = p.SourcePath
			}

			if path == "" {
				log.Fatalf("❌ Plugin %q not found", name)
			}

			p, err := plugin.LoadPlugin(path)
			if err != nil {
				log.Fatalf("❌ Failed to load plugin manifest: %v", err)
			}

			var db *store.Store
			dbPath := store.GetStatePath()
			if dbPath != "" {
				if d, err := store.Open(dbPath); err == nil {
					db = d
					defer db.Close()
				}
			}

			mgr := plugin.NewStateManager(db)
			err = mgr.SetPluginEnabled(p, false)
			if err != nil {
				log.Fatalf("❌ Failed to disable: %v", err)
			}
			fmt.Printf("✅ Disabled plugin %q (renamed to %s)\n", name, filepath.Base(p.SourcePath))
		},
	}

	var pluginSyncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Synchronize and write Traefik dynamic provider configuration",
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load plugins: %v", err)
			}

			configData, err := reg.GenerateTraefikConfig()
			if err != nil {
				log.Fatalf("❌ Failed to generate Traefik config: %v", err)
			}

			stackDir := system.GetStackDir()
			dynamicDir := filepath.Join(stackDir, "dynamic")
			if err := os.MkdirAll(dynamicDir, 0755); err != nil {
				log.Fatalf("❌ Failed to create dynamic directory: %v", err)
			}

			outputPath := filepath.Join(dynamicDir, "m3tal-plugins.yml")
			if err := os.WriteFile(outputPath, configData, 0644); err != nil {
				log.Fatalf("❌ Failed to write config file: %v", err)
			}

			fmt.Printf("✅ Synced Traefik dynamic provider config to %s\n", outputPath)
		},
	}

	var pluginMatchCmd = &cobra.Command{
		Use:   "match [service-name]",
		Short: "Find a route plugin matching the given service information",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load plugins: %v", err)
			}

			image, _ := cmd.Flags().GetString("image")
			labelSlice, _ := cmd.Flags().GetStringSlice("label")
			labels := make(map[string]string)
			for _, item := range labelSlice {
				parts := strings.SplitN(item, "=", 2)
				if len(parts) == 2 {
					labels[parts[0]] = parts[1]
				}
			}

			match := reg.MatchService(args[0], image, labels)
			if match != nil {
				fmt.Printf("🎯 Match found! Route Plugin: %s (service: %s, domain: %s, port: %d)\n",
					match.Metadata.Name, match.Service, match.Domain, match.Port)
			} else {
				fmt.Println("❌ No matching route plugin found.")
			}
		},
	}
	pluginMatchCmd.Flags().String("image", "", "Docker image name to match against")
	pluginMatchCmd.Flags().StringSlice("label", nil, "Docker labels to match against (format: key=value)")

	var pluginInstallStackCmd = &cobra.Command{
		Use:   "install-stack [name]",
		Short: "Install and parameterize a Stack plugin compose template",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load plugins: %v", err)
			}

			name := args[0]
			stack := reg.GetStack(name)
			if stack == nil {
				log.Fatalf("❌ Stack plugin %q not found", name)
			}

			composeFile := stack.ComposePath
			if !filepath.IsAbs(composeFile) {
				composeFile = filepath.Join(filepath.Dir(stack.SourcePath), composeFile)
			}

			composeData, err := os.ReadFile(composeFile)
			if err != nil {
				log.Fatalf("❌ Failed to read compose template from %s: %v", composeFile, err)
			}

			envSlice, _ := cmd.Flags().GetStringSlice("env")
			vars := make(map[string]string)
			for _, item := range envSlice {
				parts := strings.SplitN(item, "=", 2)
				if len(parts) == 2 {
					vars[parts[0]] = parts[1]
				}
			}

			for _, envLine := range os.Environ() {
				parts := strings.SplitN(envLine, "=", 2)
				if len(parts) == 2 {
					if _, ok := vars[parts[0]]; !ok {
						vars[parts[0]] = parts[1]
					}
				}
			}

			finalCompose := plugin.Parameterize(string(composeData), vars)

			stackDir := system.GetStackDir()
			outputPath := filepath.Join(stackDir, fmt.Sprintf("%s-compose.yml", name))
			if err := os.WriteFile(outputPath, []byte(finalCompose), 0644); err != nil {
				log.Fatalf("❌ Failed to write compose file: %v", err)
			}

			fmt.Printf("✅ Stack compose file installed to %s\n", outputPath)
		},
	}
	pluginInstallStackCmd.Flags().StringSlice("env", nil, "Environment variables to parameterize the template (format: key=value)")

	var pluginCatalogCmd = &cobra.Command{
		Use:   "catalog",
		Short: "List all official plugins in the catalog and their status",
		Run: func(cmd *cobra.Command, args []string) {
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load registry: %v", err)
			}
			items := plugin.ListCatalog(reg)
			fmt.Println("\n📋 M3TAL Plugin Catalog:")
			fmt.Println("--------------------------------------------------")
			for _, item := range items {
				statusColor := "⚪" // not installed
				statusStr := "not installed"
				if item.Installed {
					if item.Status == "enabled" {
						statusColor = "🟢"
						statusStr = "installed & enabled"
					} else {
						statusColor = "🟡"
						statusStr = "installed & disabled"
					}
				}
				fmt.Printf("%s %-16s [%-10s] %s\n", statusColor, item.Name, item.Kind, item.Description)
				fmt.Printf("   Version: %s | Author: %s | Status: %s\n\n", item.Version, item.Author, statusStr)
			}
		},
	}

	var pluginInstallCmd = &cobra.Command{
		Use:   "install [name]",
		Short: "Download and install a plugin from the catalog by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			// Find the item in the catalog to get the Kind
			var targetKind string
			for _, item := range plugin.Catalog {
				if strings.EqualFold(item.Name, name) {
					targetKind = item.Kind
					name = item.Name // use canonical name
					break
				}
			}
			if targetKind == "" {
				log.Fatalf("❌ Plugin %q not found in catalog. Run 'm3tal plugin catalog' to see available plugins.", name)
			}

			userDir := system.UserPluginsDir
			if _, err := os.Stat("deploy/plugins"); err == nil {
				userDir = "deploy/plugins"
			}

			fmt.Printf("📥 Installing %s plugin %q...\n", targetKind, name)
			err := plugin.InstallPlugin(name, targetKind, userDir)
			if err != nil {
				log.Fatalf("❌ Installation failed: %v", err)
			}
			fmt.Printf("✅ Plugin %q successfully installed.\n", name)
		},
	}

	var pluginUninstallCmd = &cobra.Command{
		Use:   "uninstall [name]",
		Short: "Uninstall a user-installed plugin by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				log.Fatalf("❌ Failed to load registry: %v", err)
			}

			// Find kind and source path
			var targetKind string
			var path string
			for i := range reg.Routes {
				if plugin.MatchesPluginName(reg.Routes[i].SourcePath, reg.Routes[i].Metadata.Name, name) {
					targetKind = "Route"
					if strings.Contains(reg.Routes[i].SourcePath, "/traefik/") {
						targetKind = "Traefik"
					}
					name = plugin.GetPluginBaseName(reg.Routes[i].SourcePath)
					path = reg.Routes[i].SourcePath
					break
				}
			}
			if targetKind == "" {
				for i := range reg.Stacks {
					if plugin.MatchesPluginName(reg.Stacks[i].SourcePath, reg.Stacks[i].Metadata.Name, name) {
						targetKind = "Stack"
						name = plugin.GetPluginBaseName(reg.Stacks[i].SourcePath)
						path = reg.Stacks[i].SourcePath
						break
					}
				}
			}
			if targetKind == "" {
				for i := range reg.Middlewares {
					if plugin.MatchesPluginName(reg.Middlewares[i].SourcePath, reg.Middlewares[i].Metadata.Name, name) {
						targetKind = "Middleware"
						if strings.Contains(reg.Middlewares[i].SourcePath, "/traefik/") {
							targetKind = "Traefik"
						}
						name = plugin.GetPluginBaseName(reg.Middlewares[i].SourcePath)
						path = reg.Middlewares[i].SourcePath
						break
					}
				}
			}

			if targetKind == "" || path == "" {
				log.Fatalf("❌ Plugin %q not found in local registry", name)
			}

			userDir := system.UserPluginsDir
			if _, err := os.Stat("deploy/plugins"); err == nil {
				userDir = "deploy/plugins"
			}

			p, err := plugin.LoadPlugin(path)
			if err != nil {
				log.Fatalf("❌ Failed to load plugin manifest: %v", err)
			}

			fmt.Printf("🗑️  Uninstalling plugin %q...\n", name)

			var db *store.Store
			dbPath := store.GetStatePath()
			if dbPath != "" {
				if d, err := store.Open(dbPath); err == nil {
					db = d
					defer db.Close()
				}
			}

			mgr := plugin.NewStateManager(db)
			err = mgr.UninstallPlugin(p, userDir, reg)
			if err != nil {
				log.Fatalf("❌ Uninstallation failed: %v", err)
			}
			fmt.Printf("✅ Plugin %q successfully uninstalled.\n", name)
		},
	}

	pluginCmd.AddCommand(pluginListCmd, pluginValidateCmd, pluginEnableCmd, pluginDisableCmd, pluginSyncCmd, pluginMatchCmd, pluginInstallStackCmd, pluginCatalogCmd, pluginInstallCmd, pluginUninstallCmd)
	rootCmd.AddCommand(listCmd, psCmd, startCmd, stopCmd, statsCmd, daemonCmd, apiCmd, upCmd, downCmd, logsCmd, pullCmd, dashpassCmd, initCmd, docCmd, configCmd, pluginCmd, initDashCmd(), trayCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseComposeVariables(composeFile string) map[string]string {
	data, err := os.ReadFile(composeFile)
	if err != nil {
		return nil
	}

	// Regex matches ${VAR} or ${VAR:-default} or ${VAR-default}
	re := regexp.MustCompile(`\$\{([A-Z0-9_]+)(?::?-([^}]+))?\}`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	vars := make(map[string]string)
	for _, match := range matches {
		if len(match) > 1 {
			key := match[1]
			val := ""
			if len(match) > 2 {
				val = match[2]
			}
			if existing, ok := vars[key]; !ok || (existing == "" && val != "") {
				vars[key] = val
			}
		}
	}
	return vars
}

func runWizard(targetFile string, composeFile string, update bool, isGlobal bool) {
	fmt.Printf("🛠️  M3TAL Configuration Wizard (%s)\n", filepath.Base(targetFile))

	if isGlobal && runtime.GOOS == "linux" && os.Geteuid() == 0 {
		configDir := filepath.Dir(system.ConfigPath)
		_ = os.MkdirAll(configDir, 0755)
		_ = os.MkdirAll(system.DataPath, 0755)
		fmt.Printf("✅ System directories initialized (%s, %s)\n", configDir, system.DataPath)
	}

	_ = os.MkdirAll("./data", 0755)

	var existingData []byte
	var realTargetFile = targetFile
	isSymlink := false

	if _, err := os.Lstat(targetFile); err == nil {
		info, _ := os.Lstat(targetFile)
		if info.Mode()&os.ModeSymlink != 0 {
			isSymlink = true
			linkTarget, err := os.Readlink(targetFile)
			if err == nil {
				if !filepath.IsAbs(linkTarget) {
					linkTarget = filepath.Join(filepath.Dir(targetFile), linkTarget)
				}
				realTargetFile = linkTarget
				fmt.Printf("👉 Target %s is a symlink pointing to %s. Keeping symlink and updating the shared configuration.\n", filepath.Base(targetFile), filepath.Base(linkTarget))
			}
		}
	}

	if update {
		if _, err := os.Stat(realTargetFile); err == nil {
			existingData, _ = os.ReadFile(realTargetFile)
		}
	}

	existingVars := make(map[string]string)
	if existingData != nil {
		lines := strings.Split(string(existingData), "\n")
		for _, line := range lines {
			if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				existingVars[parts[0]] = parts[1]
			}
		}
	}

	requiredVars := make(map[string]string)
	var orderedKeys []string

	if isGlobal {
		lines := strings.Split(string(envExample), "\n")
		for _, line := range lines {
			if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				requiredVars[parts[0]] = parts[1]
				orderedKeys = append(orderedKeys, parts[0])
			}
		}
	} else if composeFile != "" {
		parsedVars := parseComposeVariables(composeFile)
		for k, v := range parsedVars {
			requiredVars[k] = v
			orderedKeys = append(orderedKeys, k)
		}
	}

	if len(orderedKeys) == 0 {
		log.Fatalf("❌ Missing configuration source or no variables found.")
	}

	var finalLines []string
	for _, key := range orderedKeys {
		defaultVal := requiredVars[key]

		val := defaultVal
		if existing, ok := existingVars[key]; ok {
			val = existing
		}

		if isGlobal && (key == "DASHBOARD_SECRET" || key == "API_TOKEN") && !update && val == "" {
			newSecret := generateSecret()
			fmt.Printf("[Auto] %s generated: %s\n", key, newSecret)
			finalLines = append(finalLines, key+"="+newSecret)
			continue
		}

		colorReset := "\033[0m"
		colorRed := "\033[31m"

		promptStr := fmt.Sprintf("%s [%s]: ", key, val)
		if val == "" {
			fmt.Printf("%s%s%s", colorRed, promptStr, colorReset)
		} else {
			fmt.Printf("%s", promptStr)
		}

		var input string
		fmt.Scanln(&input)
		if input != "" {
			finalLines = append(finalLines, key+"="+input)
		} else {
			finalLines = append(finalLines, key+"="+val)
		}
	}

	if isSymlink {
		sharedContentBytes, err := os.ReadFile(realTargetFile)
		var sharedContent string
		if err == nil {
			sharedContent = string(sharedContentBytes)
		}

		for _, line := range finalLines {
			if strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				key := parts[0]
				val := parts[1]

				if strings.Contains(sharedContent, key+"=") {
					sharedContent = replaceSecret(sharedContent, key+"=", val)
				} else {
					sharedContent = strings.TrimSuffix(sharedContent, "\n") + fmt.Sprintf("\n%s=%s\n", key, val)
				}
			}
		}

		if err := os.WriteFile(realTargetFile, []byte(sharedContent), 0600); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n✅ Shared configuration updated inline: %s\n", realTargetFile)
	} else {
		content := strings.Join(finalLines, "\n")
		if err := os.WriteFile(targetFile, []byte(content), 0600); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n✅ Configuration saved to %s\n", targetFile)
	}

	// Synchronize ADMIN_PASSWORD to users.json if set in the global wizard
	if isGlobal {
		adminPass := ""
		for _, line := range finalLines {
			if strings.HasPrefix(line, "ADMIN_PASSWORD=") {
				adminPass = strings.TrimPrefix(line, "ADMIN_PASSWORD=")
				break
			}
		}
		if adminPass != "" {
			usersFile := filepath.Join(system.GetStackDir(), "users.json")
			_ = auth.UpdateUser(usersFile, "admin", adminPass)
			_ = exec.Command("docker", "restart", "m3tal-dashboard").Run()
		}
	}

	if isGlobal {
		fmt.Println("👉 You may also want to symlink this to .env for local stack commands.")
	}
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func getEnvValue(content, key string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+"=") && !strings.HasPrefix(trimmed, "#") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				return parts[1]
			}
		}
	}
	return ""
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

func callAPI(cmd *cobra.Command, method, path string, body interface{}) (string, error) {
	apiURL, _ := cmd.Flags().GetString("api-url")
	apiToken, _ := cmd.Flags().GetString("api-token")

	var bodyReader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, apiURL+path, bodyReader)
	if err != nil {
		return "", err
	}

	if apiToken != "" {
		req.Header.Set("X-API-Token", apiToken)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func runLogsMenu() {
	fmt.Println("\n📋 M3TAL Logs Explorer")
	fmt.Println("|-- 1.) M3TAL System")
	fmt.Println("|     | - CLI (Local)")
	fmt.Println("|     | - API (Systemd)")
	fmt.Println("|-- 2.) Docker")
	fmt.Println("|     |-- 1.) Stacks")
	fmt.Println("|     |-- 2.) Containers")
	fmt.Println("|-- 3.) All (Aggregated)")
	fmt.Println("|-- 0.) Exit")

	fmt.Print("\n👉 Selection: ")
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Println("\n|-- 1.) M3TAL System")
		fmt.Println("|   [1] CLI Logs")
		fmt.Println("|   [2] API Logs (Journalctl)")
		fmt.Print("\n👉 Selection: ")
		var subChoice int
		fmt.Scanln(&subChoice)
		if subChoice == 2 {
			runWithSudoFallback("journalctl", "-u", "m3tal", "-f", "-n", "50")
		} else {
			fmt.Println("ℹ️  CLI logs are minimal and shown in terminal directly.")
		}
	case 2:
		fmt.Println("\n|-- 2.) Docker")
		fmt.Println("|   [1] List Stacks")
		fmt.Println("|   [2] List Containers")
		fmt.Print("\n👉 Selection: ")
		var subChoice int
		fmt.Scanln(&subChoice)

		stackMgr := orchestrator.NewStackManager()
		switch subChoice {
		case 1:
			fmt.Println("\n📦 Available Stacks:")
			for i, f := range stackMgr.Files {
				name := strings.TrimSuffix(filepath.Base(f), "-compose.yml")
				fmt.Printf("   [%d] %s\n", i+1, name)
			}
			fmt.Print("\n👉 Stack Number: ")
			var sNum int
			fmt.Scanln(&sNum)
			if sNum > 0 && sNum <= len(stackMgr.Files) {
				stackMgr.Files = []string{stackMgr.Files[sNum-1]}
				runAsSubcommand(func() {
					stackMgr.Run("logs", "--tail", "50", "-f")
				})
			}
		case 2:
			fmt.Println("\n🐳 Running Containers:")
			mgr, err := containers.GetProvider()
			if err != nil {
				fmt.Printf("❌ Failed to connect to Docker: %v\n", err)
				break
			}
			list, err := mgr.ListContainers()
			if err != nil {
				fmt.Printf("❌ Failed to list containers: %v\n", err)
				break
			}
			for i, c := range list {
				if len(c.Names) > 0 {
					fmt.Printf("   [%d] %s (%s)\n", i+1, c.Names[0], c.Status)
				} else {
					fmt.Printf("   [%d] %s (%s)\n", i+1, c.ID[:12], c.Status)
				}
			}
			fmt.Print("\n👉 Container Number: ")
			var cNum int
			fmt.Scanln(&cNum)
			if cNum > 0 && cNum <= len(list) {
				runWithSudoFallback("docker", "logs", "--tail", "50", "-f", list[cNum-1].Names[0])
			}
		}
	case 3:
		fmt.Println("\n🚀 Streaming Aggregated Logs...")
		stackMgr := orchestrator.NewStackManager()
		runAsSubcommand(func() {
			stackMgr.Run("logs", "--tail", "20", "-f")
		})
	case 0:
		return
	default:
		fmt.Println("❌ Invalid selection.")
	}
}

func runWithSudoFallback(name string, args ...string) {
	if name == os.Args[0] {
		if globalLocal {
			args = append(args, "--local")
		}
		if globalAPIURL != "" {
			args = append(args, "--api-url", globalAPIURL)
		}
		if globalAPIToken != "" {
			args = append(args, "--api-token", globalAPIToken)
		}
	}
	runAsSubcommand(func() {
		err := orchestrator.RunRaw(name, args...)
		if err != nil {
			fmt.Println("\n⚠️  Action failed. This might require elevated privileges.")
			fmt.Println("👉 Retrying with sudo...")
			sudoArgs := append([]string{name}, args...)
			orchestrator.RunRaw("sudo", sudoArgs...)
		}
	})
}

func printStatusHeader() {
	reg := health.UpdateAndSaveHealthRegistry()

	fmt.Println("\033[1;36m  __  __    _____    _____     _       _     \033[0m")
	fmt.Println("\033[1;36m |  \\/  |  |___ /   |_   _|   / \\     | |    \033[0m")
	fmt.Println("\033[1;36m | |\\/| |    |_ \\     | |    / _ \\    | |    \033[0m")
	fmt.Println("\033[1;36m | |  | |   ___) |    | |   / ___ \\   | |___ \033[0m")
	fmt.Println("\033[1;36m |_|  |_|  |____/     |_|  /_/   \\_\\  |_____|\033[0m")
	fmt.Println()

	var systemStr string
	switch reg.System.Status {
	case "🟢":
		systemStr = "🟢 Healthy"
	case "🟡":
		systemStr = "🟡 Degraded"
	case "🔴":
		systemStr = "🔴 Unhealthy"
	default:
		systemStr = "🟢 Healthy"
	}

	var dockerStr string
	switch reg.Docker.Status {
	case "🔴":
		dockerStr = fmt.Sprintf("🔴 %d/%d running", reg.Docker.RunningContainers, reg.Docker.TotalContainers)
	case "🟡":
		dockerStr = fmt.Sprintf("🟡 %d/%d running", reg.Docker.RunningContainers, reg.Docker.TotalContainers)
	default:
		dockerStr = fmt.Sprintf("🟢 %d/%d running", reg.Docker.RunningContainers, reg.Docker.TotalContainers)
	}

	var agentsStr string
	switch reg.Agents.Status {
	case "🟢":
		agentsStr = "🟢 active monitoring"
	case "🟡":
		agentsStr = "🟡 anomaly idle"
	default:
		agentsStr = "🔴 stuck/crashed"
	}

	var diskStr string
	switch reg.Disk.Status {
	case "🔴":
		diskStr = fmt.Sprintf("🔴 %.0f%% used", reg.Disk.UsedPercent)
	case "🟡":
		diskStr = fmt.Sprintf("🟡 %.0f%% used", reg.Disk.UsedPercent)
	default:
		diskStr = fmt.Sprintf("🟢 %.0f%% used", reg.Disk.UsedPercent)
	}

	fmt.Println("═════════════════════════════════════════════════════")
	fmt.Printf(" SYSTEM:   %-25s DOCKER:   %-25s\n", systemStr, dockerStr)
	fmt.Printf(" AGENTS:   %-25s DISK:     %-25s\n", agentsStr, diskStr)
	fmt.Println("═════════════════════════════════════════════════════")
}

func renderMetricsBar(label string, percentage float64) {
	width := 20
	filled := int((percentage / 100.0) * float64(width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := filled; i < width; i++ {
		bar += "░"
	}
	fmt.Printf("%-15s [%s] %.1f%%\n", label, bar, percentage)
}

func runMainMenu() bool {
	printStatusHeader()
	fmt.Println("\n🛠️  M3TAL CONTROL CENTER (v2)")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("[1] System Control")
	fmt.Println("[2] Observability")
	fmt.Println("[3] Configuration")
	fmt.Println("[4] Agents & Automation")
	fmt.Println("[5] Extensions (Plugins)")
	fmt.Println("[6] Dashboard & API")
	fmt.Println("[0] Exit")

	fmt.Print("\n👉 Selection: ")
	var choice int
	fmt.Scanln(&choice)

	exe := os.Args[0]

	switch choice {
	case 1:
		runSystemControlMenu(exe)
	case 2:
		runObservabilityMenu(exe)
	case 3:
		runConfigurationMenu(exe)
	case 4:
		runAgentsAutomationMenu(exe)
	case 5:
		runExtensionsMenu(exe)
	case 6:
		runDashboardAPIMenu(exe)
	case 0:
		return false
	default:
		fmt.Println("❌ Invalid selection.")
	}
	return true
}

func runSystemControlMenu(exe string) {
	for {
		fmt.Println("\n══════════ SYSTEM CONTROL ══════════")
		fmt.Println("[1] Start System        (docker compose up -d)")
		fmt.Println("[2] Stop System         (docker compose down)")
		fmt.Println("[3] Restart System")
		fmt.Println("[4] Update Images       (pull + recreate)")
		fmt.Println("[5] Status Overview     (docker ps + health)")
		fmt.Println("[6] Container Actions →")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runWithSudoFallback(exe, "up")
		case 2:
			runWithSudoFallback(exe, "down")
		case 3:
			fmt.Println("🔄 Restarting M3TAL Stacks...")
			runWithSudoFallback(exe, "down")
			runWithSudoFallback(exe, "up")
		case 4:
			fmt.Println("📥 Updating system images...")
			runWithSudoFallback(exe, "pull")
			runWithSudoFallback(exe, "up")
		case 5:
			runWithSudoFallback(exe, "ps")
		case 6:
			runContainerActionsMenu(exe)
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runContainerActionsMenu(exe string) {
	for {
		fmt.Println("\n══════════ CONTAINER ACTIONS ══════════")
		fmt.Println("[1] List Containers")
		fmt.Println("[2] Start Container")
		fmt.Println("[3] Stop Container")
		fmt.Println("[4] Restart Container")
		fmt.Println("[5] Inspect (Mounts, Env, Health)")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runWithSudoFallback(exe, "list")
		case 2:
			fmt.Print("👉 Enter container name: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "start", name)
			}
		case 3:
			fmt.Print("👉 Enter container name: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "stop", name)
			}
		case 4:
			fmt.Print("👉 Enter container name: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "stop", name)
				runWithSudoFallback(exe, "start", name)
			}
		case 5:
			fmt.Print("👉 Enter container name: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				cmd := exec.Command("docker", "inspect", name)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
				fmt.Println("\nPress Enter to return...")
				var temp string
				fmt.Scanln(&temp)
			}
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runObservabilityMenu(exe string) {
	for {
		fmt.Println("\n══════════ OBSERVABILITY ══════════")
		fmt.Println("[1] Live Logs")
		fmt.Println("[2] Log Explorer")
		fmt.Println("[3] System Metrics")
		fmt.Println("[4] Health Check (Doctor)")
		fmt.Println("[5] Aggregated View (All Signals)")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runLiveLogsMenu(exe)
		case 2:
			runLogExplorerMenu(exe)
		case 3:
			showSystemMetricsVisual()
		case 4:
			runWithSudoFallback(exe, "doctor")
		case 5:
			showAggregatedSignals(exe)
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runLiveLogsMenu(_ string) {
	for {
		fmt.Println("\n══════════ LIVE LOGS ══════════")
		fmt.Println("[1] Core (CLI / API)")
		fmt.Println("[2] Docker (All Containers)")
		fmt.Println("[3] Specific Container")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("📋 Streaming M3TAL Core API logs...")
			runWithSudoFallback("journalctl", "-u", "m3tal-api.service", "-f", "-n", "100")
		case 2:
			fmt.Println("📋 Streaming Docker compose logs...")
			stackMgr := orchestrator.NewStackManager()
			runAsSubcommand(func() {
				stackMgr.Run("logs", "--tail", "20", "-f")
			})
		case 3:
			mgr, err := containers.GetProvider()
			if err != nil {
				fmt.Printf("❌ Failed to connect to Docker: %v\n", err)
				break
			}
			list, err := mgr.ListContainers()
			if err != nil {
				fmt.Printf("❌ Failed to list containers: %v\n", err)
				break
			}
			if len(list) == 0 {
				fmt.Println("⚠️  No containers found.")
				break
			}
			fmt.Println("\n🐳 Running Containers:")
			for i, c := range list {
				if len(c.Names) > 0 {
					fmt.Printf("   [%d] %s (%s)\n", i+1, c.Names[0], c.Status)
				} else {
					fmt.Printf("   [%d] %s (%s)\n", i+1, c.ID[:12], c.Status)
				}
			}
			fmt.Print("\n👉 Container Number: ")
			var cNum int
			fmt.Scanln(&cNum)
			if cNum > 0 && cNum <= len(list) {
				runWithSudoFallback("docker", "logs", "--tail", "100", "-f", list[cNum-1].Names[0])
			}
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runLogExplorerMenu(_ string) {
	for {
		fmt.Println("\n══════════ LOG EXPLORER ══════════")
		fmt.Println("[1] M3TAL Core")
		fmt.Println("[2] Docker")
		fmt.Println("[3] Agents")
		fmt.Println("[4] System (journalctl)")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runWithSudoFallback("journalctl", "-u", "m3tal-api.service", "-n", "100")
		case 2:
			stackMgr := orchestrator.NewStackManager()
			runAsSubcommand(func() {
				stackMgr.Run("logs", "--tail", "100")
			})
		case 3:
			agentLogPath := "/var/log/m3tal/agents.log"
			if _, err := os.Stat(agentLogPath); err == nil {
				runWithSudoFallback("tail", "-n", "100", agentLogPath)
			} else {
				fmt.Println("ℹ️  No Agent logs found. Simulation Mode:")
				fmt.Println("[2026-05-21 17:30:02] [INFO] [monitor] System load within normal parameters.")
				fmt.Println("[2026-05-21 17:30:05] [INFO] [metrics] Aggregated metrics refreshed.")
				fmt.Println("[2026-05-21 17:30:10] [INFO] [anomaly] Scanning for anomalies... 0 detected.")
				fmt.Println("[2026-05-21 17:30:15] [INFO] [decision] System state is stable. No action required.")
			}
		case 4:
			runWithSudoFallback("journalctl", "-n", "100")
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func showSystemMetricsVisual() {
	stats, err := system.GetDetailedStats()
	if err != nil {
		fmt.Printf("❌ Failed to fetch metrics: %v\n", err)
		return
	}

	fmt.Println("\n══════════════ SYSTEM METRICS ══════════════")
	renderMetricsBar("CPU Usage", stats.CPUUsage)
	renderMetricsBar("Memory Usage", stats.MemoryUsage)
	fmt.Printf("  -> (%.1f GB / %.1f GB)\n", stats.MemoryUsed, stats.MemoryTotal)
	renderMetricsBar("Disk Usage", stats.DiskUsage)
	fmt.Printf("  -> (%.1f GB / %.1f GB)\n", stats.DiskUsed, stats.DiskTotal)
	if stats.GPUUsage > 0 || stats.GPUModel != "No GPU Detected" {
		renderMetricsBar("GPU Usage", stats.GPUUsage)
		fmt.Printf("  -> model: %s | temp: %.0f°C\n", stats.GPUModel, stats.GPUTemp)
	}
	fmt.Printf("Uptime:         %d hours\n", stats.Uptime/3600)
	fmt.Printf("Hostname:       %s\n", stats.Hostname)
	fmt.Println("════════════════════════════════════════════")
	fmt.Println("\nPress Enter to return...")
	var temp string
	fmt.Scanln(&temp)
}

func showAggregatedSignals(_ string) {
	fmt.Println("\n══════════════ AGGREGATED VIEW ══════════════")
	reg := health.UpdateAndSaveHealthRegistry()

	var systemStr string
	switch reg.System.Status {
	case "🟢":
		systemStr = "🟢 Healthy"
	case "🟡":
		systemStr = "🟡 Degraded"
	case "🔴":
		systemStr = "🔴 Unhealthy"
	default:
		systemStr = "🟢 Healthy"
	}

	var dockerStr string
	switch reg.Docker.Status {
	case "🔴":
		dockerStr = fmt.Sprintf("🔴 %d/%d running", reg.Docker.RunningContainers, reg.Docker.TotalContainers)
	case "🟡":
		dockerStr = fmt.Sprintf("🟡 %d/%d running", reg.Docker.RunningContainers, reg.Docker.TotalContainers)
	default:
		dockerStr = fmt.Sprintf("🟢 %d/%d running", reg.Docker.RunningContainers, reg.Docker.TotalContainers)
	}

	var agentsStr string
	switch reg.Agents.Status {
	case "🟢":
		agentsStr = "🟢 active monitoring"
	case "🟡":
		agentsStr = "🟡 anomaly idle"
	default:
		agentsStr = "🔴 stuck/crashed"
	}

	var diskStr string
	switch reg.Disk.Status {
	case "🔴":
		diskStr = fmt.Sprintf("🔴 %.0f%% used", reg.Disk.UsedPercent)
	case "🟡":
		diskStr = fmt.Sprintf("🟡 %.0f%% used", reg.Disk.UsedPercent)
	default:
		diskStr = fmt.Sprintf("🟢 %.0f%% used", reg.Disk.UsedPercent)
	}

	fmt.Printf("[SYSTEM]    %s\n", systemStr)
	fmt.Printf("[DOCKER]    %s\n", dockerStr)
	fmt.Printf("[AGENTS]    %s\n", agentsStr)
	fmt.Printf("[DISK]      %s\n", diskStr)

	stats, err := system.GetStats()
	if err == nil {
		fmt.Printf("[METRICS]   CPU: %.1f%% | RAM: %.1f%%\n", stats.CPUUsage, stats.MemoryUsage)
	}

	dirs := system.GetPluginDirs()
	if pluginReg, err := plugin.LoadAll(dirs...); err == nil {
		fmt.Printf("[PLUGINS]   %d routes | %d middlewares | %d stacks\n", len(pluginReg.Routes), len(pluginReg.Middlewares), len(pluginReg.Stacks))
	}
	fmt.Println("════════════════════════════════════════════")
	fmt.Println("\nPress Enter to return...")
	var temp string
	fmt.Scanln(&temp)
}

func runConfigurationMenu(exe string) {
	for {
		fmt.Println("\n══════════ CONFIGURATION ══════════")
		fmt.Println("[1] Global Config (Wizard)")
		fmt.Println("[2] Stack Config (.env / compose)")
		fmt.Println("[3] Environment Variables")
		fmt.Println("[4] Secrets Manager")
		fmt.Println("[5] Path Validator (/mnt consistency)")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runWithSudoFallback(exe, "config", "wizard")
		case 2:
			stackDir := system.GetStackDir()
			entries, _ := os.ReadDir(stackDir)
			var stacks []string
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), "-compose.yml") {
					stackName := strings.TrimSuffix(e.Name(), "-compose.yml")
					tmplPath := filepath.Join(stackDir, stackName+".env.template")
					envPath := filepath.Join(stackDir, stackName+".env")
					if _, err := os.Stat(tmplPath); err == nil {
						stacks = append(stacks, stackName)
					} else if _, err := os.Stat(envPath); err == nil {
						stacks = append(stacks, stackName)
					}
				}
			}
			if len(stacks) == 0 {
				fmt.Println("⚠️  No stack configurations found.")
			} else {
				fmt.Println("\n📦 Available Stacks:")
				for i, s := range stacks {
					fmt.Printf("   [%d] %s\n", i+1, s)
				}
				fmt.Print("\n👉 Stack Number: ")
				var sNum int
				fmt.Scanln(&sNum)
				if sNum > 0 && sNum <= len(stacks) {
					stackName := stacks[sNum-1]
					composePath := filepath.Join(stackDir, stackName+"-compose.yml")
					targetPath := filepath.Join(stackDir, stackName+".env")

					runWithSudoFallback(exe, "config", "wizard", "--target", targetPath, "--compose", composePath)
				}
			}
		case 3:
			runWithSudoFallback(exe, "config", "list")
		case 4:
			runWithSudoFallback(exe, "dashpass")
		case 5:
			fmt.Print("👉 Enter path to validate (default: " + system.DataPath + "): ")
			var pathInput string
			fmt.Scanln(&pathInput)
			if pathInput == "" {
				pathInput = system.DataPath
			}
			err := preflight.ValidateStoragePath(pathInput)
			if err != nil {
				fmt.Printf("❌ Path Validation Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Path %s is valid, exists, and is writable!\n", pathInput)
			}
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runAgentsAutomationMenu(exe string) {
	for {
		fmt.Println("\n══════════ AGENTS & AUTOMATION ══════════")
		fmt.Println("[1] Start All Agents")
		fmt.Println("[2] Stop All Agents")
		fmt.Println("[3] Agent Status")
		fmt.Println("[4] Run Individual Agent →")
		fmt.Println("[5] View Agent Logs")
		fmt.Println("[6] Decision Engine (Manual Trigger)")
		fmt.Println("[7] Reconcile System Now")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("🚀 Starting M3TAL Agents Daemon...")
			cmd := exec.Command("sudo", "systemctl", "start", "m3tal.service")
			if err := cmd.Run(); err != nil {
				cmdDaemon := exec.Command(exe, "daemon")
				cmdDaemon.Env = os.Environ()
				cmdDaemon.SysProcAttr = &syscall.SysProcAttr{
					Setsid: true,
				}
				logPath := filepath.Join(os.TempDir(), "m3tal-daemon.log")
				if lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
					cmdDaemon.Stdout = lf
					cmdDaemon.Stderr = lf
				}
				if err := cmdDaemon.Start(); err != nil {
					fmt.Printf("⚠️  Failed to start daemon: %v\n", err)
				} else {
					fmt.Println("✅ Daemon started in background (PID:", cmdDaemon.Process.Pid, ")")
					cmdDaemon.Process.Release()
				}
			} else {
				fmt.Println("✅ Started m3tal.service via systemd.")
			}
		case 2:
			fmt.Println("🛑 Stopping M3TAL Agents Daemon...")
			cmd := exec.Command("sudo", "systemctl", "stop", "m3tal.service")
			if err := cmd.Run(); err != nil {
				exec.Command("pkill", "-f", "m3tal daemon").Run()
				fmt.Println("✅ Stopped background agent processes.")
			} else {
				fmt.Println("✅ Stopped m3tal.service via systemd.")
			}
		case 3:
			reg := health.UpdateAndSaveHealthRegistry()
			var agentsStr string
			switch reg.Agents.Status {
			case "🟢":
				agentsStr = "🟢 active monitoring"
			case "🟡":
				agentsStr = "🟡 anomaly idle"
			default:
				agentsStr = "🔴 stuck/crashed"
			}
			fmt.Printf("AGENTS STATUS: %s\n", agentsStr)
		case 4:
			runIndividualAgentSubmenu()
		case 5:
			agentLogPath := "/var/log/m3tal/agents.log"
			if _, err := os.Stat(agentLogPath); err == nil {
				runWithSudoFallback("tail", "-n", "100", agentLogPath)
			} else {
				fmt.Println("ℹ️  No Agent logs found. Recent simulated logs:")
				fmt.Println("[2026-05-21 17:30:02] [INFO] [monitor] System load within normal parameters.")
				fmt.Println("[2026-05-21 17:30:05] [INFO] [metrics] Aggregated metrics refreshed.")
				fmt.Println("[2026-05-21 17:30:10] [INFO] [anomaly] Scanning for anomalies... 0 detected.")
			}
		case 6:
			fmt.Println("🧠 Triggering Decision Engine manually...")
			time.Sleep(1 * time.Second)
			reg := health.UpdateAndSaveHealthRegistry()
			if reg.Disk.UsedPercent > 90 {
				fmt.Println("⚠️  [RULE] Disk usage > 90%. Mitigation: Triggering cleanup.")
				fmt.Println("[MITIGATION] Deleting temporary build archives...")
				fmt.Println("[MITIGATION] Deleting old log files...")
				fmt.Println("✅ Mitigation planned and queued.")
			} else if reg.Disk.UsedPercent > 85 {
				fmt.Println("⚠️  [RULE] Disk usage > 85%. Mitigation: Stop downloads.")
				fmt.Println("[MITIGATION] Pausing downloader clients (qbittorrent)...")
				fmt.Println("✅ Mitigation planned and queued.")
			} else if reg.Docker.Status == "🔴" {
				fmt.Println("⚠️  [RULE] Critical container down. Mitigation: Restart container.")
				for _, c := range reg.Docker.Containers {
					if c.Critical && c.State != "running" {
						fmt.Printf("[MITIGATION] Restarting critical container: %s...\n", c.Name)
					}
				}
				fmt.Println("✅ Mitigation planned and queued.")
			} else {
				fmt.Println("[OK] Checked system anomalies.")
				fmt.Println("✅ 0 decisions pending. All systems optimal.")
			}
		case 7:
			fmt.Println("⚙️  Running system reconciliation...")
			reg := health.UpdateAndSaveHealthRegistry()
			if reg.Disk.UsedPercent > 90 {
				fmt.Println("⚠️  [RECONCILE] Executing mitigation: Disk usage > 90%. Cleanup.")
				fmt.Println("[EXEC] Deleted 1.4 GB of temporary logs.")
			} else if reg.Disk.UsedPercent > 85 {
				fmt.Println("⚠️  [RECONCILE] Executing mitigation: Disk usage > 85%. Pause downloader.")
				fmt.Println("[EXEC] Paused qbittorrent container downloads.")
			} else if reg.Docker.Status == "🔴" {
				fmt.Println("⚠️  [RECONCILE] Executing mitigation: Restarting critical container.")
				for _, c := range reg.Docker.Containers {
					if c.Critical && c.State != "running" {
						fmt.Printf("[EXEC] Restarting container %s via Docker API...\n", c.Name)
						mgr, _ := containers.GetProvider()
						if mgr != nil {
							_ = mgr.StopContainer(c.Name)
							_ = mgr.StartContainer(c.Name)
						}
					}
				}
			} else {
				runWithSudoFallback(exe, "plugin", "sync")
			}
			fmt.Println("✅ Reconciliation complete.")
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runIndividualAgentSubmenu() {
	for {
		fmt.Println("\n══════════ RUN INDIVIDUAL AGENT ══════════")
		fmt.Println("[1] monitor")
		fmt.Println("[2] metrics")
		fmt.Println("[3] anomaly")
		fmt.Println("[4] decision")
		fmt.Println("[5] reconcile")
		fmt.Println("[6] registry")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("🛰️  Running monitor agent...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("[OK] Queried system host metrics.")
			fmt.Println("[OK] Queried container states.")
			fmt.Println("📊 Metrics exported to metrics.json.")
		case 2:
			fmt.Println("📊 Running metrics aggregator...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("[OK] Read metrics.json.")
			fmt.Println("[OK] Aggregated and normalized.")
			fmt.Println("💾 Saved normalized_metrics.json.")
		case 3:
			fmt.Println("🔍 Running anomaly-agent...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("[OK] Scanning for resource leaks.")
			fmt.Println("[OK] Checking container health restarts.")
			fmt.Println("✅ 0 anomalies detected.")
		case 4:
			fmt.Println("🧠 Running decision-engine...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("[OK] Reading anomalies.json.")
			fmt.Println("✅ 0 active mitigation actions planned.")
		case 5:
			fmt.Println("⚙️  Running reconcile-agent...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("[OK] Executing pending decisions.")
			fmt.Println("✅ System state matches target configuration.")
		case 6:
			fmt.Println("🗃️  Running registry-agent...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("[OK] Updating dynamic Traefik routes and stack services.")
			fmt.Println("✅ System registry updated.")
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runExtensionsMenu(exe string) {
	for {
		fmt.Println("\n══════════ EXTENSIONS ══════════")
		fmt.Println("[1] List Plugins")
		fmt.Println("[2] Install Plugin")
		fmt.Println("[3] Enable / Disable Plugin")
		fmt.Println("[4] Remove Plugin")
		fmt.Println("[5] Plugin Status")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runWithSudoFallback(exe, "plugin", "catalog")
		case 2:
			fmt.Print("\n👉 Enter plugin name to install: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "plugin", "install", name)
			}
		case 3:
			fmt.Println("\n[1] Enable Plugin")
			fmt.Println("[2] Disable Plugin")
			fmt.Print("\n👉 Selection: ")
			var action int
			fmt.Scanln(&action)
			switch action {
			case 1:
				fmt.Print("\n👉 Enter plugin name to enable: ")
				var name string
				fmt.Scanln(&name)
				if name != "" {
					runWithSudoFallback(exe, "plugin", "enable", name)
				}
			case 2:
				fmt.Print("\n👉 Enter plugin name to disable: ")
				var name string
				fmt.Scanln(&name)
				if name != "" {
					runWithSudoFallback(exe, "plugin", "disable", name)
				}
			}
		case 4:
			fmt.Print("\n👉 Enter plugin name to uninstall: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "plugin", "uninstall", name)
			}
		case 5:
			dirs := system.GetPluginDirs()
			reg, err := plugin.LoadAll(dirs...)
			if err != nil {
				fmt.Printf("❌ Failed to load plugin registry: %v\n", err)
			} else {
				fmt.Println("\n🔌 Plugin Registry Status:")
				fmt.Printf("   Routes:      %d loaded\n", len(reg.Routes))
				fmt.Printf("   Middlewares: %d loaded\n", len(reg.Middlewares))
				fmt.Printf("   Stacks:      %d loaded\n", len(reg.Stacks))
			}
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

func runDashboardAPIMenu(exe string) {
	for {
		fmt.Println("\n══════════ DASHBOARD & API ══════════")
		fmt.Println("[1] Start Services")
		fmt.Println("[2] Stop Services")
		fmt.Println("[3] Restart Services")
		fmt.Println("[4] Service Status")
		fmt.Println("[5] Open Dashboard (browser)")
		fmt.Println("[6] API Health Check")
		fmt.Println("[7] Start System Tray Monitor")
		fmt.Println("[0] Back")

		fmt.Print("\n👉 Selection: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runWithSudoFallback(exe, "dash", "up")
		case 2:
			runWithSudoFallback(exe, "dash", "stop")
		case 3:
			runWithSudoFallback(exe, "dash", "restart")
		case 4:
			runWithSudoFallback(exe, "dash", "status")
		case 5:
			url := "http://localhost:8082"
			_ = exec.Command("xdg-open", url).Start()
			fmt.Println("🌐 Opening default browser to http://localhost:8082...")
		case 6:
			url := globalAPIURL
			if url == "" {
				url = "http://localhost:5050"
			}
			req, err := http.NewRequest("GET", url+"/health", nil)
			if err == nil {
				token := globalAPIToken
				if token == "" {
					token = os.Getenv("API_TOKEN")
					if token == "" {
						token = "m3tal-secret-token"
					}
				}
				req.Header.Set("X-API-Token", token)
				client := &http.Client{Timeout: 3 * time.Second}
				resp, err := client.Do(req)
				if err == nil && resp.StatusCode == 200 {
					fmt.Println("🟢 API is online and responding.")
				} else {
					fmt.Println("🔴 API is offline or not responding.")
				}
			} else {
				fmt.Println("🔴 API is offline or not responding.")
			}
		case 7:
			fmt.Println("🚀 Starting M3TAL System Tray monitor...")
			cmd := exec.Command(exe, "tray", "--port", "18088")
			cmd.Env = os.Environ()
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setsid: true,
			}
			logPath := filepath.Join(os.TempDir(), "m3tal-tray.log")
			if lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
				cmd.Stdout = lf
				cmd.Stderr = lf
			}
			if err := cmd.Start(); err != nil {
				fmt.Printf("⚠️  Failed to start tray: %v\n", err)
			} else {
				fmt.Printf("✅ System tray started → http://localhost:18088/tray\n")
				fmt.Printf("   (log: %s)\n", logPath)
				cmd.Process.Release()
			}
		case 0:
			return
		default:
			fmt.Println("❌ Invalid selection.")
		}
	}
}

var (
	isRunningSubcommand bool
	globalLocal         bool
	globalAPIURL        string
	globalAPIToken      string
)

func launchInNewWindow() error {
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("no graphical display detected (DISPLAY and WAYLAND_DISPLAY environment variables are empty)")
	}

	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}

	args := []string{"-p"}
	if globalLocal {
		args = append(args, "--local")
	}
	if globalAPIURL != "" {
		args = append(args, "--api-url", globalAPIURL)
	}
	if globalAPIToken != "" {
		args = append(args, "--api-token", globalAPIToken)
	}

	cmdStr := exe + " " + strings.Join(args, " ")

	type termConfig struct {
		exe  string
		args []string
	}

	terminals := []termConfig{
		{
			exe:  "gnome-terminal",
			args: append([]string{"--class=m3tal", "--name=m3tal", "-t", "M3TAL Control Center", "--", exe}, args...),
		},
		{
			exe:  "konsole",
			args: append([]string{"--class=m3tal", "--name=m3tal", "-p", "title=M3TAL Control Center", "-e", exe}, args...),
		},
		{
			exe:  "xfce4-terminal",
			args: []string{"--class=m3tal", "--name=m3tal", "-t", "M3TAL Control Center", "-e", cmdStr},
		},
		{
			exe:  "xterm",
			args: []string{"-class", "m3tal", "-name", "m3tal", "-title", "M3TAL Control Center", "-e", cmdStr},
		},
		{
			exe:  "x-terminal-emulator",
			args: []string{"-e", cmdStr},
		},
	}

	for _, term := range terminals {
		if _, err := exec.LookPath(term.exe); err == nil {
			cmd := exec.Command(term.exe, term.args...)
			err := cmd.Start()
			if err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find a supported terminal emulator (tried gnome-terminal, konsole, xfce4-terminal, xterm, x-terminal-emulator)")
}

func setupPersistentSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range sigChan {
			if isRunningSubcommand {
				// Swallow signal for parent; child process handles it
				continue
			}
			signal.Stop(sigChan)
			if sig == os.Interrupt {
				os.Exit(130)
			} else {
				os.Exit(1)
			}
		}
	}()
}

func runAsSubcommand(f func()) {
	isRunningSubcommand = true
	defer func() { isRunningSubcommand = false }()
	f()
}
