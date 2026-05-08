package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

// M3talState holds the unified state of the control plane in memory
type M3talState struct {
	mu         sync.RWMutex
	System     SystemMetrics     `json:"system"`
	Containers []ContainerMetric `json:"containers"`
	Network    NetworkMetrics    `json:"network"`
	Anomalies  []Anomaly         `json:"issues"`    // For anomalies.json
	Decisions  []Decision        `json:"actions"`   // For decisions.json
	MutedUntil int64             `json:"muted_until"` // For mute_state.json
	GPU        GpuStats          `json:"gpu"`        // For gpu.json
	Temp       TempStats         `json:"temp"`       // For temp.json
	Storage    StorageStats      `json:"storage"`    // For storage.json
	NetworkL   []NetworkRoute    `json:"links"`      // For network.json
	Timestamp  int64             `json:"timestamp"`
	CPU        float64           `json:"cpu"`
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

type Decision struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Reason string `json:"reason"`
	Time   int64  `json:"timestamp"`
}

type Anomaly struct {
// ... existing structs ...
	Type   string `json:"type"`
	Target string `json:"target"`
	Reason string `json:"reason"`
}

type SystemMetrics struct {
	CPU       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	MemGB     float64 `json:"mem_gb"`
	Timestamp int64   `json:"timestamp"`
}

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From struct {
			ID int64 `json:"id"`
		} `json:"from"`
		Text string `json:"text"`
	} `json:"message"`
}

type ContainerMetric struct {
	Name   string  `json:"name"`
	CPU    float64 `json:"cpu"`
	Mem    float64 `json:"mem"`
	Status string  `json:"status"` // e.g. "running"
	State  string  `json:"state"`  // e.g. "exited"
}

type NetworkMetrics struct {
	Down string  `json:"down"`
	Up   string  `json:"up"`
	Load float64 `json:"load"`
}

func main() {
	state := &M3talState{
		Containers: []ContainerMetric{},
	}

	// Resolve state path
	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		stateDir = filepath.Join("..", "state")
	}
	metricsPath := filepath.Join(stateDir, "metrics.json")
	anomalyPath := filepath.Join(stateDir, "anomalies.json")
	decisionPath := filepath.Join(stateDir, "decisions.json")
	gpuPath := filepath.Join(stateDir, "gpu.json")
	tempPath := filepath.Join(stateDir, "temp.json")
	storagePath := filepath.Join(stateDir, "storage.json")
	networkPath := filepath.Join(stateDir, "network.json")

	fmt.Println("🚀 M3TAL Go Backend (Linux-Ready) starting...")
	fmt.Printf("📂 Target State: %s\n", metricsPath)

	ctx := context.Background()

	// 1. Metrics Agent (System)
	go metricsAgent(state)

	// 2. Docker Agent (Containers)
	go dockerAgent(ctx, state)

	// 3. Anomaly Agent (Detection)
	go anomalyAgent(state)

	// 4. Healer Agent (Auto-Recovery)
	go healerAgent(ctx, state)

	// 5. Notify Agent (Telegram)
	go notifyAgent(state)

	// 6. Listener Agent (Interactive)
	go listenerAgent(ctx, state)

	// 7. Hardware Agent (GPU/Temp)
	go hardwareAgent(state)

	// 8. Storage Agent (Disk/IO)
	go storageAgent(state)

	// 9. Scout Agent (Network Routes)
	go scoutAgent(ctx, state)

	// 10. Persistence Agent (JSON Sync)
	go saveAgent(state, metricsPath, anomalyPath, decisionPath, gpuPath, tempPath, storagePath, networkPath)

	// Keep alive
	select {}
}

