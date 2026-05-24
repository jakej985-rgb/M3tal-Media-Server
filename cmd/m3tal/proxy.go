package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jakej985-rgb/m3tal-core/internal/proxy"
	"github.com/jakej985-rgb/m3tal-core/internal/store"
	"github.com/spf13/cobra"
)

var (
	exposeDomain      string
	exposePort        int
	exposeSSL         bool
	exposeMiddlewares []string
)

func initProxyCmds() *cobra.Command {
	var proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "Manage reverse proxy routes, auto-discovery, and SSL",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var proxyStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "List all active reverse proxy routes",
		Run: func(cmd *cobra.Command, args []string) {
			local, _ := cmd.Flags().GetBool("local")
			if !local {
				if resp, err := callAPI(cmd, "GET", "/api/v2/routes", nil); err == nil {
					var list []store.RouteRecord
					if err := json.Unmarshal([]byte(resp), &list); err == nil {
						printRoutesTable(list)
						return
					}
				}
			}

			// Local fallback
			dbPath := store.GetStatePath()
			db, err := store.Open(dbPath)
			if err != nil {
				log.Fatalf("❌ Failed to open database: %v", err)
			}
			defer db.Close()

			list, err := db.ListRoutes()
			if err != nil {
				log.Fatalf("❌ Failed to list routes: %v", err)
			}
			printRoutesTable(list)
		},
	}

	var proxyDiscoverCmd = &cobra.Command{
		Use:   "discover",
		Short: "Auto-discover proxy candidate containers running on the host",
		Run: func(cmd *cobra.Command, args []string) {
			local, _ := cmd.Flags().GetBool("local")
			if !local {
				if resp, err := callAPI(cmd, "GET", "/api/v2/proxy/discover", nil); err == nil {
					var list []proxy.DiscoverableService
					if err := json.Unmarshal([]byte(resp), &list); err == nil {
						printDiscoverTable(list)
						return
					}
				}
			}

			// Local fallback
			dbPath := store.GetStatePath()
			db, err := store.Open(dbPath)
			if err != nil {
				log.Fatalf("❌ Failed to open database: %v", err)
			}
			defer db.Close()

			mgr := proxy.NewManager(db)
			list, err := mgr.DiscoverServices()
			if err != nil {
				log.Fatalf("❌ Failed to discover services: %v", err)
			}
			printDiscoverTable(list)
		},
	}

	var proxyExposeCmd = &cobra.Command{
		Use:   "expose [container-name]",
		Short: "Expose a container service via Traefik routing rule",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			service := args[0]
			local, _ := cmd.Flags().GetBool("local")

			if exposeDomain == "" {
				log.Fatal("❌ Domain must be specified via --domain")
			}
			if exposePort == 0 {
				log.Fatal("❌ Internal port must be specified via --port")
			}

			if !local {
				body := map[string]interface{}{
					"service":     service,
					"domain":      exposeDomain,
					"port":        exposePort,
					"ssl":         exposeSSL,
					"middlewares": exposeMiddlewares,
				}
				if _, err := callAPI(cmd, "POST", "/api/v2/proxy/expose", body); err == nil {
					fmt.Printf("✅ Exposed %s on https://%s (port %d, via API)\n", service, exposeDomain, exposePort)
					return
				}
			}

			// Local fallback
			dbPath := store.GetStatePath()
			db, err := store.Open(dbPath)
			if err != nil {
				log.Fatalf("❌ Failed to open database: %v", err)
			}
			defer db.Close()

			mgr := proxy.NewManager(db)
			err = mgr.ExposeService(service, exposeDomain, exposePort, exposeSSL, exposeMiddlewares)
			if err != nil {
				log.Fatalf("❌ Expose failed: %v", err)
			}

			scheme := "http"
			if exposeSSL {
				scheme = "https"
			}
			fmt.Printf("✅ Exposed %s on %s://%s (port %d)\n", service, scheme, exposeDomain, exposePort)
		},
	}

	proxyExposeCmd.Flags().StringVar(&exposeDomain, "domain", "", "Public/internal domain name (e.g. app.local)")
	proxyExposeCmd.Flags().IntVar(&exposePort, "port", 0, "Internal container port to forward to")
	proxyExposeCmd.Flags().BoolVar(&exposeSSL, "ssl", false, "Enable Let's Encrypt SSL/TLS automation for this route")
	proxyExposeCmd.Flags().StringSliceVar(&exposeMiddlewares, "middlewares", []string{}, "Comma-separated list of middlewares to attach")

	var proxyUnexposeCmd = &cobra.Command{
		Use:   "unexpose [domain-name]",
		Short: "Remove a routing rule and close the proxy door",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			local, _ := cmd.Flags().GetBool("local")

			if !local {
				body := map[string]interface{}{
					"domain": domain,
				}
				if _, err := callAPI(cmd, "POST", "/api/v2/proxy/unexpose", body); err == nil {
					fmt.Printf("✅ Unexposed domain %s (via API)\n", domain)
					return
				}
			}

			// Local fallback
			dbPath := store.GetStatePath()
			db, err := store.Open(dbPath)
			if err != nil {
				log.Fatalf("❌ Failed to open database: %v", err)
			}
			defer db.Close()

			mgr := proxy.NewManager(db)
			err = mgr.UnexposeService(domain)
			if err != nil {
				log.Fatalf("❌ Unexpose failed: %v", err)
			}

			fmt.Printf("✅ Unexposed domain %s\n", domain)
		},
	}

	var proxySecureCmd = &cobra.Command{
		Use:   "secure [email-address]",
		Short: "Configure automated Let's Encrypt TLS generation & expose ports 80/443 on the host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			email := args[0]
			local, _ := cmd.Flags().GetBool("local")

			if !local {
				body := map[string]interface{}{
					"email": email,
				}
				if _, err := callAPI(cmd, "POST", "/api/v2/proxy/secure", body); err == nil {
					fmt.Println("✅ Let's Encrypt SSL configured and routing stack redeployed (via API)")
					return
				}
			}

			// Local fallback
			dbPath := store.GetStatePath()
			db, err := store.Open(dbPath)
			if err != nil {
				log.Fatalf("❌ Failed to open database: %v", err)
			}
			defer db.Close()

			mgr := proxy.NewManager(db)
			err = mgr.ConfigureSSL(email)
			if err != nil {
				log.Fatalf("❌ SSL configuration failed: %v", err)
			}

			fmt.Println("✅ Let's Encrypt SSL configured and routing stack redeployed successfully.")
		},
	}

	proxyCmd.AddCommand(proxyStatusCmd, proxyDiscoverCmd, proxyExposeCmd, proxyUnexposeCmd, proxySecureCmd)
	return proxyCmd
}

