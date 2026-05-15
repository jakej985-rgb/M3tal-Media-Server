package api

import (
	"log"
	"net/http"

	"github.com/jakej985-rgb/m3tal-core/internal/containers"
)

// StartServer starts the API server on the specified port
func StartServer(port string, token string) error {
	srv := NewServer(token)

	log.Printf("🚀 M3TAL API Interface starting on :%s...\n", port)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/containers", srv.AuthMiddleware(srv.ListContainers))
	mux.HandleFunc("/api/containers/list", srv.AuthMiddleware(srv.ListContainers))
	mux.HandleFunc("/api/ps", srv.AuthMiddleware(srv.ListContainers)) // Added ps alias
	mux.HandleFunc("/api/metrics", srv.AuthMiddleware(srv.GetStats))

	mux.HandleFunc("/api/containers/start", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mgr, _ := containers.GetProvider()
		srv.HandleContainerAction(w, r, mgr.StartContainer)
	}))
	mux.HandleFunc("/api/containers/stop", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mgr, _ := containers.GetProvider()
		srv.HandleContainerAction(w, r, mgr.StopContainer)
	}))
	mux.HandleFunc("/api/containers/restart", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mgr, _ := containers.GetProvider()
		srv.HandleContainerAction(w, r, mgr.RestartContainer)
	}))

	mux.HandleFunc("/api/health", srv.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		srv.ListContainers(w, r)
	}))

	return http.ListenAndServe(":"+port, mux)
}