func metricsAgent(s *M3talState) {
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

		var down, up, load float64
		if len(netIO) > 0 {
			if !lastTime.IsZero() {
				dt := now.Sub(lastTime).Seconds()
				if dt > 0 {
					down = float64(netIO[0].BytesRecv-lastRecv) / (1024 * 1024) / dt
					up = float64(netIO[0].BytesSent-lastSent) / (1024 * 1024) / dt
					
					// Assuming 1Gbps (125MB/s) capacity for load calculation
					capacity := 125.0 
					load = ((down + up) / capacity) * 100
					if load > 100 {
						load = 100
					}
				}
			}
			lastRecv = netIO[0].BytesRecv
			lastSent = netIO[0].BytesSent
			lastTime = now
		}

		s.mu.Lock()
		s.System = SystemMetrics{
			CPU:       cpuVal,
			Mem:       vm.UsedPercent,
			MemGB:     float64(vm.Used) / (1024 * 1024 * 1024),
			Timestamp: now.Unix(),
		}
		s.Network = NetworkMetrics{
			Down: formatSpeed(down * 1024 * 1024), // down is in MB/s currently, convert back to bytes for formatter
			Up:   formatSpeed(up * 1024 * 1024),
			Load: load,
		}
		s.CPU = cpuVal
		s.Timestamp = now.Unix()
		s.mu.Unlock()
	}
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

func dockerAgent(ctx context.Context, s *M3talState) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("⚠️ Docker SDK Error: %v", err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		// Latest Moby uses client-specific options and returns a Result struct
		res, err := cli.ContainerList(ctx, client.ContainerListOptions{})
		if err != nil {
			log.Printf("⚠️ Failed to list containers: %v", err)
			continue
		}

		newStats := []ContainerMetric{}
		// Range over the Items slice in the ContainerListResult
		for _, c := range res.Items {
			name := "unknown"
			if len(c.Names) > 0 {
				name = c.Names[0][1:] // Remove leading slash
			}
			newStats = append(newStats, ContainerMetric{
				Name:   name,
				CPU:    0.0,
				Mem:    0.0,
				Status: c.Status,
				State:  string(c.State),
			})
		}

		s.mu.Lock()
		s.Containers = newStats
		s.mu.Unlock()
	}
}

func anomalyAgent(s *M3talState) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		issues := []Anomaly{}

		s.mu.RLock()
		cpu := s.System.CPU
		mem := s.System.Mem
		containers := s.Containers
		s.mu.RUnlock()

		// 1. Host Resource Checks
		if cpu > 90 {
			issues = append(issues, Anomaly{Type: "transient", Target: "host", Reason: fmt.Sprintf("CPU saturation: %.1f%%", cpu)})
		}
		if mem > 95 {
			issues = append(issues, Anomaly{Type: "critical", Target: "host", Reason: fmt.Sprintf("Memory saturation: %.1f%%", mem)})
		}

		// 2. Container Resource Checks
		for _, c := range containers {
			if c.CPU > 90 {
				issues = append(issues, Anomaly{Type: "resource_spike", Target: c.Name, Reason: fmt.Sprintf("High Container CPU: %.1f%%", c.CPU)})
			}
		}

		s.mu.Lock()
		s.Anomalies = issues
		s.mu.Unlock()
	}
}

func healerAgent(ctx context.Context, s *M3talState) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		s.mu.RLock()
		containers := s.Containers
		s.mu.RUnlock()

		for _, c := range containers {
			// If container is exited or dead, try to restart it
			if c.State == "exited" || c.State == "dead" {
				log.Printf("🛡️ HEALER: Detected crashed container: %s. Restarting...", c.Name)
				
				_, err := cli.ContainerRestart(ctx, c.Name, client.ContainerRestartOptions{})
				if err != nil {
					log.Printf("❌ HEALER: Failed to restart %s: %v", c.Name, err)
					continue
				}

				// Log the decision
				s.mu.Lock()
				s.Decisions = append(s.Decisions, Decision{
					Type:   "restart",
					Target: c.Name,
					Reason: fmt.Sprintf("State was '%s'", c.State),
					Time:   time.Now().Unix(),
				})
				// Limit log size to last 50 actions
				if len(s.Decisions) > 50 {
					s.Decisions = s.Decisions[1:]
				}
				s.mu.Unlock()
			}
		}
	}
}

