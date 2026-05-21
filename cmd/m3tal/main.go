package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"bytes"
	_ "embed"
	"io"
	"net/http"

	"github.com/jakej985-rgb/m3tal-core/internal/api"
	"github.com/jakej985-rgb/m3tal-core/internal/auth"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/orchestrator"
	"github.com/jakej985-rgb/m3tal-core/internal/plugin"
	"github.com/jakej985-rgb/m3tal-core/internal/preflight"
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
	rootCmd.PersistentFlags().String("api-url", "http://localhost:8080", "M3TAL API URL")
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
			if err := api.StartServer(port, token); err != nil {
				log.Fatalf("❌ API server failed: %v", err)
			}
		},
	}
	apiCmd.Flags().String("port", "8080", "Port to listen on")

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

			newPath, err := plugin.EnablePlugin(path)
			if err != nil {
				log.Fatalf("❌ Failed to enable: %v", err)
			}
			fmt.Printf("✅ Enabled plugin %q (renamed to %s)\n", name, filepath.Base(newPath))
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

			newPath, err := plugin.DisablePlugin(path)
			if err != nil {
				log.Fatalf("❌ Failed to disable: %v", err)
			}
			fmt.Printf("✅ Disabled plugin %q (renamed to %s)\n", name, filepath.Base(newPath))
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

			// Find kind
			var targetKind string
			for i := range reg.Routes {
				if plugin.MatchesPluginName(reg.Routes[i].SourcePath, reg.Routes[i].Metadata.Name, name) {
					targetKind = "Route"
					name = plugin.GetPluginBaseName(reg.Routes[i].SourcePath)
					break
				}
			}
			if targetKind == "" {
				for i := range reg.Stacks {
					if plugin.MatchesPluginName(reg.Stacks[i].SourcePath, reg.Stacks[i].Metadata.Name, name) {
						targetKind = "Stack"
						name = plugin.GetPluginBaseName(reg.Stacks[i].SourcePath)
						break
					}
				}
			}
			if targetKind == "" {
				for i := range reg.Middlewares {
					if plugin.MatchesPluginName(reg.Middlewares[i].SourcePath, reg.Middlewares[i].Metadata.Name, name) {
						targetKind = "Middleware"
						name = plugin.GetPluginBaseName(reg.Middlewares[i].SourcePath)
						break
					}
				}
			}

			if targetKind == "" {
				log.Fatalf("❌ Plugin %q not found in local registry", name)
			}

			userDir := system.UserPluginsDir
			if _, err := os.Stat("deploy/plugins"); err == nil {
				userDir = "deploy/plugins"
			}

			fmt.Printf("🗑️  Uninstalling plugin %q...\n", name)
			err = plugin.UninstallPlugin(name, targetKind, userDir, reg)
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
			mgr, _ := containers.GetProvider()
			list, _ := mgr.ListContainers()
			for i, c := range list {
				fmt.Printf("   [%d] %s (%s)\n", i+1, c.Names[0], c.Status)
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

func runMainMenu() bool {
	fmt.Println("\n🛠️  M3TAL Control Center")
	fmt.Println("|-- 1.) Container Management")
	fmt.Println("|-- 2.) View Logs Explorer")
	fmt.Println("|-- 3.) Start Dashboard & API")
	fmt.Println("|-- 4.) Configuration & Secrets")
	fmt.Println("|-- 5.) System Health Check (Doctor)")
	fmt.Println("|-- 6.) Start System Tray Monitor")
	fmt.Println("|-- 7.) Manage Plugins")
	fmt.Println("|-- 0.) Exit")

	fmt.Print("\n👉 Selection: ")
	var choice int
	fmt.Scanln(&choice)

	exe := os.Args[0]

	switch choice {
	case 1:
		fmt.Println("\n|-- 1.) Container Management")
		fmt.Println("|   [1] Start Environment (Up)")
		fmt.Println("|   [2] Stop Environment (Down)")
		fmt.Println("|   [3] Update Images (Pull)")
		fmt.Println("|   [4] System Status (PS)")
		fmt.Print("\n👉 Selection: ")
		var subChoice int
		fmt.Scanln(&subChoice)
		switch subChoice {
		case 1:
			runWithSudoFallback(exe, "up")
		case 2:
			runWithSudoFallback(exe, "down")
		case 3:
			runWithSudoFallback(exe, "pull")
		case 4:
			runWithSudoFallback(exe, "ps")
		default:
			fmt.Println("❌ Invalid selection.")
		}
	case 2:
		runWithSudoFallback(exe, "logs")
	case 3:
		runWithSudoFallback(exe, "dash", "up")
	case 4:
		fmt.Println("\n|-- 4.) Configuration & Secrets")
		fmt.Println("|   [1] Edit Global Configuration (Wizard)")
		fmt.Println("|   [2] Edit Stack Configuration")
		fmt.Println("|   [3] Scan & List All Variables")
		fmt.Println("|   [4] Manage Dashboard Users")
		fmt.Print("\n👉 Selection: ")
		var subChoice int
		fmt.Scanln(&subChoice)
		switch subChoice {
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
		default:
			fmt.Println("❌ Invalid selection.")
		}
	case 5:
		runWithSudoFallback(exe, "doctor")
	case 6:
		// Launch tray in the background so the menu returns immediately.
		// Must run as the current user (not sudo) to keep DBUS_SESSION_BUS_ADDRESS.
		// Stdout/Stderr are redirected to a log file – NOT the parent terminal,
		// otherwise the tray's startup messages appear inside the menu window.
		cmd := exec.Command(exe, "tray", "--port", "18088")
		cmd.Env = os.Environ() // inherit full user environment (incl. D-Bus)
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
			cmd.Process.Release() // detach – don't wait for it
		}
	case 7:
		fmt.Println("\n|-- 7.) Plugin Management")
		fmt.Println("|   [1] List Discovered & Catalog Plugins")
		fmt.Println("|   [2] Download & Install Plugin")
		fmt.Println("|   [3] Enable Plugin")
		fmt.Println("|   [4] Disable Plugin")
		fmt.Println("|   [5] Uninstall Plugin")
		fmt.Println("|   [6] Sync Gateway Configuration")
		fmt.Print("\n👉 Selection: ")
		var subChoice int
		fmt.Scanln(&subChoice)
		switch subChoice {
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
			fmt.Print("\n👉 Enter plugin name to enable: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "plugin", "enable", name)
			}
		case 4:
			fmt.Print("\n👉 Enter plugin name to disable: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "plugin", "disable", name)
			}
		case 5:
			fmt.Print("\n👉 Enter plugin name to uninstall: ")
			var name string
			fmt.Scanln(&name)
			if name != "" {
				runWithSudoFallback(exe, "plugin", "uninstall", name)
			}
		case 6:
			runWithSudoFallback(exe, "plugin", "sync")
		default:
			fmt.Println("❌ Invalid selection.")
		}
	case 0:
		return false
	default:
		fmt.Println("❌ Invalid selection.")
	}
	return true
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

	terminals := [][]string{
		{"x-terminal-emulator", "-e", cmdStr},
		{"gnome-terminal", "--", exe},
		{"konsole", "-e", cmdStr},
		{"xfce4-terminal", "-e", cmdStr},
		{"xterm", "-e", cmdStr},
	}

	for _, term := range terminals {
		termExe := term[0]
		if _, err := exec.LookPath(termExe); err == nil {
			var cmd *exec.Cmd
			if termExe == "gnome-terminal" {
				gnomeArgs := append([]string{"--", exe}, args...)
				cmd = exec.Command("gnome-terminal", gnomeArgs...)
			} else {
				cmd = exec.Command(termExe, term[1:]...)
			}
			err := cmd.Start()
			if err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find a supported terminal emulator (tried x-terminal-emulator, gnome-terminal, konsole, xfce4-terminal, xterm)")
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
