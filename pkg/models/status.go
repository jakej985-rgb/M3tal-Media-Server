package models

// Status represents the health status of system components.
type Status struct {
	Status     string            `json:"status"`
	Components map[string]string `json:"components"`
	Details    map[string]string `json:"details,omitempty"`
}