func hardwareAgent(s *M3talState) {
	ticker := time.NewTicker(10 * time.Second)
	gpuRe := regexp.MustCompile(`gpu\s+([\d\.]+)%`)
	vramRe := regexp.MustCompile(`vram\s+[\d\.]+% ([\d\.]+)mb`)

	for range ticker.C {
		var cpuT, gpuT float64
		var gpuStats GpuStats

		// 1. CPU Temp (Multi-zone Scan)
		for i := 0; i < 6; i++ {
			path := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", i)
			if data, err := os.ReadFile(path); err == nil {
				if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
					cpuT = float64(val) / 1000.0
					if cpuT > 20 && cpuT < 110 { // Basic sanity check for CPU temp
						break
					}
				}
			}
		}

		// 2. AMD GPU Temp (Linux Fallback)
		paths := []string{
			"/sys/class/drm/card0/device/hwmon/hwmon0/temp1_input",
			"/sys/class/drm/card0/device/hwmon/hwmon1/temp1_input",
			"/sys/class/drm/card0/device/hwmon/hwmon2/temp1_input",
		}
		for _, p := range paths {
			if data, err := os.ReadFile(p); err == nil {
				if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
					gpuT = float64(val) / 1000.0
					break
				}
			}
		}

		// 3. Radeontop Stats
		cmd := exec.Command("radeontop", "-d", "-", "-l", "1")
		if output, err := cmd.CombinedOutput(); err == nil {
			line := string(output)
			gpuStats.Active = true
			gpuStats.Name = "AMD Radeon HD 5770"
			gpuStats.Temp = gpuT
			gpuStats.MemTotal = 1024

			if m := gpuRe.FindStringSubmatch(line); len(m) > 1 {
				if f, err := strconv.ParseFloat(m[1], 64); err == nil {
					gpuStats.Load = int(f)
				}
			}
			if m := vramRe.FindStringSubmatch(line); len(m) > 1 {
				if f, err := strconv.ParseFloat(m[1], 64); err == nil {
					gpuStats.MemUsed = int(f)
				}
			}
		}

		status := "healthy"
		if cpuT > 85 || gpuT > 85 {
			status = "critical"
		} else if cpuT > 75 || gpuT > 75 {
			status = "warning"
		}

		s.mu.Lock()
		s.Temp = TempStats{
			CPUTemp:   cpuT,
			GPUTemp:   gpuT,
			Timestamp: time.Now().Unix(),
			Status:    status,
		}
		s.GPU = gpuStats
		s.mu.Unlock()
	}
}

func storageAgent(s *M3talState) {
	ticker := time.NewTicker(30 * time.Second)
	smartRe := regexp.MustCompile(`(?i)(?:Temperature_Celsius|Airflow_Temperature_Cel|Composite\s+Temperature|Current\s+Drive\s+Temperature:).*?(\d+)`)

	for range ticker.C {
		disks := make(map[string]DiskInfo)
		highestUsage := 0.0

		// 1. Get partitions and usage
		parts, _ := disk.Partitions(false)
		for _, p := range parts {
			// Skip special filesystems
			if strings.HasPrefix(p.Mountpoint, "/proc") || strings.HasPrefix(p.Mountpoint, "/dev") || strings.HasPrefix(p.Mountpoint, "/sys") {
				continue
			}

			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}

			label := p.Mountpoint
			if label == "/" {
				label = "System"
			} else {
				label = filepath.Base(label)
			}

			// 2. Drive Temperature (smartctl)
			var driveT float64
			dev := p.Device
			if strings.HasPrefix(dev, "/dev/") {
				// Try to get physical device (e.g. /dev/sda1 -> /dev/sda)
				phys := dev
				if len(dev) > 8 && (dev[7] >= '0' && dev[7] <= '9') {
					phys = dev[:7]
				}
				cmd := exec.Command("smartctl", "-a", phys)
				if output, err := cmd.CombinedOutput(); err == nil {
					if m := smartRe.FindStringSubmatch(string(output)); len(m) > 1 {
						if f, err := strconv.ParseFloat(m[1], 64); err == nil {
							driveT = f
						}
					}
				}
			}

			disks[label] = DiskInfo{
				Free:    fmt.Sprintf("%.1f", float64(usage.Free)/(1024*1024*1024)),
				Temp:    driveT,
				Percent: usage.UsedPercent,
			}
			if usage.UsedPercent > highestUsage {
				highestUsage = usage.UsedPercent
			}
		}

		// 3. IO Counters
		var ioStats *DiskIO
		if io, err := disk.IOCounters(); err == nil {
			// Sum all disks for global IO
			var rC, wC, rB, wB uint64
			for _, stats := range io {
				rC += stats.ReadCount
				wC += stats.WriteCount
				rB += stats.ReadBytes
				wB += stats.WriteBytes
			}
			ioStats = &DiskIO{
				ReadCount:  rC,
				WriteCount: wC,
				ReadBytes:  rB,
				WriteBytes: wB,
			}
		}

		status := "healthy"
		if highestUsage > 95 {
			status = "critical"
		} else if highestUsage > 85 {
			status = "warning"
		}

		s.mu.Lock()
		s.Storage = StorageStats{
			Disks:     disks,
			IO:        ioStats,
			Timestamp: time.Now().Unix(),
			Status:    status,
		}
		s.mu.Unlock()
	}
}

