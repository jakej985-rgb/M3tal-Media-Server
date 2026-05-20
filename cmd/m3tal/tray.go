package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/godbus/dbus/v5"
	"github.com/gogpu/systray"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/system"
	"github.com/spf13/cobra"
)

// Minimal placeholder HTML for the tray interface
var TrayHTML = `<!DOCTYPE html><html><head><title>M3TAL Tray</title></head><body></body></html>`

// Base64-encoded 1x1 transparent PNG icon
var IconBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNgYAAAAAMAAWgmWQ0AAAAAElFTkSuQmCC"

var IconData []byte

func init() {
	decoded, err := base64.StdEncoding.DecodeString(IconBase64)
	if err == nil {
		IconData = decoded
	}
}

// hasStatusNotifier checks if org.kde.StatusNotifierWatcher is registered on the session bus
func hasStatusNotifier() bool {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return false
	}
	defer conn.Close()

	var owner string
	err = conn.BusObject().Call("org.freedesktop.DBus.GetNameOwner", 0, "org.kde.StatusNotifierWatcher").Store(&owner)
	if err != nil {
		return false
	}
	return owner != ""
}

var trayCmd = &cobra.Command{
	Use:   "tray",
	Short: "Run the M3TAL system tray monitor",
	Long:  "Starts a system tray icon that launches a browser-based stats monitor popup showing real-time CPU, GPU, storage, and container metrics.",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		runTray(port)
	},
}

func init() {
	trayCmd.Flags().StringP("port", "p", "18088", "Port to run the tray web server on")
}

func runTray(port string) {
	// Bind a TCP listener on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		log.Fatalf("Tray web server failed to bind to port %s: %v", port, err)
	}

	fmt.Printf("🚀 Starting M3TAL System Tray monitor on port %s...\n", port)
	fmt.Printf("👉 Access stats directly at http://localhost:%s/tray\n", port)

	if runtime.GOOS == "linux" && os.Getenv("DBUS_SESSION_BUS_ADDRESS") == "" {
		fmt.Println("\nℹ️  Note: DBUS_SESSION_BUS_ADDRESS is not set (common when running under sudo).")
		fmt.Println("   The system tray icon may not appear in your desktop bar, but the")
		fmt.Printf("   stats web interface remains fully accessible at http://localhost:%s/tray\n", port)
	}

	// Start HTTP server serving the tray interface and API endpoints
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
		if err := http.Serve(listener, mux); err != nil {
			log.Fatalf("Tray web server failed: %v", err)
		}
	}()

	// Run systray (pure Go, zero CGO, blocks on Run)
	mode := os.Getenv("M3TAL_TRAY")
	if mode == "off" {
		log.Println("[tray] M3TAL_TRAY is set to off; system tray disabled")
		return
	}
	if mode != "force" && !hasStatusNotifier() {
		log.Println("[tray] StatusNotifierWatcher not available; system tray disabled")
		return
	}

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
