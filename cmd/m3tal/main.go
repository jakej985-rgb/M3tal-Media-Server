package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jakej985-rgb/m3tal-core/pkg/containers"
	"github.com/jakej985-rgb/m3tal-core/pkg/system"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var rootCmd = &cobra.Command{Use: "m3tal"}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			mgr, err := containers.NewManager()
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
			mgr, err := containers.NewManager()
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
			mgr, err := containers.NewManager()
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
			runCompose("up", "-d")
		},
	}

	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Stop all M3TAL stacks",
		Run: func(cmd *cobra.Command, args []string) {
			runCompose("down")
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull latest images",
		Run: func(cmd *cobra.Command, args []string) {
			runCompose("pull")
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

			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("✅ Updating user %s...\n", username)
			usersFile := "source/m3tal-stack/users.json"
			updateUser(usersFile, username, string(hash))
		},
	}

	rootCmd.AddCommand(listCmd, startCmd, stopCmd, statsCmd, daemonCmd, upCmd, downCmd, pullCmd, dashpassCmd)

	// Default to daemon if no args
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "daemon")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func runCompose(args ...string) {
	stacks := []string{
		"source/m3tal-stack/network-compose.yml",
		"source/m3tal-stack/routing-compose.yml",
		"source/m3tal-stack/m3tal-compose.yml",
	}

	for _, stack := range stacks {
		if _, err := os.Stat(stack); os.IsNotExist(err) {
			continue
		}
		fmt.Printf("🚀 Running docker compose %v on %s...\n", args, stack)
		cmdArgs := append([]string{"compose", "-f", stack}, args...)
		c := exec.Command("docker", cmdArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("❌ Failed to run compose on %s: %v\n", stack, err)
		}
	}
}

func updateUser(filePath, username, hash string) {
	type User struct {
		Username string `json:"username"`
		Hash     string `json:"token_hash"`
		Role     string `json:"role"`
	}
	var users []User
	data, err := os.ReadFile(filePath)
	if err == nil {
		json.Unmarshal(data, &users)
	}

	found := false
	for i, u := range users {
		if u.Username == username {
			users[i].Hash = hash
			found = true
			break
		}
	}
	if !found {
		users = append(users, User{Username: username, Hash: hash, Role: "admin"})
	}

	newData, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(filePath, newData, 0644)
}
