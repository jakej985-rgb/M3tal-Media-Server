package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/moby/moby/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// M3talState holds the unified state of the control plane in memory
type M3talState struct {
	mu         sync.RWMutex
	System     SystemMetrics     `json:"system"`
	Containers []ContainerMetric `json:"containers"`
	Network    NetworkMetrics    `json:"network"`
	Anomalies  []Anomaly         `json:"issues"` // For anomalies.json
	Timestamp  int64             `json:"timestamp"`
	CPU        float64           `json:"cpu"`
}

type Anomaly struct {
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

type ContainerMetric struct {
	Name string  `json:"name"`
	CPU  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
}

type NetworkMetrics struct {
	Down float64 `json:"down"`
	Up   float64 `json:"up"`
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

	fmt.Println("🚀 M3TAL Go Backend (Linux-Ready) starting...")
	fmt.Printf("📂 Target State: %s\n", metricsPath)

	ctx := context.Background()

	// 1. Metrics Agent (System)
	go metricsAgent(state)

	// 2. Docker Agent (Containers)
	go dockerAgent(ctx, state)

	// 3. Anomaly Agent (Detection)
	go anomalyAgent(state)

	// 4. Persistence Agent (JSON Sync)
	anomalyPath := filepath.Join(stateDir, "anomalies.json")
	go saveAgent(state, metricsPath, anomalyPath)

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
			Down: down,
			Up:   up,
			Load: load,
		}
		s.CPU = cpuVal
		s.Timestamp = now.Unix()
		s.mu.Unlock()
	}
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
				Name: name,
				CPU:  0.0,
				Mem:  0.0,
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

func saveAgent(s *M3talState, mPath, aPath string) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		// Save Metrics
		s.mu.RLock()
		mBytes, _ := json.MarshalIndent(s, "", "  ")
		
		// Save Anomalies (Subset)
		aData := struct {
			Issues []Anomaly `json:"issues"`
		}{Issues: s.Anomalies}
		aBytes, _ := json.MarshalIndent(aData, "", "  ")
		s.mu.RUnlock()

		// Atomic write for metrics
		tmpM := mPath + ".tmp"
		if err := os.WriteFile(tmpM, mBytes, 0644); err == nil {
			_ = os.Rename(tmpM, mPath)
		}

		// Atomic write for anomalies
		tmpA := aPath + ".tmp"
		if err := os.WriteFile(tmpA, aBytes, 0644); err == nil {
			_ = os.Rename(tmpA, aPath)
		}
	}
}