func printRoutesTable(routes []store.RouteRecord) {
	if len(routes) == 0 {
		fmt.Println("No reverse proxy routes configured.")
		return
	}

	fmt.Printf("\n%-20s %-28s %-6s %-5s %-20s\n", "SERVICE", "DOMAIN", "PORT", "SSL", "MIDDLEWARES")
	fmt.Println(strings.Repeat("-", 84))
	for _, r := range routes {
		sslStr := "No"
		if r.SSL {
			sslStr = "Yes"
		}
		mwStr := r.Middlewares
		if mwStr == "" {
			mwStr = "-"
		}
		fmt.Printf("%-20s %-28s %-6d %-5s %-20s\n", r.Service, r.Domain, r.Port, sslStr, mwStr)
	}
	fmt.Println()
}

func printDiscoverTable(services []proxy.DiscoverableService) {
	if len(services) == 0 {
		fmt.Println("No services discovered.")
		return
	}

	fmt.Printf("\n%-20s %-24s %-8s %-8s %-20s\n", "CONTAINER", "IMAGE", "STATE", "EXPOSED", "PORTS")
	fmt.Println(strings.Repeat("-", 84))
	for _, s := range services {
		expStr := "No"
		if s.Exposed {
			expStr = "Yes (" + s.Domain + ")"
		}

		var portStrs []string
		for _, p := range s.Ports {
			portStrs = append(portStrs, strconv.Itoa(p))
		}
		portsJoin := strings.Join(portStrs, ",")
		if portsJoin == "" {
			portsJoin = "-"
		}

		fmt.Printf("%-20s %-24s %-8s %-8s %-20s\n", s.ContainerName, truncateString(s.Image, 24), s.State, expStr, portsJoin)
	}
	fmt.Println()
}

func truncateString(s string, l int) string {
	if len(s) > l {
		return s[:l-3] + "..."
	}
	return s
}
