package output

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/jakej985-rgb/m3tal-core/internal/registry"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
)

// PrintContainersTable formats a list of containers as a table to standard output.
func PrintContainersTable(containers []models.Container) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "CONTAINER ID\tNAMES\tIMAGE\tSTATE\tSTATUS")
	fmt.Fprintln(w, "------------\t-----\t-----\t-----\t------")
	for _, c := range containers {
		id := c.ID
		if len(id) > 12 {
			id = id[:12]
		}
		names := strings.Join(c.Names, ",")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", id, names, c.Image, c.State, c.Status)
	}
	w.Flush()
}

// PrintStacksTable formats a list of stacks as a table to standard output.
func PrintStacksTable(stacks []models.Stack) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tPATH\tSERVICES")
	fmt.Fprintln(w, "----\t------\t----\t--------")
	for _, s := range stacks {
		svcs := strings.Join(s.Services, ",")
		if svcs == "" {
			svcs = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.Status, s.ComposePath, svcs)
	}
	w.Flush()
}

// PrintStats prints system metrics in a styled dashboard layout.
func PrintStats(s *models.MetricsResponse) {
	fmt.Printf("💻 Hostname: %s\n", s.Hostname)
	uptimeDur := time.Duration(s.Uptime) * time.Second
	fmt.Printf("⏱️  Uptime:   %s\n", uptimeDur.String())
	fmt.Printf("🧠 Memory:   %.1f%%\n", s.MemoryUsage)
	fmt.Printf("💿 Disk:     %.1f%%\n", s.DiskUsage)
	fmt.Printf("🔥 CPU:      %.1f%%\n", s.CPUUsage)
}

// PrintRoutesTable displays a list of active reverse proxy routes.
func PrintRoutesTable(routes []models.RouteRecord) {
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

// PrintDiscoverTable prints discoverable candidate containers.
func PrintDiscoverTable(services []models.DiscoverableService) {
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

// PrintVPNStatus formats and prints the VPN status structure.
func PrintVPNStatus(status *models.VPNStatus) {
	fmt.Println("🌐 VPN Connection Status:")
	fmt.Println("----------------------------------------")
	if status.Connected {
		fmt.Println("🟢 Status:      Connected (Running)")
		fmt.Printf("🔒 Provider:    %s\n", status.Provider)
		fmt.Printf("🌍 Region:      %s\n", status.Region)
		fmt.Printf("📬 External IP: %s\n", status.ExternalIP)
		if status.ForwardedPort > 0 {
			fmt.Printf("🔌 Port:        %d (Forwarded)\n", status.ForwardedPort)
		}
	} else {
		fmt.Printf("🔴 Status:      Disconnected (%s)\n", status.StatusText)
	}
}

// PrintVPNLeakReport formats and prints the leak check report.
func PrintVPNLeakReport(report *models.VPNLeakReport) {
	fmt.Println("\n🛡️  VPN Leak Detection Report:")
	fmt.Println("----------------------------------------")
	fmt.Printf("🏠 Host Public IP: %s\n", report.HostIP)
	fmt.Printf("🔒 VPN Outbound IP: %s\n", report.VPNIP)

	if report.Leak {
		fmt.Println("🚨 RESULT: LEAK DETECTED! Your traffic is NOT protected!")
		if len(report.StoppedContainers) > 0 {
			fmt.Println("⚠️  Activating kill switch (stopping dependent containers)...")
			fmt.Printf("🛑 Successfully stopped containers: %s\n", strings.Join(report.StoppedContainers, ", "))
		} else {
			fmt.Println("✅ No active dependent containers found running.")
		}
	} else {
		fmt.Println("✅ RESULT: SAFE. Outbound IP is protected by VPN.")
	}
}

// PrintRegistrySearchTable formats registry search results.
func PrintRegistrySearchTable(stacks []registry.StackMetadata) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tCATEGORY\tVERSION\tDESCRIPTION")
	for _, s := range stacks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.Category, s.Version, s.Description)
	}
	w.Flush()
}

// FatalError prints the error to stderr and calls log.Fatalf.
func FatalError(err error) {
	log.Fatalf("❌ Error: %v", err)
}

// FatalErrorMsg prints a formatted error message and exits.
func FatalErrorMsg(msg string, args ...any) {
	log.Fatalf("❌ "+msg, args...)
}

// PrintError prints the error to stderr without exiting.
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
}

func truncateString(s string, l int) string {
	if len(s) > l {
		return s[:l-3] + "..."
	}
	return s
}
