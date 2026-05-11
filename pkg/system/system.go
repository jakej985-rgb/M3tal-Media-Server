package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"time"
)

// SystemStats represents system metrics
type SystemStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Uptime      uint64  `json:"uptime"`
	Hostname    string  `json:"hostname"`
}

// GetStats collects system metrics
func GetStats() (*SystemStats, error) {
	stats := &SystemStats{}

	// CPU
	percentages, err := cpu.Percent(time.Second, false)
	if err == nil && len(percentages) > 0 {
		stats.CPUUsage = percentages[0]
	}

	// Memory
	vm, err := mem.VirtualMemory()
	if err == nil {
		stats.MemoryUsage = vm.UsedPercent
	}

	// Disk (root)
	usage, err := disk.Usage("/")
	if err == nil {
		stats.DiskUsage = usage.UsedPercent
	}

	// Host info
	info, err := host.Info()
	if err == nil {
		stats.Uptime = info.Uptime
		stats.Hostname = info.Hostname
	}

	return stats, nil
}
