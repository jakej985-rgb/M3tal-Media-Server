package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// --- Types ---

type Anomaly struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Reason string `json:"reason"`
}

type Decision struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Reason string `json:"reason"`
	Time   int64  `json:"timestamp"`
}

type HealthReport struct {
	Score     int              `json:"score"`
	Mode      string           `json:"mode"`
	Verdict   string           `json:"verdict"`
	Issues    []string         `json:"issues"`
	Uptime    string           `json:"uptime"`
	Timestamp int64            `json:"timestamp"`
	Agents    map[string]Agent `json:"agent_health"`
}

type Agent struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type M3talState struct {
	dirty      bool
	System     SystemMetrics
	Network    NetworkMetrics
	CPU        float64
	Timestamp  int64
	Storage    StorageStats
	GPU        GpuStats
	Temp       TempStats
	Anomalies  []Anomaly
	Decisions  []Decision
	Health     HealthReport
	Cooldowns  map[string]time.Time
	Containers []ContainerMetric
	MutedUntil int64
	NetworkL   []NetworkRoute

	mu sync.RWMutex
}

type NetworkRoute struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	Image     string `json:"image"`
	Icon      string `json:"icon"`
	Container string `json:"container"`
}

type StorageStats struct {
	Disks     map[string]DiskInfo `json:"disks"`
	IO        *DiskIO             `json:"io"`
	Timestamp int64               `json:"timestamp"`
	Status    string              `json:"status"`
}

type DiskInfo struct {
	Free    string  `json:"free"`
	Temp    float64 `json:"temp"`
	Percent float64 `json:"percent"`
}

type DiskIO struct {
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
}

type GpuStats struct {
	Name     string  `json:"name"`
	Temp     float64 `json:"temp"`
	Load     int     `json:"load"`
	MemUsed  int     `json:"mem_used"`
	MemTotal int     `json:"mem_total"`
	Active   bool    `json:"active"`
}

type TempStats struct {
	CPUTemp   float64 `json:"cpu_temp"`
	GPUTemp   float64 `json:"gpu_temp"`
	Timestamp int64   `json:"timestamp"`
	Status    string  `json:"status"`
}

type SystemMetrics struct {
	CPU       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	MemGB     float64 `json:"mem_gb"`
	MemTotal  float64 `json:"mem_total"`
	Timestamp int64   `json:"timestamp"`
}

type ContainerMetric struct {
	Name     string  `json:"name"`
	CPU      float64 `json:"cpu"`
	Mem      float64 `json:"mem"`
	MemUsage uint64  `json:"mem_usage"`
	MemLimit uint64  `json:"mem_limit"`
	Status       string  `json:"status"`
	State        string  `json:"state"`
	HealthStatus string  `json:"health_status"` // "healthy", "unhealthy", "starting", or ""
	Managed      bool    `json:"managed"`
	NetRx        uint64  `json:"net_rx"`
	NetTx        uint64  `json:"net_tx"`
}

type NetworkMetrics struct {
	Down    string  `json:"down"`
	Up      string  `json:"up"`
	DownRaw float64 `json:"down_raw"`
	UpRaw   float64 `json:"up_raw"`
	Load    float64 `json:"load"`
}

// --- Persistence ---

func (s *M3talState) Save(stateDir string) {
	s.mu.Lock()
	if !s.dirty {
		s.mu.Unlock()
		return
	}
	s.dirty = false
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()

	saveAtomic(filepath.Join(stateDir, "metrics.json"), s.System)
	saveAtomic(filepath.Join(stateDir, "storage.json"), s.Storage)
	saveAtomic(filepath.Join(stateDir, "gpu.json"), s.GPU)
	saveAtomic(filepath.Join(stateDir, "temp.json"), s.Temp)
	saveAtomic(filepath.Join(stateDir, "network.json"), map[string]interface{}{"metrics": s.Network, "links": s.NetworkL})
	saveAtomic(filepath.Join(stateDir, "anomalies.json"), map[string]interface{}{"issues": s.Anomalies})
	saveAtomic(filepath.Join(stateDir, "decisions.json"), map[string]interface{}{"actions": s.Decisions})
	saveAtomic(filepath.Join(stateDir, "health.json"), map[string]interface{}{
		"status":     "online",
		"containers": s.Containers,
		"timestamp":  s.Timestamp,
	})
	saveAtomic(filepath.Join(stateDir, "health_report.json"), s.Health)
	
	// Unified state for dashboard and API discovery
	saveAtomic(filepath.Join(stateDir, "system.json"), s)
}

