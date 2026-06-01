package models

// SystemInfo represents basic system information.
type SystemInfo struct {
	Hostname string `json:"hostname"`
	Uptime   uint64 `json:"uptime"`
}

// DetailedStats represents extended system metrics for dashboard/tray.
type DetailedStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	CPUTemp     float64 `json:"cpu_temp"`
	MemoryUsage float64 `json:"memory_usage"`
	MemoryUsed  float64 `json:"memory_used"`  // GB
	MemoryTotal float64 `json:"memory_total"` // GB
	DiskUsage   float64 `json:"disk_usage"`
	DiskUsed    float64 `json:"disk_used"`  // GB
	DiskTotal   float64 `json:"disk_total"` // GB
	GPUUsage    float64 `json:"gpu_usage"`
	GPUTemp     float64 `json:"gpu_temp"`
	GPUMemUsed  float64 `json:"gpu_mem_used"`  // MB
	GPUMemTotal float64 `json:"gpu_mem_total"` // MB
	GPUModel    string  `json:"gpu_model"`
	Uptime      uint64  `json:"uptime"`
	Hostname    string  `json:"hostname"`
}
