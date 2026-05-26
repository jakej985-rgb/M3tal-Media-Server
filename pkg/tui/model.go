package tui

import (
	"time"

	"github.com/jakej985-rgb/m3tal-core/pkg/client"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
)

// Tab represents each TUI view pane.
type Tab int

const (
	TabStacks Tab = iota
	TabLogs
	TabAI
	TabPlugins
)

// Message types for Bubble Tea
type errMsg error
type tickMsg time.Time

type metricsMsg struct {
	metrics *models.MetricsResponse
	err     error
}

type statusMsg struct {
	status *models.Status
	err    error
}

type stacksMsg struct {
	stacks []models.Stack
	err    error
}

type containersMsg struct {
	containers []models.Container
	err        error
}

type logsMsg struct {
	container string
	logs      string
	err       error
}

type aiMsg struct {
	queue  []models.JobRecord
	models []string
	err    error
}

type pluginsMsg struct {
	plugins *models.PluginsResponse
	catalog []models.CatalogItemStatus
	err     error
}

type actionResultMsg struct {
	message string
	err     error
}

// Model defines the state model for M3TAL Go TUI dashboard.
type Model struct {
	client      *client.Client
	activeTab   Tab
	width       int
	height      int
	lastUpdated time.Time
	err         error

	// Top notification bar
	notification        string
	notificationTimeout time.Time

	// API Data
	metrics *models.MetricsResponse
	status  *models.Status

	// Stacks Tab State
	stacks            []models.Stack
	selectedStackIdx  int
	containers        []models.Container
	selectedContainerIdx int
	focusOnStacks     bool // true = stacks list, false = services list

	// Logs Tab State
	logContainers            []models.Container
	selectedLogContainerIdx  int
	logs                     string
	logScrollOffset          int
	logScrollHeight          int
	lastLogContainer         string
	lastLogContent           string

	// AI Queue & Models State
	queue             []models.JobRecord
	selectedJobIdx    int
	aiModels          []string
	selectedModelIdx  int
	focusOnQueue      bool // true = job queue list, false = Ollama models list

	// Plugins Tab State
	plugins           *models.PluginsResponse
	catalog           []models.CatalogItemStatus
	showCatalog       bool // true = Catalog view, false = Installed view
	selectedPluginIdx int
}

// NewModel initializes the TUI state model.
func NewModel(c *client.Client) Model {
	return Model{
		client:        c,
		activeTab:     TabStacks,
		focusOnStacks: true,
		focusOnQueue:  true,
		showCatalog:   false,
	}
}

// SetNotification triggers a status bar toast notification.
func (m *Model) SetNotification(msg string, duration time.Duration) {
	m.notification = msg
	m.notificationTimeout = time.Now().Add(duration)
}