func saveAtomic(path string, data interface{}) {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return
	}
	f.Close()
	os.Rename(tmp, path)
}

// --- Main ---

func main() {
	state := &M3talState{
		Cooldowns:  make(map[string]time.Time),
		Containers: []ContainerMetric{},
	}

	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		stateDir = filepath.Join("..", "..", "state")
	}

	// Setup logging
	logsDir := filepath.Join(stateDir, "logs")
	os.MkdirAll(logsDir, 0755)
	if logFile, err := os.OpenFile(filepath.Join(logsDir, "m3tal-core.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("🚀 M3TAL Core (The Brain) starting...")

	// Launch 6 Core Pillars (Agents)
	go registryAgent(ctx, state, stateDir)
	go monitorAgent(ctx, state)
	go metricsAgent(state)
	go anomalyAgent(ctx, state)
	go decisionAgent(ctx, state)
	go reconcileAgent(ctx, state)

	select {}
}

// --- Helpers ---

func getAgentLogger(name string) *log.Logger {
	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		stateDir = filepath.Join("..", "..", "state")
	}
	logsDir := filepath.Join(stateDir, "logs")
	os.MkdirAll(logsDir, 0755)
	logFile, err := os.OpenFile(filepath.Join(logsDir, name+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return log.New(os.Stdout, "["+name+"] ", log.LstdFlags)
	}
	return log.New(io.MultiWriter(os.Stdout, logFile), "["+name+"] ", log.LstdFlags)
}

func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec < 1024 {
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	} else if bytesPerSec < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", bytesPerSec/1024)
	} else if bytesPerSec < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB/s", bytesPerSec/(1024*1024))
	}
	return fmt.Sprintf("%.1f GB/s", bytesPerSec/(1024*1024*1024))
}

func getUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "—"
	}
	parts := strings.Fields(string(data))
	if len(parts) == 0 {
		return "—"
	}
	sec, _ := strconv.ParseFloat(parts[0], 64)
	days := int(sec / 86400)
	hours := int(sec) / 3600 % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	return fmt.Sprintf("%dh", hours)
}

// --- Core Agents (6 Pillars) ---

func registryAgent(ctx context.Context, s *M3talState, stateDir string) {
	logger := getAgentLogger("registry")
	logger.Println("Agent started")
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.Save(stateDir)
		}
	}
}

