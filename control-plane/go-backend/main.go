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
	Anomalies  []Anomaly         `json:"issues"`    // For anomalies.json
	Decisions  []Decision        `json:"actions"`   // For decisions.json
	Timestamp  int64             `json:"timestamp"`
	CPU        float64           `json:"cpu"`
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

type ContainerMetric struct {
	Name   string  `json:"name"`
	CPU    float64 `json:"cpu"`
	Mem    float64 `json:"mem"`
	Status string  `json:"status"` // e.g. "running"
	State  string  `json:"state"`  // e.g. "exited"
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

	// 4. Healer Agent (Auto-Recovery)
	go healerAgent(ctx, state)

	// 5. Persistence Agent (JSON Sync)
	anomalyPath := filepath.Join(stateDir, "anomalies.json")
	decisionPath := filepath.Join(stateDir, "decisions.json")
	go saveAgent(state, metricsPath, anomalyPath, decisionPath)

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

func saveAgent(s *M3talState, mPath, aPath, dPath string) {
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

		// Save Decisions (Subset)
		dData := struct {
			Actions []Decision `json:"actions"`
		}{Actions: s.Decisions}
		dBytes, _ := json.MarshalIndent(dData, "", "  ")
		s.mu.RUnlock()

		// Atomic writes
		writeAtomically(mPath, mBytes)
		writeAtomically(aPath, aBytes)
		writeAtomically(dPath, dBytes)
	}
}

func writeAtomically(path string, data []byte) {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err == nil {
		_ = os.Rename(tmpPath, path)
	}
}