func scoutAgent(ctx context.Context, s *M3talState) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ticker := time.NewTicker(60 * time.Second)
	hostRe := regexp.MustCompile(`Host\(` + "`" + `([^` + "`" + `]+)` + "`" + `\)`)
	
	hc := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for range ticker.C {
		res, err := cli.ContainerList(ctx, client.ContainerListOptions{})
		if err != nil {
			continue
		}

		links := []NetworkRoute{}
		seen := make(map[string]bool)
		blacklist := []string{"dashboard", "api", "traefik", "m3tal"}

		for _, c := range res.Items {
			labels := ""
			for k, v := range c.Labels {
				labels += k + "=" + v + ","
			}

			if m := hostRe.FindStringSubmatch(labels); len(m) > 1 {
				host := m[1]
				if seen[host] {
					continue
				}

				serviceKey := strings.ToLower(strings.Split(host, ".")[0])
				readableName := strings.Title(strings.ReplaceAll(serviceKey, "-", " "))
				
				skip := false
				for _, b := range blacklist {
					if strings.Contains(strings.ToLower(readableName), b) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}

				targetURL := "https://" + host
				status := "enabled"
				if resp, err := hc.Head(targetURL); err != nil || resp.StatusCode >= 500 {
					status = "disabled"
				}

				links = append(links, NetworkRoute{
					Name:      readableName,
					URL:       targetURL,
					Status:    status,
					Image:     c.Image,
					Icon:      fmt.Sprintf("https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/%s.png", serviceKey),
					Container: c.Names[0][1:],
				})
				seen[host] = true
			}
		}

		s.mu.Lock()
		s.NetworkL = links
		s.mu.Unlock()
	}
}

func saveAgent(s *M3talState, mPath, aPath, dPath, gPath, tPath, sPath, nPath string) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		s.mu.RLock()
		mBytes, _ := json.MarshalIndent(s, "", "  ")
		
		aData := struct {Issues []Anomaly `json:"issues"`}{Issues: s.Anomalies}
		aBytes, _ := json.MarshalIndent(aData, "", "  ")

		dData := struct {Actions []Decision `json:"actions"`}{Actions: s.Decisions}
		dBytes, _ := json.MarshalIndent(dData, "", "  ")

		gBytes, _ := json.MarshalIndent(s.GPU, "", "  ")
		tBytes, _ := json.MarshalIndent(s.Temp, "", "  ")
		sBytes, _ := json.MarshalIndent(s.Storage, "", "  ")
		nBytes, _ := json.MarshalIndent(s.NetworkL, "", "  ")
		s.mu.RUnlock()

		writeAtomically(mPath, mBytes)
		writeAtomically(aPath, aBytes)
		writeAtomically(dPath, dBytes)
		writeAtomically(gPath, gBytes)
		writeAtomically(tPath, tBytes)
		writeAtomically(sPath, sBytes)
		writeAtomically(nPath, nBytes)
	}
}