func monitorAgent(ctx context.Context, s *M3talState) {
	logger := getAgentLogger("monitor")
	logger.Println("Agent started")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Printf("❌ Failed to init Docker client: %v", err)
		return
	}
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		res, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true})
		if err != nil {
			continue
		}
		newStats := []ContainerMetric{}
		for _, c := range res.Items {
			name := "unknown"
			if len(c.Names) > 0 {
				name = c.Names[0][1:]
			}
			
			cpuPerc := 0.0
			memPerc := 0.0
			var memUsage uint64
			var memLimit uint64
			var netRx uint64
			var netTx uint64

			if strings.ToLower(string(c.State)) == "running" {
				stats, err := cli.ContainerStats(ctx, c.ID, client.ContainerStatsOptions{Stream: false})
				if err == nil {
					var v struct {
						CPUStats struct {
							CPUUsage struct {
								TotalUsage uint64 `json:"total_usage"`
							} `json:"cpu_usage"`
							SystemCPUUsage uint64 `json:"system_cpu_usage"`
							OnlineCPUs     uint64 `json:"online_cpus"`
						} `json:"cpu_stats"`
						PreCPUStats struct {
							CPUUsage struct {
								TotalUsage uint64 `json:"total_usage"`
							} `json:"cpu_usage"`
							SystemCPUUsage uint64 `json:"system_cpu_usage"`
						} `json:"precpu_stats"`
						MemoryStats struct {
							Usage uint64            `json:"usage"`
							Limit uint64            `json:"limit"`
							Stats map[string]uint64 `json:"stats"`
						} `json:"memory_stats"`
						Networks map[string]struct {
							RxBytes uint64 `json:"rx_bytes"`
							TxBytes uint64 `json:"tx_bytes"`
						} `json:"networks"`
					}
					if err := json.NewDecoder(stats.Body).Decode(&v); err == nil {
						cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
						sysDelta := float64(v.CPUStats.SystemCPUUsage) - float64(v.PreCPUStats.SystemCPUUsage)
						cpus := float64(v.CPUStats.OnlineCPUs)
						if cpus == 0 { cpus = 1 }
						if sysDelta > 0 && cpuDelta > 0 {
							cpuPerc = (cpuDelta / sysDelta) * cpus * 100.0
						}
						if v.MemoryStats.Limit > 0 {
							cache := v.MemoryStats.Stats["inactive_file"]
							if cache == 0 { cache = v.MemoryStats.Stats["cache"] }
							memUsage = v.MemoryStats.Usage - cache
							memLimit = v.MemoryStats.Limit
							memPerc = (float64(memUsage) / float64(memLimit)) * 100.0
						}
						for _, nw := range v.Networks {
							netRx += nw.RxBytes
							netTx += nw.TxBytes
						}
					}
					stats.Body.Close()
				}
			}

			managed := false
			if _, ok := c.Labels["m3tal.managed"]; ok {
				managed = true
			} else if _, ok := c.Labels["m3tal.stack"]; ok {
				managed = true
			}
			health := ""
			if c.State == "running" && c.Health != nil {
				health = c.Health.Status
			}

			newStats = append(newStats, ContainerMetric{
				Name:         name,
				CPU:          cpuPerc,
				Mem:          memPerc,
				MemUsage:     memUsage,
				MemLimit:     memLimit,
				Status:       c.Status,
				State:        string(c.State),
				HealthStatus: health,
				Managed:      managed,
				NetRx:        netRx,
				NetTx:        netTx,
			})
		}
		s.mu.Lock()
		s.Containers = newStats
		s.dirty = true
		s.mu.Unlock()
	}
}

func metricsAgent(s *M3talState) {
	go hostMetricsLoop(s)
	go gpuMetricsLoop(s)
	go storageMetricsLoop(s)
}

func hostMetricsLoop(s *M3talState) {
	logger := getAgentLogger("metrics")
	logger.Println("Host metrics loop started")
	ticker := time.NewTicker(2 * time.Second)
	var lastRecv, lastSent uint64
	var lastTime time.Time

	for range ticker.C {
		now := time.Now()
		cpuPerc, _ := cpu.Percent(0, false)
		vm, _ := mem.VirtualMemory()
		netIO, _ := net.IOCounters(false)

		cpuVal := 0.0
		if len(cpuPerc) > 0 {
			cpuVal = cpuPerc[0]
		}

		var totalRecv, totalSent uint64
		for _, io := range netIO {
			totalRecv += io.BytesRecv
			totalSent += io.BytesSent
		}

		var down, up, load float64
		if !lastTime.IsZero() {
			dt := now.Sub(lastTime).Seconds()
			if dt > 0 {
				down = float64(totalRecv-lastRecv) / (1024 * 1024) / dt
				up = float64(totalSent-lastSent) / (1024 * 1024) / dt
				capacity := 125.0
				load = ((down + up) / capacity) * 100
				if load > 100 { load = 100 }
			}
		}
		lastRecv = totalRecv
		lastSent = totalSent
		lastTime = now

		s.mu.Lock()
		s.System = SystemMetrics{
			CPU:       cpuVal,
			Mem:       vm.UsedPercent,
			MemGB:     float64(vm.Used) / (1024 * 1024 * 1024),
			MemTotal:  float64(vm.Total) / (1024 * 1024 * 1024),
			Timestamp: now.Unix(),
		}
		s.Network = NetworkMetrics{
			Down:    formatSpeed(down * 1024 * 1024),
			Up:      formatSpeed(up * 1024 * 1024),
			DownRaw: down,
			UpRaw:   up,
			Load:    load,
		}
		s.CPU = cpuVal
		s.Timestamp = now.Unix()
		s.dirty = true
		s.mu.Unlock()
	}
}

