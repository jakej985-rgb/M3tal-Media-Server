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

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// M3talState holds the unified state of the control plane in memory
type M3talState struct {
	mu         sync.RWMutex
	System     SystemMetrics     `json:"system"`
	Containers []ContainerMetric `json:"containers"`
	Network    NetworkMetrics    `json:"network"`
	Timestamp  int64             `json:"timestamp"`
	CPU        float64           `json:"cpu"`
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

	// 3. Persistence Agent (JSON Sync)
	go saveAgent(state, metricsPath)

	// Keep alive
	select {}
}

func metricsAgent(s *M3talState) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		now := time.Now().Unix()
		cpuPerc, _ := cpu.Percent(0, false)
		vm, _ := mem.VirtualMemory()
		
		cpuVal := 0.0
		if len(cpuPerc) > 0 {
			cpuVal = cpuPerc[0]
		}

		s.mu.Lock()
		s.System = SystemMetrics{
			CPU:       cpuVal,
			Mem:       vm.UsedPercent,
			MemGB:     float64(vm.Used) / (1024 * 1024 * 1024),
			Timestamp: now,
		}
		s.CPU = cpuVal
		s.Timestamp = now
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
		containers, err := cli.ContainerList(ctx, container.ListOptions{})
		if err != nil {
			log.Printf("⚠️ Failed to list containers: %v", err)
			continue
		}

		newStats := []ContainerMetric{}
		for _, c := range containers {
			name := "unknown"
			if len(c.Names) > 0 {
				name = c.Names[0][1:] // Remove leading slash
			}
			newStats = append(newStats, ContainerMetric{
				Name: name,
				CPU:  0.0, // Stats require a separate call, keeping it simple for now
				Mem:  0.0,
			})
		}

		s.mu.Lock()
		s.Containers = newStats
		s.mu.Unlock()
	}
}

func saveAgent(s *M3talState, path string) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		s.mu.RLock()
		data, err := json.MarshalIndent(s, "", "  ")
		s.mu.RUnlock()

		if err != nil {
			continue
		}

		tmpPath := path + ".tmp"
		if err := os.WriteFile(tmpPath, data, 0644); err == nil {
			_ = os.Rename(tmpPath, path)
		}
	}
}
