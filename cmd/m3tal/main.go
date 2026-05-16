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

	"bytes"
	_ "embed"
	"io"
	"net/http"

	"github.com/jakej985-rgb/m3tal-core/internal/api"
	"github.com/jakej985-rgb/m3tal-core/internal/auth"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/orchestrator"
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

	var rootCmd = &cobra.Command{Use: "m3tal"}
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
			usersFile := filepath.Join(system.GetStackDir(), "users.json")
			if err := auth.UpdateUser(usersFile, username, password); err != nil {
				log.Fatal(err)
			}
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
			runWizard(false)

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
				Template string `json:"template"`
			}
			stacks := make(map[string]stackInfo)
			for _, e := range entries {
				name := e.Name()
				if strings.HasSuffix(name, "-compose.yml") {
					stack := strings.TrimSuffix(name, "-compose.yml")
					composePath := filepath.Join(stackDir, name)
					templatePath := filepath.Join(stackDir, stack+".env.template")
					if _, err := os.Stat(templatePath); os.IsNotExist(err) {
						continue
					}
					stacks[stack] = stackInfo{Compose: composePath, Template: templatePath}
				}
			}
			printJSON(stacks)
		},
	}
	configCmd.AddCommand(configListCmd, configSetCmd, configGetCmd, configScanCmd, configWizardCmd)
	rootCmd.AddCommand(listCmd, psCmd, startCmd, stopCmd, statsCmd, daemonCmd, apiCmd, upCmd, downCmd, logsCmd, pullCmd, dashpassCmd, initCmd, docCmd, configCmd, initDashCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runWizard(update bool) {
	fmt.Println("🛠️  M3TAL Configuration Wizard")

	targetFile := system.GetConfigPath()
	isSystem := false
	if runtime.GOOS == "linux" && os.Geteuid() == 0 {
		isSystem = true
		configDir := filepath.Dir(system.ConfigPath)
		_ = os.MkdirAll(configDir, 0755)
		_ = os.MkdirAll(system.DataPath, 0755)
		fmt.Printf("✅ System directories initialized (%s, %s)\n", configDir, system.DataPath)
	}

	_ = os.MkdirAll("./data", 0755)

	sourceFile := ".env.example"
	var data []byte
	var err error

	if update {
		if _, err = os.Stat(targetFile); err == nil {
			sourceFile = targetFile
			data, err = os.ReadFile(sourceFile)
		} else if _, err = os.Stat(".env"); err == nil {
			sourceFile = ".env"
			data, err = os.ReadFile(sourceFile)
		}
	}

	if data == nil {
		// Use embedded example if not an update or if update source missing
		data = []byte(envExample)
	}

	if err != nil && data == nil {
		log.Fatalf("❌ Missing configuration source.")
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
			orchestrator.RunRaw("journalctl", "-u", "m3tal", "-f", "-n", "50")
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
		if subChoice == 1 {
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
				stackMgr.Run("logs", "--tail", "50", "-f")
			}
		} else if subChoice == 2 {
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
				orchestrator.RunRaw("docker", "logs", "--tail", "50", "-f", list[cNum-1].Names[0])
			}
		}
	case 3:
		fmt.Println("\n🚀 Streaming Aggregated Logs...")
		stackMgr := orchestrator.NewStackManager()
		stackMgr.Run("logs", "--tail", "20", "-f")
	case 0:
		return
	default:
		fmt.Println("❌ Invalid selection.")
	}
}
