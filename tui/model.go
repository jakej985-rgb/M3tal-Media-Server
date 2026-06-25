package tui

import (
	"time"

	"github.com/jakej985-rgb/m3tal-core/pkg/client"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
)

// ConfigFile represents a dynamic configuration file discovered on disk.
type ConfigFile struct {
	Name string
	Path string
}

// Tab represents each TUI view pane.
type Tab int

const (
	TabDashboard Tab = iota
	TabContainers
	TabLogs
	TabEditor
	TabSystem
	TabTerminal
)

// Message types for Bubble Tea
type tickMsg time.Time

type detailedMetricsMsg struct {
	stats *models.DetailedStats
	err   error
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

type actionResultMsg struct {
	message string
	err     error
}

type configMsg struct {
	config      map[string]string
	cloudflared string
	envRaw      string
	err         error
}

type editFinishedMsg struct {
	content   string
	isEnvEdit bool
	err       error
}

type doctorReportMsg struct {
	report models.DoctorReport
	err    error
}

type dockerImagesMsg struct {
	images []models.DockerImage
	err    error
}

type dockerVolumesMsg struct {
	volumes []models.DockerVolume
	err     error
}

type dockerNetworksMsg struct {
	networks []models.DockerNetwork
	err      error
}

type systemServicesMsg struct {
	services []models.ServiceStatus
	err      error
}

type systemStorageMsg struct {
	partitions []models.DiskPartition
	samba      []models.SambaShare
	err        error
}

type cronJobsMsg struct {
	jobs []models.CronJob
	err  error
}

type systemUpdatesMsg struct {
	updates *models.SystemUpdates
	err     error
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

	// API Data (Detailed)
	detailedStats *models.DetailedStats
	status        *models.Status

	// Tab 1: Dashboard State
	doctorReport           models.DoctorReport
	dashboardFocusIndex    int // 0 = metrics/health, 1 = doctor alerts, 2 = quick actions
	selectedQuickActionIdx int

	// Tab 2: Containers & Docker Controls State
	containersTabFocus    int // 0 = stacks list, 1 = services list, 2 = images, 3 = volumes, 4 = networks
	stacks                []models.Stack
	selectedStackIdx      int
	containers            []models.Container
	selectedContainerIdx  int
	dockerImages          []models.DockerImage
	selectedImageIdx      int
	dockerVolumes         []models.DockerVolume
	selectedVolumeIdx     int
	dockerNetworks        []models.DockerNetwork
	selectedNetworkIdx    int

	// Tab 3: Logs State
	logContainers           []models.Container
	selectedLogContainerIdx int
	logs                    string
	logScrollOffset         int
	logScrollHeight         int
	lastLogContainer        string
	logLevelFilter          string // "ALL", "INFO", "WARN", "ERROR"
	logSearchQuery          string
	logFollow               bool
	showLogSearchPrompt     bool // state of search prompt input

	// Tab 4: Editor State
	configFiles             []ConfigFile
	selectedConfigIdx       int
	selectedConfigContent   string
	configScrollOffset      int
	focusOnConfig           bool // true = left panel (configs list), false = right panel (content viewer)
	diffContent             string
	showDiff                bool
	cloudflaredContent      string
	envRawContent           string

	// Tab 5: System Admin State
	systemTabFocus          int // 0 = services, 1 = disks/mounts, 2 = cron, 3 = updates
	systemServices          []models.ServiceStatus
	selectedServiceIdx      int
	diskPartitions          []models.DiskPartition
	selectedPartitionIdx    int
	sambaShares             []models.SambaShare
	selectedSambaIdx        int
	cronJobs                []models.CronJob
	selectedCronIdx         int
	systemUpdates           *models.SystemUpdates
	selectedUpdateIdx       int

	// Tab 6: Terminal & Integrations State
	terminalTabFocus        int // 0 = Terminal Launcher, 1 = Saved SSH, 2 = Env variables, 3 = AI Queue, 4 = Plugins
	sshProfiles             []string
	selectedSSHIdx          int
	envVars                 []string
	selectedEnvVarIdx       int
	savedAliases            []string
	selectedAliasIdx        int

	// AI Queue & Models State (integrated under Terminal/Shell view)
	queue            []models.JobRecord
	selectedJobIdx   int
	aiModels         []string
	selectedModelIdx int

	// Plugins State (integrated under Terminal/Shell view)
	plugins           *models.PluginsResponse
	catalog           []models.CatalogItemStatus
	showCatalog       bool // true = Catalog view, false = Installed view
	selectedPluginIdx int
}

// NewModel initializes the TUI state model.
func NewModel(c *client.Client) Model {
	return Model{
		client:             c,
		activeTab:          TabDashboard,
		containersTabFocus: 0,
		logFollow:          true,
		logLevelFilter:     "ALL",
		focusOnConfig:      true,
		systemTabFocus:     0,
		terminalTabFocus:   0,
		sshProfiles:        []string{"localhost", "staging-server", "production-edge"},
		envVars:            []string{"BASE_STORAGE_PATH=/docker", "API_TOKEN=m3tal-secret", "DASHBOARD_PORT=8082"},
		savedAliases:       []string{"dc=docker compose", "dps=docker ps --format 'table {{.Names}}\t{{.Status}}'", "dprune=docker system prune -a --volumes"},
		showCatalog:        false,
	}
}

// SetNotification triggers a status bar toast notification.
func (m *Model) SetNotification(msg string, duration time.Duration) {
	m.notification = msg
	m.notificationTimeout = time.Now().Add(duration)
}
