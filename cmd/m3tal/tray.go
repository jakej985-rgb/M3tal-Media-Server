package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gogpu/systray"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
	"github.com/spf13/cobra"
)

var trayCmd = &cobra.Command{
	Use:   "tray",
	Short: "Run the M3TAL system tray monitor",
	Long:  "Starts a system tray icon that launches a browser-based stats monitor popup showing real-time CPU, GPU, storage, and container metrics.",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		fmt.Printf("🚀 Starting M3TAL System Tray monitor on port %s...\n", port)
		fmt.Printf("👉 Access stats directly at http://localhost:%s/tray\n", port)

		if runtime.GOOS == "linux" && os.Getenv("DBUS_SESSION_BUS_ADDRESS") == "" {
			fmt.Println("\nℹ️  Note: DBUS_SESSION_BUS_ADDRESS is not set (common when running under sudo).")
			fmt.Println("   The system tray icon may not appear in your desktop bar, but the")
			fmt.Printf("   stats web interface remains fully accessible at http://localhost:%s/tray\n", port)
		}
		runTray(port)
	},
}

func init() {
	trayCmd.Flags().StringP("port", "p", "18088", "Port to run the tray web server on")
}

func runTray(port string) {
	// Start HTTP server in a goroutine
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/tray", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(TrayHTML))
		})
		mux.HandleFunc("/tray/api/stats", func(w http.ResponseWriter, r *http.Request) {
			stats, err := system.GetDetailedStats()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(stats)
		})
		mux.HandleFunc("/tray/api/containers", func(w http.ResponseWriter, r *http.Request) {
			provider, err := containers.GetProvider()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			list, err := provider.ListContainers()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(list)
		})
		if err := http.ListenAndServe(net.JoinHostPort("127.0.0.1", port), mux); err != nil {
			log.Fatalf("Tray web server failed: %v", err)
		}
	}()

	// Run systray (pure Go, zero CGO, blocks on Run)
	tray := systray.New()
	tray.SetIcon(IconData)
	tray.SetTooltip("M3TAL Core System Tray")

	menu := systray.NewMenu()
	menu.Add("Open Monitor", func() {
		openBrowser(fmt.Sprintf("http://localhost:%s/tray", port))
	})
	menu.Add("Open Dashboard", func() {
		openBrowser("http://localhost:8082")
	})
	menu.AddSeparator()
	menu.Add("Quit", func() {
		tray.Remove()
		os.Exit(0)
	})

	tray.SetMenu(menu)
	tray.Show()
	tray.Run()
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
