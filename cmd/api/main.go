package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jakej985-rgb/m3tal-core/internal/api"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
)

func main() {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		apiToken = "m3tal-secret-token"
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5050"
	}

	// Start agents daemon if not already running
	startAgentsDaemonIfMissing()

	// Initialize SQLite store
	dbPath := store.GetStatePath()
	db, err := store.Open(dbPath)
	if err != nil {
		log.Printf("⚠️  Could not open state database at %s: %v", dbPath, err)
		log.Println("⚠️  v2 engine endpoints will be disabled. Starting with v1 only.")
		if err := api.StartServer(port, apiToken); err != nil {
			log.Fatalf("❌ API server failed: %v", err)
		}
		return
	}
	defer db.Close()

	log.Printf("📦 State database: %s\n", dbPath)

	if err := api.StartServerWithStore(port, apiToken, db); err != nil {
		log.Fatalf("❌ API server failed: %v", err)
	}
}

func startAgentsDaemonIfMissing() {
	go func() {
		// Give the API a moment to bind and settle, then check/start agents
		time.Sleep(1 * time.Second)

		// Check if agents daemon is running
		cmdCheck := exec.Command("pgrep", "-f", "m3tal daemon")
		if err := cmdCheck.Run(); err != nil {
			// not running, start it
			log.Println("🚀 Starting M3TAL Agents Daemon...")
			// try systemd first
			cmdSys := exec.Command("sudo", "systemctl", "start", "m3tal.service")
			if err := cmdSys.Run(); err != nil {
				// fallback to running binary directly
				exe, err := os.Executable()
				if err == nil {
					dir := filepath.Dir(exe)
					m3talExe := filepath.Join(dir, "m3tal")
					if _, err := os.Stat(m3talExe); err != nil {
						m3talExe = "m3tal"
					}
					cmdDaemon := exec.Command(m3talExe, "daemon")
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
						log.Printf("⚠️  Failed to start agents daemon: %v\n", err)
					} else {
						log.Printf("✅ Agents daemon started in background (PID: %d)\n", cmdDaemon.Process.Pid)
						_ = cmdDaemon.Process.Release()
					}
				}
			} else {
				log.Println("✅ Started m3tal.service via systemd.")
			}
		}
	}()
}
