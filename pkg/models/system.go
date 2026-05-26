package models

// SystemInfo represents basic system information.
type SystemInfo struct {
	Hostname string `json:"hostname"`
	Uptime   uint64 `json:"uptime"`
}
