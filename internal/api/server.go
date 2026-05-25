package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jakej985-rgb/m3tal-core/internal/containers"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
)

// StartServer starts the API server on the specified port.
// This is the legacy entrypoint that creates an in-memory server without a store.
func StartServer(port string, token string) error {
	return StartServerWithStore(port, token, nil)
}

// StartServerWithStore starts the API server with an optional SQLite store
// for the v2 engine endpoints.
func StartServerWithStore(port string, token string, db *store.Store) error {
	srv := NewServer(token)

	log.Printf("🚀 M3TAL API Interface starting on :%s...\n", port)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// ─── v1 Legacy Endpoints (backward compatible) ───
	r.Group(func(r chi.Router) {
		r.Use(srv.chiAuthMiddleware)

		// Original endpoints
		r.Get("/health", srv.GetHealth)
		r.Get("/services", srv.GetServices)
		r.Get("/stack", srv.GetStack)
		r.Get("/config", srv.GetConfig)

		// /api/* prefixed duplicates
		r.Get("/api/health", srv.GetHealth)
		r.Get("/api/services", srv.GetServices)
		r.Get("/api/stack", srv.GetStack)
		r.Get("/api/config", srv.GetConfig)
		r.Get("/api/containers", srv.GetServices)
		r.Get("/api/containers/list", srv.GetServices)
		r.Get("/api/metrics", srv.GetStats)
		r.Post("/ai/run", srv.AIRun)
		r.Post("/api/v2/ai/run", srv.AIRun)

		// Container actions
		r.Post("/api/containers/start", func(w http.ResponseWriter, r *http.Request) {
			mgr, err := containers.GetProvider()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Docker provider unavailable: "+err.Error())
				return
			}
			srv.HandleContainerAction(w, r, mgr.StartContainer)
		})
		r.Post("/api/containers/stop", func(w http.ResponseWriter, r *http.Request) {
			mgr, err := containers.GetProvider()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Docker provider unavailable: "+err.Error())
				return
			}
			srv.HandleContainerAction(w, r, mgr.StopContainer)
		})
		r.Post("/api/containers/restart", func(w http.ResponseWriter, r *http.Request) {
			mgr, err := containers.GetProvider()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Docker provider unavailable: "+err.Error())
				return
			}
			srv.HandleContainerAction(w, r, mgr.RestartContainer)
		})
	})

	// ─── v2 Engine Endpoints ───
	if db != nil {
		routeH := &RouteHandlers{Store: db}
		stackH := &StackHandlers{Store: db}
		mwH := &MiddlewareHandlers{Store: db}
		pluginH := NewPluginHandlers(db)
		composeH := &ComposeHandlers{Store: db}
		vpnH := NewVPNHandlers()
		proxyH := NewProxyHandlers(db)

		r.Route("/api/v2", func(r chi.Router) {
			r.Use(srv.chiAuthMiddleware)

			// Routes
			r.Get("/routes", routeH.ListRoutes)
			r.Post("/routes", routeH.CreateRoute)
			r.Delete("/routes/{id}", routeH.DeleteRoute)
			r.Post("/routes/preview", routeH.GetRouteLabels)

			// Proxy
			r.Get("/proxy/discover", proxyH.DiscoverServices)
			r.Post("/proxy/expose", proxyH.ExposeService)
			r.Post("/proxy/unexpose", proxyH.UnexposeService)
			r.Post("/proxy/secure", proxyH.ConfigureSSL)

			// Stacks
			r.Get("/stacks", stackH.ListStacks)
			r.Post("/stacks/load", stackH.LoadStack)
			r.Post("/stacks/deploy", stackH.DeployStack)
			r.Post("/stacks/scan", stackH.ScanStacks)

			// Services (pass-through to Docker provider)
			r.Get("/services", srv.GetServices)

			// Middleware
			r.Get("/middleware", mwH.ListMiddleware)
			r.Post("/middleware", mwH.CreateMiddleware)

			// Plugins
			r.Get("/plugins", pluginH.ListPlugins)
			r.Get("/plugins/catalog", pluginH.ListCatalog)
			r.Get("/plugins/{kind}", pluginH.ListPluginsByKind)
			r.Post("/plugins/reload", pluginH.Reload)
			r.Post("/plugins/enable", pluginH.Enable)
			r.Post("/plugins/disable", pluginH.Disable)
			r.Post("/plugins/sync", pluginH.Sync)
			r.Post("/plugins/install", pluginH.Install)
			r.Post("/plugins/uninstall", pluginH.Uninstall)

			// Compose
			r.Post("/compose/validate", composeH.ValidateCompose)
			r.Post("/compose/fix", composeH.FixCompose)
			r.Get("/compose/templates", composeH.ListTemplates)
			r.Post("/compose/generate", composeH.GenerateTemplate)
			r.Post("/compose/save", composeH.SaveCompose)

			// VPN
			r.Get("/vpn/status", vpnH.GetStatus)
			r.Post("/vpn/control", vpnH.ControlVPN)
			r.Post("/vpn/region", vpnH.SwitchRegion)
			r.Post("/vpn/sync-port", vpnH.SyncPort)
			r.Get("/vpn/check-leak", vpnH.CheckLeak)

			// AI
			r.Post("/ai/run", srv.AIRun)
		})

		log.Println("✅ v2 engine endpoints enabled (SQLite store active)")
	} else {
		// Plugin endpoints work without a store
		pluginH := NewPluginHandlers(nil)
		composeH := &ComposeHandlers{Store: nil}
		vpnH := NewVPNHandlers()
		proxyH := NewProxyHandlers(nil)
		r.Route("/api/v2", func(r chi.Router) {
			r.Use(srv.chiAuthMiddleware)
			// Plugins
			r.Get("/plugins", pluginH.ListPlugins)
			r.Get("/plugins/catalog", pluginH.ListCatalog)
			r.Get("/plugins/{kind}", pluginH.ListPluginsByKind)
			r.Post("/plugins/reload", pluginH.Reload)
			r.Post("/plugins/enable", pluginH.Enable)
			r.Post("/plugins/disable", pluginH.Disable)
			r.Post("/plugins/sync", pluginH.Sync)
			r.Post("/plugins/install", pluginH.Install)
			r.Post("/plugins/uninstall", pluginH.Uninstall)

			// Proxy
			r.Get("/proxy/discover", proxyH.DiscoverServices)
			r.Post("/proxy/expose", proxyH.ExposeService)
			r.Post("/proxy/unexpose", proxyH.UnexposeService)
			r.Post("/proxy/secure", proxyH.ConfigureSSL)

			// Compose
			r.Post("/compose/validate", composeH.ValidateCompose)
			r.Post("/compose/fix", composeH.FixCompose)
			r.Get("/compose/templates", composeH.ListTemplates)
			r.Post("/compose/generate", composeH.GenerateTemplate)
			r.Post("/compose/save", composeH.SaveCompose)

			// VPN
			r.Get("/vpn/status", vpnH.GetStatus)
			r.Post("/vpn/control", vpnH.ControlVPN)
			r.Post("/vpn/region", vpnH.SwitchRegion)
			r.Post("/vpn/sync-port", vpnH.SyncPort)
			r.Get("/vpn/check-leak", vpnH.CheckLeak)

			// AI
			r.Post("/ai/run", srv.AIRun)
		})

		log.Println("⚠️  v2 engine endpoints disabled (no store configured), plugin endpoints still available")
	}

	return http.ListenAndServe(":"+port, r)
}

// chiAuthMiddleware is a chi-compatible version of AuthMiddleware.
func (s *Server) chiAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-Token")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token != s.APIToken {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// writeJSON encodes a value as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