func gpuMetricsLoop(s *M3talState) {
	logger := getAgentLogger("metrics")
	logger.Println("GPU metrics loop started")
	ticker := time.NewTicker(10 * time.Second)
	gpuRe := regexp.MustCompile(`gpu\s+([\d\.]+)%`)
	vramRe := regexp.MustCompile(`vram\s+[\d\.]+% ([\d\.]+)mb`)

	for range ticker.C {
		var cpuT, gpuT float64
		var gpuStats GpuStats

		if entries, err := os.ReadDir("/sys/class/hwmon"); err == nil {
			for _, entry := range entries {
				namePath := filepath.Join("/sys/class/hwmon", entry.Name(), "name")
				tempPath := filepath.Join("/sys/class/hwmon", entry.Name(), "temp1_input")
				if name, err := os.ReadFile(namePath); err == nil {
					n := strings.TrimSpace(string(name))
					if data, err := os.ReadFile(tempPath); err == nil {
						if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
							t := float64(val) / 1000.0
							switch n {
							case "amdgpu", "radeon":
								gpuT = t
							case "coretemp", "cpu_thermal", "k10temp", "zenpower", "acpitz", "it87":
								if t > cpuT { cpuT = t }
							default:
								if strings.Contains(strings.ToLower(n), "temp") && t > cpuT { cpuT = t }
							}
						}
					}
				}
			}
		}

		if cpuT == 0 {
			for i := 0; i < 5; i++ {
				path := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", i)
				if data, err := os.ReadFile(path); err == nil {
					if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
						cpuT = float64(val) / 1000.0
						break
					}
				}
			}
		}

		gpuStats.Name = "AMD Radeon HD 5770"
		gpuStats.Temp = gpuT
		gpuStats.MemTotal = 1024

		radeontopPath := "/usr/bin/radeontop"
		cmd := exec.Command(radeontopPath, "-d", "-", "-l", "1")
		if output, err := cmd.CombinedOutput(); err == nil {
			line := string(output)
			gpuStats.Active = true
			if m := gpuRe.FindStringSubmatch(line); len(m) > 1 {
				if f, err := strconv.ParseFloat(m[1], 64); err == nil { gpuStats.Load = int(f) }
			}
			if m := vramRe.FindStringSubmatch(line); len(m) > 1 {
				if f, err := strconv.ParseFloat(m[1], 64); err == nil { gpuStats.MemUsed = int(f) }
			}
		}

		status := "healthy"
		if cpuT > 85 || gpuT > 85 { status = "critical" } else if cpuT > 75 || gpuT > 75 { status = "warning" }

		s.mu.Lock()
		s.Temp = TempStats{CPUTemp: cpuT, GPUTemp: gpuT, Timestamp: time.Now().Unix(), Status: status}
		s.GPU = gpuStats
		s.dirty = true
		s.mu.Unlock()
	}
}

