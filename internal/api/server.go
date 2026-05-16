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

	// Required endpoints from master plan
	mux.HandleFunc("/health", srv.AuthMiddleware(srv.GetHealth))
	mux.HandleFunc("/services", srv.AuthMiddleware(srv.GetServices))
	mux.HandleFunc("/stack", srv.AuthMiddleware(srv.GetStack))
	mux.HandleFunc("/config", srv.AuthMiddleware(srv.GetConfig))

	// Backward compatibility and convenience
	mux.HandleFunc("/api/health", srv.AuthMiddleware(srv.GetHealth))
	mux.HandleFunc("/api/services", srv.AuthMiddleware(srv.GetServices))
	mux.HandleFunc("/api/stack", srv.AuthMiddleware(srv.GetStack))
	mux.HandleFunc("/api/config", srv.AuthMiddleware(srv.GetConfig))
	mux.HandleFunc("/api/containers", srv.AuthMiddleware(srv.GetServices))

	// Actions
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

	return http.ListenAndServe(":"+port, mux)
}
