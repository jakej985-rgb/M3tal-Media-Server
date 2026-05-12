package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jakej985-rgb/m3tal-core/pkg/api"
	"github.com/jakej985-rgb/m3tal-core/pkg/containers"
)

func main() {
	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		stateDir = filepath.Join("..", "state")
	}

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		apiToken = "m3tal-secret-token"
	}

	srv := api.NewServer(apiToken)

	log.Println("🚀 M3TAL API Interface starting on :5050...")

	http.HandleFunc("/api/containers", srv.AuthMiddleware(srv.ListContainers))
	http.HandleFunc("/api/containers/list", srv.AuthMiddleware(srv.ListContainers))
	http.HandleFunc("/api/metrics", srv.AuthMiddleware(srv.GetStats))

	http.HandleFunc("/api/containers/start", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mgr, _ := containers.GetProvider()
		srv.HandleContainerAction(w, r, mgr.StartContainer)
	}))
	http.HandleFunc("/api/containers/stop", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mgr, _ := containers.GetProvider()
		srv.HandleContainerAction(w, r, mgr.StopContainer)
	}))
	http.HandleFunc("/api/containers/restart", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mgr, _ := containers.GetProvider()
		srv.HandleContainerAction(w, r, mgr.RestartContainer)
	}))

	http.HandleFunc("/api/health", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		srv.ListContainers(w, r)
	}))

	if err := http.ListenAndServe(":5050", nil); err != nil {
		log.Fatalf("❌ API server failed: %v", err)
	}
}
