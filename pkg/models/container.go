package models

// PortInfo represents container port mapping info.
type PortInfo struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort int    `json:"private_port"`
	PublicPort  int    `json:"public_port,omitempty"`
	Type        string `json:"type"`
}

// Container represents basic container information.
type Container struct {
	ID       string            `json:"id"`
	Names    []string          `json:"names"`
	Image    string            `json:"image"`
	Status   string            `json:"status"`
	State    string            `json:"state"`
	CPU      float64           `json:"cpu"`
	Memory   float64           `json:"memory"`
	Labels   map[string]string `json:"labels,omitempty"`
	Ports    []PortInfo        `json:"ports,omitempty"`
	Networks []string          `json:"networks,omitempty"`
}