func storageMetricsLoop(s *M3talState) {
	logger := getAgentLogger("metrics")
	logger.Println("Storage metrics loop started")
	ticker := time.NewTicker(30 * time.Second)
	tempKeywords := []string{"Temperature_Celsius", "Airflow_Temperature_Cel", "Composite Temperature", "Current Drive Temperature:"}

	for range ticker.C {
		disksMap := make(map[string]DiskInfo)
		hostRoot := "/host"
		if _, err := os.Stat(hostRoot); err != nil { hostRoot = "" }

		var mountEntries []struct { Device, Mountpoint string }
		if data, err := os.ReadFile("/proc/mounts"); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				fields := strings.Fields(line)
				if len(fields) < 3 { continue }
				dev, mnt, fstype := fields[0], fields[1], fields[2]
				if !strings.HasPrefix(dev, "/dev/") { continue }
				if fstype == "squashfs" || fstype == "tmpfs" || fstype == "devtmpfs" { continue }
				mountEntries = append(mountEntries, struct{ Device, Mountpoint string }{dev, mnt})
			}
		}

		seen := make(map[string]bool)
		for _, entry := range mountEntries {
			if seen[entry.Device] { continue }
			seen[entry.Device] = true
			usage, err := disk.Usage(entry.Mountpoint)
			if err != nil || usage.Total == 0 { continue }
			if usage.Total < 1024*1024*1024 { continue }

			var driveT float64
			phys := entry.Device
			cmd := exec.Command("chroot", "/host", "smartctl", "-a", phys)
			if output, err := cmd.CombinedOutput(); err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					matched := false
					for _, kw := range tempKeywords {
						if strings.Contains(strings.ToLower(line), strings.ToLower(kw)) { matched = true; break }
					}
					if matched {
						fields := strings.Fields(line)
						if len(fields) > 0 {
							for i := len(fields) - 1; i >= 0; i-- {
								f := regexp.MustCompile(`[^\d].*`).ReplaceAllString(fields[i], "")
								if val, err := strconv.ParseFloat(f, 64); err == nil && val > 0 && val < 150 {
									driveT = val; break
								}
							}
						}
					}
				}
			}
			disksMap[entry.Device] = DiskInfo{Free: formatSpeed(float64(usage.Free)), Temp: driveT, Percent: usage.UsedPercent}
		}
		s.mu.Lock()
		s.Storage = StorageStats{Disks: disksMap, Timestamp: time.Now().Unix(), Status: "healthy"}
		s.dirty = true
		s.mu.Unlock()
	}
}

func anomalyAgent(ctx context.Context, s *M3talState) {
	logger := getAgentLogger("anomaly")
	logger.Println("Agent started")
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ctx.Done(): return
		case <-ticker.C:
			res, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true})
			if err != nil { continue }
			var anomalies []Anomaly
			for _, c := range res.Items {
				if strings.ToLower(string(c.State)) == "exited" {
					anomalies = append(anomalies, Anomaly{Type: "crash", Target: c.Names[0][1:], Reason: c.Status})
				}
			}
			s.mu.Lock()
			s.Anomalies = anomalies
			s.dirty = true
			s.mu.Unlock()
		}
	}
}

