package docker

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// PortInfo represents container port mapping info
type PortInfo struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort int    `json:"private_port"`
	PublicPort  int    `json:"public_port,omitempty"`
	Type        string `json:"type"`
}

// ContainerInfo represents basic container information
type ContainerInfo struct {
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

// ContainerEvent represents a container lifecycle event
type ContainerEvent struct {
	Action        string // "start", "stop", "die", etc.
	ContainerName string
}

// Provider defines the interface for container management
type Provider interface {
	ListContainers() ([]ContainerInfo, error)
	StartContainer(name string) error
	StopContainer(name string) error
	RestartContainer(name string) error
	Logs(name string, tail string) (string, error)
	SubscribeEvents(ctx context.Context) (<-chan ContainerEvent, error)
	StreamLogs(ctx context.Context, name string, tail string) (io.ReadCloser, error)
}

var globalProvider Provider

// SetProvider overrides the global container provider (Docker, Podman, Mock)
func SetProvider(p Provider) {
	globalProvider = p
}

// GetProvider returns the active container provider
func GetProvider() (Provider, error) {
	if globalProvider == nil {
		p, err := NewDockerManager()
		if err != nil {
			return nil, err
		}
		globalProvider = p
	}
	return globalProvider, nil
}

// DockerManager handles Docker container operations (implements Provider)
type DockerManager struct {
	cli *client.Client
}

// NewDockerManager creates a new Docker manager
func NewDockerManager() (*DockerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerManager{cli: cli}, nil
}

func (m *DockerManager) ListContainers() ([]ContainerInfo, error) {
	ctx := context.Background()
	containers, err := m.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var info []ContainerInfo
	for _, c := range containers {
		var ports []PortInfo
		for _, p := range c.Ports {
			ports = append(ports, PortInfo{
				IP:          p.IP,
				PrivatePort: int(p.PrivatePort),
				PublicPort:  int(p.PublicPort),
				Type:        p.Type,
			})
		}
		var networks []string
		if c.NetworkSettings != nil {
			for netName := range c.NetworkSettings.Networks {
				networks = append(networks, netName)
			}
		}
		info = append(info, ContainerInfo{
			ID:       c.ID,
			Names:    c.Names,
			Image:    c.Image,
			Status:   c.Status,
			State:    c.State,
			Labels:   c.Labels,
			Ports:    ports,
			Networks: networks,
		})
	}
	return info, nil
}

func (m *DockerManager) StartContainer(name string) error {
	ctx := context.Background()
	return m.cli.ContainerStart(ctx, name, types.ContainerStartOptions{})
}

func (m *DockerManager) StopContainer(name string) error {
	ctx := context.Background()
	return m.cli.ContainerStop(ctx, name, container.StopOptions{})
}

func (m *DockerManager) RestartContainer(name string) error {
	ctx := context.Background()
	return m.cli.ContainerRestart(ctx, name, container.StopOptions{})
}

func (m *DockerManager) Logs(name string, tail string) (string, error) {
	ctx := context.Background()
	options := types.ContainerLogsOptions{
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

func (m *DockerManager) SubscribeEvents(ctx context.Context) (<-chan ContainerEvent, error) {
	msgCh, errCh := m.cli.Events(ctx, types.EventsOptions{})
	ch := make(chan ContainerEvent, 100)

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errCh:
				if err != nil {
					return
				}
			case msg := <-msgCh:
				if string(msg.Type) == "container" {
					name := msg.Actor.Attributes["name"]
					if name == "" {
						name = msg.Actor.ID
					}
					ch <- ContainerEvent{
						Action:        string(msg.Action),
						ContainerName: name,
					}
				}
			}
		}
	}()

	return ch, nil
}

func (m *DockerManager) StreamLogs(ctx context.Context, name string, tail string) (io.ReadCloser, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       tail,
	}
	return m.cli.ContainerLogs(ctx, name, options)
}

// MockProvider for testing (implements Provider)
type MockProvider struct {
	Containers []ContainerInfo
}

func (m *MockProvider) ListContainers() ([]ContainerInfo, error) { return m.Containers, nil }
func (m *MockProvider) StartContainer(name string) error         { return nil }
func (m *MockProvider) StopContainer(name string) error          { return nil }
func (m *MockProvider) RestartContainer(name string) error       { return nil }
func (m *MockProvider) Logs(name string, tail string) (string, error) {
	return "mock logs for " + name, nil
}
func (m *MockProvider) SubscribeEvents(ctx context.Context) (<-chan ContainerEvent, error) {
	ch := make(chan ContainerEvent)
	return ch, nil
}
func (m *MockProvider) StreamLogs(ctx context.Context, name string, tail string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("mock streaming logs for " + name + "\n")), nil
}
