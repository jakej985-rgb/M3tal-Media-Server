package containers

import (
	"context"
	"github.com/moby/moby/client"
	"io"
)

// ContainerInfo represents basic container information
type ContainerInfo struct {
	ID     string   `json:"id"`
	Names  []string `json:"names"`
	Image  string   `json:"image"`
	Status string   `json:"status"`
	State  string   `json:"state"`
	CPU    float64  `json:"cpu"`
	Memory float64  `json:"memory"`
}

// Provider defines the interface for container management
type Provider interface {
	ListContainers() ([]ContainerInfo, error)
	StartContainer(name string) error
	StopContainer(name string) error
	RestartContainer(name string) error
	Logs(name string, tail string) (string, error)
}

// Manager handles Docker container operations (implements Provider)
type Manager struct {
	cli *client.Client
}

// NewManager creates a new Docker manager
func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Manager{cli: cli}, nil
}

// ListContainers returns all containers
func (m *Manager) ListContainers() ([]ContainerInfo, error) {
	ctx := context.Background()
	result, err := m.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var info []ContainerInfo
	for _, c := range result.Items {
		info = append(info, ContainerInfo{
			ID:     c.ID,
			Names:  c.Names,
			Image:  c.Image,
			Status: c.Status,
			State:  string(c.State),
			CPU:    0.0, // TODO: Fetch from stats
			Memory: 0.0, // TODO: Fetch from stats
		})
	}
	return info, nil
}

// StartContainer starts a container by name or ID
func (m *Manager) StartContainer(name string) error {
	ctx := context.Background()
	_, err := m.cli.ContainerStart(ctx, name, client.ContainerStartOptions{})
	return err
}

// StopContainer stops a container by name or ID
func (m *Manager) StopContainer(name string) error {
	ctx := context.Background()
	_, err := m.cli.ContainerStop(ctx, name, client.ContainerStopOptions{})
	return err
}

// RestartContainer restarts a container
func (m *Manager) RestartContainer(name string) error {
	ctx := context.Background()
	_, err := m.cli.ContainerRestart(ctx, name, client.ContainerRestartOptions{})
	return err
}

// Logs returns the logs for a container
func (m *Manager) Logs(name string, tail string) (string, error) {
	ctx := context.Background()
	options := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}
	result, err := m.cli.ContainerLogs(ctx, name, options)
	if err != nil {
		return "", err
	}
	defer result.Close()

	content, err := io.ReadAll(result)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Client returns the underlying Docker client
func (m *Manager) Client() *client.Client {
	return m.cli
}