func decisionAgent(ctx context.Context, s *M3talState) {
	logger := getAgentLogger("decision")
	logger.Println("Agent started")
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done(): return
		case <-ticker.C:
			s.mu.Lock()
			score := 100
			var issues []string

			// --- 1. Thermal Stress Check (-15 pts each) ---
			if s.Temp.CPUTemp > 80 {
				score -= 15
				issues = append(issues, fmt.Sprintf("High CPU Temp: %.1f°C", s.Temp.CPUTemp))
			}
			if s.Temp.GPUTemp > 80 {
				score -= 15
				issues = append(issues, fmt.Sprintf("High GPU Temp: %.1f°C", s.Temp.GPUTemp))
			}

			// --- 2. Resource Saturation (-10 pts each) ---
			if s.System.CPU > 90 {
				score -= 10
				issues = append(issues, "CPU Saturation (>90%)")
			}
			if s.System.Mem > 90 {
				score -= 10
				issues = append(issues, "Memory Pressure (>90%)")
			}

			// --- 3. Service Instability (-20 pts each) ---
			unhealthyCount := 0
			crashedCount := 0
			for _, c := range s.Containers {
				if !c.Managed { continue }
				if c.HealthStatus == "unhealthy" {
					unhealthyCount++
				}
				if c.State == "exited" || c.State == "dead" {
					crashedCount++
				}
			}
			if unhealthyCount > 0 {
				penalty := unhealthyCount * 20
				score -= penalty
				issues = append(issues, fmt.Sprintf("%d Unhealthy Services", unhealthyCount))
			}
			if crashedCount > 0 {
				penalty := crashedCount * 20
				score -= penalty
				issues = append(issues, fmt.Sprintf("%d Crashed Services", crashedCount))
			}

			// --- 4. Disk Pressure (-5 pts each) ---
			for dev, info := range s.Storage.Disks {
				if info.Percent > 95 {
					score -= 5
					issues = append(issues, fmt.Sprintf("Disk Space Critical: %s", dev))
				}
			}

			// Finalize Score
			if score < 0 { score = 0 }
			
			mode := "OPTIMAL"
			if score < 40 {
				mode = "CRITICAL"
			} else if score < 70 {
				mode = "DEGRADED"
			} else if score < 90 {
				mode = "STABLE"
			}
			
			s.Health = HealthReport{
				Score:     score,
				Mode:      mode,
				Verdict:   mode,
				Issues:    issues,
				Uptime:    getUptime(),
				Timestamp: time.Now().Unix(),
				Agents:    map[string]Agent{"decision": {Status: "healthy", Timestamp: time.Now().Unix()}},
			}
			s.dirty = true
			s.mu.Unlock()
		}
	}
}

func reconcileAgent(ctx context.Context, s *M3talState) {
	logger := getAgentLogger("reconcile")
	logger.Println("Agent started")
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ticker := time.NewTicker(10 * time.Second)
	
	// Track recovery attempts to prevent restart loops
	recoveries := make(map[string]int)
	lastRecovery := make(map[string]time.Time)

	for range ticker.C {
		s.mu.RLock()
		containers := s.Containers
		s.mu.RUnlock()
		
		for _, c := range containers {
			if !c.Managed { continue }

			shouldHeal := false
			reason := ""

			// Rule 1: Service is Down
			if c.State == "exited" || c.State == "dead" {
				shouldHeal = true
				reason = "container_down"
			}
			
			// Rule 2: Zombie Service (Running but Unhealthy)
			if c.State == "running" && c.HealthStatus == "unhealthy" {
				shouldHeal = true
				reason = "zombie_unhealthy"
			}

			if shouldHeal {
				// Cooldown Check: Max 3 restarts per hour per service
				if time.Since(lastRecovery[c.Name]) < 1*time.Hour && recoveries[c.Name] >= 3 {
					logger.Printf("⚠️ RECONCILE: Skipping %s recovery (Rate Limit reached). Possible fatal loop.", c.Name)
					continue
				}

				logger.Printf("🛡️ RECONCILE: Healing %s (Reason: %s)", c.Name, reason)
				
				if err := cli.ContainerRestart(ctx, c.Name, client.ContainerRestartOptions{}); err != nil {
					logger.Printf("❌ RECONCILE: Failed to heal %s: %v", c.Name, err)
				} else {
					logger.Printf("✅ RECONCILE: Successfully healed %s.", c.Name)
					
					// Update recovery tracking
					if time.Since(lastRecovery[c.Name]) > 1*time.Hour {
						recoveries[c.Name] = 1
					} else {
						recoveries[c.Name]++
					}
					lastRecovery[c.Name] = time.Now()
				}
			}
		}
	}
}