func notifyAgent(s *M3talState) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	alertChat := os.Getenv("TG_ALERT_CHAT_ID")
	actionChat := os.Getenv("TG_ACTION_CHAT_ID")

	if token == "" || alertChat == "" {
		log.Println("⚠️ Telegram Notify: Missing credentials, agent disabled.")
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	alertedAnomalies := make(map[string]time.Time)
	lastDecisionIndex := 0

	for range ticker.C {
		s.mu.RLock()
		muted := time.Now().Unix() < s.MutedUntil
		anomalies := s.Anomalies
		decisions := s.Decisions
		s.mu.RUnlock()

		if muted {
			continue
		}

		// 1. Process Anomalies (Cooldown 1 hour)
		for _, a := range anomalies {
			key := a.Target + ":" + a.Reason
			if last, exists := alertedAnomalies[key]; !exists || time.Since(last) > time.Hour {
				emoji := "🚨"
				if a.Type == "transient" {
					emoji = "🟡"
				}
				msg := fmt.Sprintf("%s <b>M3TAL Anomaly</b>\n<b>Target:</b> <code>%s</code>\n<b>Reason:</b> %s", emoji, a.Target, a.Reason)
				sendTelegram(token, alertChat, msg)
				alertedAnomalies[key] = time.Now()
			}
		}

		// 2. Process New Decisions (Healer Actions)
		if len(decisions) > lastDecisionIndex {
			for i := lastDecisionIndex; i < len(decisions); i++ {
				d := decisions[i]
				msg := fmt.Sprintf("🛡️ <b>M3TAL Healer</b>\n<b>Action:</b> <code>%s</code>\n<b>Target:</b> %s\n<b>Reason:</b> %s", d.Type, d.Target, d.Reason)
				sendTelegram(token, actionChat, msg)
			}
			lastDecisionIndex = len(decisions)
		}
	}
}

func listenerAgent(ctx context.Context, s *M3talState) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	allowedUsers := os.Getenv("ALLOWED_USERS")
	
	if token == "" {
		return
	}

	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	offset := 0

	for {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=30", token, offset)
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		var result struct {
			OK     bool             `json:"ok"`
			Result []TelegramUpdate `json:"result"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if !result.OK {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, u := range result.Result {
			offset = u.UpdateID + 1
			if u.Message == nil || u.Message.Text == "" {
				continue
			}

			// Security Check
			uidStr := strconv.FormatInt(u.Message.From.ID, 10)
			if !strings.Contains(allowedUsers, uidStr) {
				sendTelegram(token, strconv.FormatInt(u.Message.Chat.ID, 10), "⛔ <b>Unauthorized</b>")
				continue
			}

			parts := strings.Fields(u.Message.Text)
			cmd := strings.ToLower(parts[0])
			chatStr := strconv.FormatInt(u.Message.Chat.ID, 10)

			switch cmd {
			case "/status":
				s.mu.RLock()
				sys := s.System
				net := s.Network
				s.mu.RUnlock()
				msg := fmt.Sprintf("🏥 <b>M3TAL Status</b>\nCPU: <b>%.1f%%</b>\nRAM: <b>%.1f%%</b>\nNet: <b>%s</b>", sys.CPU, sys.Mem, net.Down)
				sendTelegram(token, chatStr, msg)

			case "/restart":
				if len(parts) < 2 {
					sendTelegram(token, chatStr, "❓ Usage: <code>/restart &lt;name&gt;</code>")
					continue
				}
				target := parts[1]
				sendTelegram(token, chatStr, fmt.Sprintf("⏳ Restarting <code>%s</code>...", target))
				_, err := cli.ContainerRestart(ctx, target, client.ContainerRestartOptions{})
				if err != nil {
					sendTelegram(token, chatStr, fmt.Sprintf("❌ Error: %v", err))
				} else {
					sendTelegram(token, chatStr, fmt.Sprintf("✅ <code>%s</code> restarted.", target))
				}

			case "/mute":
				hours := 1
				if len(parts) > 1 {
					hours, _ = strconv.Atoi(parts[1])
				}
				s.mu.Lock()
				s.MutedUntil = time.Now().Add(time.Duration(hours) * time.Hour).Unix()
				s.mu.Unlock()
				sendTelegram(token, chatStr, fmt.Sprintf("🔇 Alerts muted for <b>%d hours</b>.", hours))

			case "/unmute":
				s.mu.Lock()
				s.MutedUntil = 0
				s.mu.Unlock()
				sendTelegram(token, chatStr, "🔔 Alerts <b>resumed</b>.")

			case "/help":
				msg := "🤖 <b>M3TAL Commands</b>\n/status - System overview\n/restart &lt;name&gt; - Restart service\n/mute &lt;h&gt; - Mute alerts\n/unmute - Resume alerts"
				sendTelegram(token, chatStr, msg)
			}
		}
	}
}

func writeAtomically(path string, data []byte) {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err == nil {
		_ = os.Rename(tmpPath, path)
	}
}

func sendTelegram(token, chatID, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	body, _ := json.Marshal(map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ Failed to send Telegram: %v", err)
		return
	}
	defer resp.Body.Close()
}
