package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
	"github.com/jakej985-rgb/m3tal-core/pkg/system"
)

// Init initializes Bubble Tea commands.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchAllDataCmd(),
		tick(),
	)
}

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles incoming messages and updates TUI state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		cmds = append(cmds, m.fetchAllDataCmd(), tick())

	case detailedMetricsMsg:
		if msg.err == nil {
			m.detailedStats = msg.stats
			m.err = nil
		} else {
			m.err = msg.err
		}

	case statusMsg:
		if msg.err == nil {
			m.status = msg.status
			m.err = nil
		} else {
			m.err = msg.err
		}

	case doctorReportMsg:
		if msg.err == nil {
			m.doctorReport = msg.report
			m.err = nil
		} else {
			m.err = msg.err
		}

	case stacksMsg:
		if msg.err == nil {
			m.stacks = msg.stacks
			m.err = nil
			if m.selectedStackIdx >= len(m.stacks) {
				m.selectedStackIdx = 0
			}
			if len(m.stacks) > 0 {
				cmds = append(cmds, m.fetchContainersCmd(m.stacks[m.selectedStackIdx].Name))
			}
		} else {
			m.err = msg.err
		}

	case containersMsg:
		if msg.err == nil {
			// Update containers in Containers view
			if m.selectedStackIdx < len(m.stacks) {
				stackName := strings.ToLower(m.stacks[m.selectedStackIdx].Name)
				trimmedStackName := strings.TrimLeft(stackName, "0123456789-")
				var filtered []models.Container
				for _, c := range msg.containers {
					proj := strings.ToLower(c.Labels["com.docker.compose.project"])
					mStack := strings.ToLower(c.Labels["m3tal.stack"])
					configFiles := strings.ToLower(c.Labels["com.docker.compose.project.config_files"])

					isMatch := proj == stackName ||
						proj == trimmedStackName ||
						mStack == trimmedStackName ||
						mStack == stackName ||
						strings.Contains(configFiles, stackName+"-compose.yml") ||
						strings.Contains(configFiles, trimmedStackName+"-compose.yml")

					if !isMatch && len(c.Names) > 0 {
						cname := strings.ToLower(strings.TrimPrefix(c.Names[0], "/"))
						isMatch = strings.HasPrefix(cname, stackName+"-") ||
							strings.HasPrefix(cname, trimmedStackName+"-") ||
							strings.HasPrefix(cname, "m3tal-"+trimmedStackName) ||
							strings.HasPrefix(cname, "m3tal-"+stackName)
					}

					if isMatch {
						filtered = append(filtered, c)
					}
				}
				m.containers = filtered
				if m.selectedContainerIdx >= len(m.containers) {
					m.selectedContainerIdx = 0
				}
			}

			// Update container logs selector list (includes all active containers)
			m.logContainers = msg.containers
			if m.selectedLogContainerIdx >= len(m.logContainers) {
				m.selectedLogContainerIdx = 0
			}
			if len(m.logContainers) > 0 && m.logs == "" {
				selectedC := m.logContainers[m.selectedLogContainerIdx]
				cname := "unknown"
				if len(selectedC.Names) > 0 {
					cname = strings.TrimPrefix(selectedC.Names[0], "/")
				}
				cmds = append(cmds, m.fetchLogsCmd(cname))
			}
			m.err = nil
		} else {
			m.err = msg.err
		}

	case dockerImagesMsg:
		if msg.err == nil {
			m.dockerImages = msg.images
			if m.selectedImageIdx >= len(m.dockerImages) {
				m.selectedImageIdx = 0
			}
		}

	case dockerVolumesMsg:
		if msg.err == nil {
			m.dockerVolumes = msg.volumes
			if m.selectedVolumeIdx >= len(m.dockerVolumes) {
				m.selectedVolumeIdx = 0
			}
		}

	case dockerNetworksMsg:
		if msg.err == nil {
			m.dockerNetworks = msg.networks
			if m.selectedNetworkIdx >= len(m.dockerNetworks) {
				m.selectedNetworkIdx = 0
			}
		}

	case systemServicesMsg:
		if msg.err == nil {
			m.systemServices = msg.services
			if m.selectedServiceIdx >= len(m.systemServices) {
				m.selectedServiceIdx = 0
			}
		}

	case systemStorageMsg:
		if msg.err == nil {
			m.diskPartitions = msg.partitions
			m.sambaShares = msg.samba
			if m.selectedPartitionIdx >= len(m.diskPartitions) {
				m.selectedPartitionIdx = 0
			}
			if m.selectedSambaIdx >= len(m.sambaShares) {
				m.selectedSambaIdx = 0
			}
		}

	case cronJobsMsg:
		if msg.err == nil {
			m.cronJobs = msg.jobs
			if m.selectedCronIdx >= len(m.cronJobs) {
				m.selectedCronIdx = 0
			}
		}

	case systemUpdatesMsg:
		if msg.err == nil {
			m.systemUpdates = msg.updates
		}

	case logsMsg:
		if msg.err == nil {
			m.logs = msg.logs
			m.err = nil
		} else {
			m.err = msg.err
		}

	case aiMsg:
		if msg.err == nil {
			m.queue = msg.queue
			m.aiModels = msg.models
			m.err = nil
			if m.selectedJobIdx >= len(m.queue) {
				m.selectedJobIdx = 0
			}
			if m.selectedModelIdx >= len(m.aiModels) {
				m.selectedModelIdx = 0
			}
		} else {
			m.err = msg.err
		}

	case pluginsMsg:
		if msg.err == nil {
			m.plugins = msg.plugins
			m.catalog = msg.catalog
			m.err = nil
			listLen := m.getPluginsListLen()
			if m.selectedPluginIdx >= listLen {
				m.selectedPluginIdx = 0
			}
		} else {
			m.err = msg.err
		}

	case configMsg:
		if msg.err == nil {
			m.cloudflaredContent = msg.cloudflared
			m.envRawContent = msg.envRaw
			m.err = nil
			m.refreshConfigFiles()
			m.loadSelectedConfigContent()
		} else {
			m.err = msg.err
		}

	case editFinishedMsg:
		if msg.err != nil {
			m.SetNotification(fmt.Sprintf("❌ Edit failed: %v", msg.err), 4*time.Second)
		} else {
			m.SetNotification("✅ Config updated!", 3*time.Second)
			m.loadSelectedConfigContent()
		}
		cmds = append(cmds, m.fetchAllDataCmd())

	case actionResultMsg:
		if msg.err != nil {
			m.SetNotification(fmt.Sprintf("❌ Error: %v", msg.err), 4*time.Second)
		} else {
			m.SetNotification(msg.message, 3*time.Second)
		}
		cmds = append(cmds, m.fetchAllDataCmd())

	case tea.KeyMsg:
		keyStr := strings.ToLower(msg.String())

		// If in log search query input mode, capture input characters
		if m.showLogSearchPrompt {
			switch keyStr {
			case "enter":
				m.showLogSearchPrompt = false
				m.logScrollOffset = 0
				m.logs = ""
				if len(m.logContainers) > 0 {
					selectedC := m.logContainers[m.selectedLogContainerIdx]
					cname := "unknown"
					if len(selectedC.Names) > 0 {
						cname = strings.TrimPrefix(selectedC.Names[0], "/")
					}
					cmds = append(cmds, m.fetchLogsCmd(cname))
				}
			case "esc":
				m.showLogSearchPrompt = false
				m.logSearchQuery = ""
			case "backspace":
				if len(m.logSearchQuery) > 0 {
					m.logSearchQuery = m.logSearchQuery[:len(m.logSearchQuery)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.logSearchQuery += msg.String()
				}
			}
			return m, nil
		}

		switch keyStr {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			m.SetNotification("Refreshing data...", 1*time.Second)
			cmds = append(cmds, m.fetchAllDataCmd())

		// Tab Switching Keys 1-6
		case "1":
			m.activeTab = TabDashboard
			cmds = append(cmds, m.fetchAllDataCmd())
		case "2":
			m.activeTab = TabContainers
			cmds = append(cmds, m.fetchAllDataCmd())
		case "3":
			m.activeTab = TabLogs
			m.logs = ""
			m.logScrollOffset = 0
			cmds = append(cmds, m.fetchAllDataCmd())
		case "4":
			m.activeTab = TabEditor
			m.refreshConfigFiles()
			m.loadSelectedConfigContent()
			cmds = append(cmds, m.fetchAllDataCmd())
		case "5":
			m.activeTab = TabSystem
			cmds = append(cmds, m.fetchAllDataCmd())
		case "6":
			m.activeTab = TabTerminal
			cmds = append(cmds, m.fetchAllDataCmd())

		case "tab":
			m.cycleFocus()

		case "up", "k":
			m.handleNavigationUp()
			cmds = append(cmds, m.postNavCmds()...)

		case "down", "j":
			m.handleNavigationDown()
			cmds = append(cmds, m.postNavCmds()...)

		// --- Contextual Actions ---
		case "enter":
			cmds = append(cmds, m.handleEnterAction()...)

		// Start Svc
		case "s":
			if m.activeTab == TabContainers && m.containersTabFocus == 1 && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🐳 Starting container %s...", cname), 5*time.Second)
				cmds = append(cmds, m.controlContainerCmd(cname, "start"))
			} else if m.activeTab == TabSystem && m.systemTabFocus == 0 && len(m.systemServices) > 0 {
				svc := m.systemServices[m.selectedServiceIdx]
				m.SetNotification(fmt.Sprintf("⚙️ Starting service %s...", svc.Name), 5*time.Second)
				cmds = append(cmds, m.controlSystemServiceCmd(svc.Name, "start"))
			}

		// Stop Svc
		case "x":
			if m.activeTab == TabContainers && m.containersTabFocus == 1 && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🛑 Stopping container %s...", cname), 5*time.Second)
				cmds = append(cmds, m.controlContainerCmd(cname, "stop"))
			} else if m.activeTab == TabSystem && m.systemTabFocus == 0 && len(m.systemServices) > 0 {
				svc := m.systemServices[m.selectedServiceIdx]
				m.SetNotification(fmt.Sprintf("🛑 Stopping service %s...", svc.Name), 5*time.Second)
				cmds = append(cmds, m.controlSystemServiceCmd(svc.Name, "stop"))
			}

		// Restart Svc
		case "t":
			if m.activeTab == TabContainers && m.containersTabFocus == 1 && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🔄 Restarting container %s...", cname), 5*time.Second)
				cmds = append(cmds, m.controlContainerCmd(cname, "restart"))
			} else if m.activeTab == TabSystem && m.systemTabFocus == 0 && len(m.systemServices) > 0 {
				svc := m.systemServices[m.selectedServiceIdx]
				m.SetNotification(fmt.Sprintf("🔄 Restarting service %s...", svc.Name), 5*time.Second)
				cmds = append(cmds, m.controlSystemServiceCmd(svc.Name, "restart"))
			}

		// Toggle Enable Systemd
		case "v":
			if m.activeTab == TabSystem && m.systemTabFocus == 0 && len(m.systemServices) > 0 {
				svc := m.systemServices[m.selectedServiceIdx]
				action := "enable"
				if svc.Enabled == "enabled" {
					action = "disable"
				}
				m.SetNotification(fmt.Sprintf("⚙️ Setting service %s to %sd...", svc.Name, action), 5*time.Second)
				cmds = append(cmds, m.controlSystemServiceCmd(svc.Name, action))
			}

		// Deploy Stack
		case "u":
			if m.activeTab == TabContainers && m.containersTabFocus == 0 && len(m.stacks) > 0 {
				stackName := m.stacks[m.selectedStackIdx].Name
				m.SetNotification(fmt.Sprintf("🚀 Deploying stack %s...", stackName), 10*time.Second)
				cmds = append(cmds, m.deployStackCmd(stackName))
			} else if m.activeTab == TabSystem && m.systemTabFocus == 3 {
				m.SetNotification("📥 Applying system updates...", 15*time.Second)
				cmds = append(cmds, m.applySystemUpdatesCmd())
			}

		// Stop Stack
		case "d":
			if m.activeTab == TabContainers && m.containersTabFocus == 0 && len(m.stacks) > 0 {
				stackName := m.stacks[m.selectedStackIdx].Name
				m.SetNotification(fmt.Sprintf("🧹 Stopping stack %s...", stackName), 10*time.Second)
				cmds = append(cmds, m.stopStackCmd(stackName))
			}

		// Exec into Container or Edit Config
		case "e":
			if m.activeTab == TabContainers && m.containersTabFocus == 1 && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🐳 Suspending TUI. Executing bash shell in container %s...", cname), 2*time.Second)
				cmds = append(cmds, m.execContainerCmd(cname))
			} else if m.activeTab == TabEditor {
				cmds = append(cmds, m.editConfigCmd())
			}

		// Pull updates / scan stacks
		case "p":
			if m.activeTab == TabContainers && m.containersTabFocus == 0 && len(m.stacks) > 0 {
				stackName := m.stacks[m.selectedStackIdx].Name
				m.SetNotification(fmt.Sprintf("📥 Pulling stack %s images...", stackName), 10*time.Second)
				cmds = append(cmds, m.pullStackCmd(stackName))
			}

		// Logs tab: Log level cycle
		case "l":
			if m.activeTab == TabLogs {
				switch m.logLevelFilter {
				case "ALL":
					m.logLevelFilter = "INFO"
				case "INFO":
					m.logLevelFilter = "WARN"
				case "WARN":
					m.logLevelFilter = "ERROR"
				case "ERROR":
					m.logLevelFilter = "ALL"
				}
				m.SetNotification(fmt.Sprintf("🔍 Filter log level: %s", m.logLevelFilter), 2*time.Second)
			}

		// Logs tab: Toggle Follow Mode
		case "f":
			if m.activeTab == TabLogs {
				m.logFollow = !m.logFollow
				stateStr := "Disabled"
				if m.logFollow {
					stateStr = "Enabled"
				}
				m.SetNotification(fmt.Sprintf("📄 Log Follow Mode %s", stateStr), 2*time.Second)
			}

		// Logs tab: Search Filter Prompt
		case "/":
			if m.activeTab == TabLogs {
				m.showLogSearchPrompt = true
				m.logSearchQuery = ""
			}

		// Logs tab: Export logs
		case "o":
			if m.activeTab == TabLogs && len(m.logContainers) > 0 {
				selectedC := m.logContainers[m.selectedLogContainerIdx]
				cname := "unknown"
				if len(selectedC.Names) > 0 {
					cname = strings.TrimPrefix(selectedC.Names[0], "/")
				}
				m.SetNotification(fmt.Sprintf("💾 Exporting logs for %s...", cname), 2*time.Second)
				cmds = append(cmds, m.exportLogsCmd(cname, m.logs))
			}

		// Scan Stacks
		case "a":
			if m.activeTab == TabContainers {
				m.SetNotification("🔍 Scanning stacks...", 5*time.Second)
				cmds = append(cmds, m.scanStacksCmd())
			}

		// Plugins catalog toggle
		case "c":
			if m.activeTab == TabTerminal && m.terminalTabFocus == 4 {
				m.showCatalog = !m.showCatalog
				m.selectedPluginIdx = 0
			}

		// Plugins enable/disable
		case "g":
			if m.activeTab == TabTerminal && m.terminalTabFocus == 4 {
				cmds = append(cmds, m.togglePluginCmd())
			}

		// Plugins install
		case "i":
			if m.activeTab == TabTerminal && m.terminalTabFocus == 4 && m.showCatalog {
				cmds = append(cmds, m.installPluginCmd())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) cycleFocus() {
	switch m.activeTab {
	case TabDashboard:
		m.dashboardFocusIndex = (m.dashboardFocusIndex + 1) % 3
	case TabContainers:
		m.containersTabFocus = (m.containersTabFocus + 1) % 5
	case TabLogs:
		m.focusOnConfig = !m.focusOnConfig // left container list vs right viewer
	case TabEditor:
		m.focusOnConfig = !m.focusOnConfig
	case TabSystem:
		m.systemTabFocus = (m.systemTabFocus + 1) % 4
	case TabTerminal:
		m.terminalTabFocus = (m.terminalTabFocus + 1) % 5
	}
}

func (m *Model) handleNavigationUp() {
	switch m.activeTab {
	case TabDashboard:
		if m.dashboardFocusIndex == 2 {
			if m.selectedQuickActionIdx > 0 {
				m.selectedQuickActionIdx--
			}
		}
	case TabContainers:
		switch m.containersTabFocus {
		case 0:
			if m.selectedStackIdx > 0 {
				m.selectedStackIdx--
			}
		case 1:
			if m.selectedContainerIdx > 0 {
				m.selectedContainerIdx--
			}
		case 2:
			if m.selectedImageIdx > 0 {
				m.selectedImageIdx--
			}
		case 3:
			if m.selectedVolumeIdx > 0 {
				m.selectedVolumeIdx--
			}
		case 4:
			if m.selectedNetworkIdx > 0 {
				m.selectedNetworkIdx--
			}
		}
	case TabLogs:
		if m.focusOnConfig {
			if m.selectedLogContainerIdx > 0 {
				m.selectedLogContainerIdx--
			}
		} else {
			m.logScrollOffset++
		}
	case TabEditor:
		if m.focusOnConfig {
			if m.selectedConfigIdx > 0 {
				m.selectedConfigIdx--
				m.configScrollOffset = 0
				m.loadSelectedConfigContent()
			}
		} else {
			if m.configScrollOffset > 0 {
				m.configScrollOffset--
			}
		}
	case TabSystem:
		switch m.systemTabFocus {
		case 0:
			if m.selectedServiceIdx > 0 {
				m.selectedServiceIdx--
			}
		case 1:
			// Disks Mounts cycle
			if m.selectedPartitionIdx > 0 {
				m.selectedPartitionIdx--
			}
		case 2:
			if m.selectedCronIdx > 0 {
				m.selectedCronIdx--
			}
		case 3:
			if m.selectedUpdateIdx > 0 {
				m.selectedUpdateIdx--
			}
		}
	case TabTerminal:
		switch m.terminalTabFocus {
		case 1:
			if m.selectedSSHIdx > 0 {
				m.selectedSSHIdx--
			}
		case 2:
			if m.selectedEnvVarIdx > 0 {
				m.selectedEnvVarIdx--
			}
		case 3:
			if m.selectedJobIdx > 0 {
				m.selectedJobIdx--
			}
		case 4:
			if m.selectedPluginIdx > 0 {
				m.selectedPluginIdx--
			}
		}
	}
}

func (m *Model) handleNavigationDown() {
	switch m.activeTab {
	case TabDashboard:
		if m.dashboardFocusIndex == 2 {
			if m.selectedQuickActionIdx < 2 { // 3 quick actions
				m.selectedQuickActionIdx++
			}
		}
	case TabContainers:
		switch m.containersTabFocus {
		case 0:
			if m.selectedStackIdx < len(m.stacks)-1 {
				m.selectedStackIdx++
			}
		case 1:
			if m.selectedContainerIdx < len(m.containers)-1 {
				m.selectedContainerIdx++
			}
		case 2:
			if m.selectedImageIdx < len(m.dockerImages)-1 {
				m.selectedImageIdx++
			}
		case 3:
			if m.selectedVolumeIdx < len(m.dockerVolumes)-1 {
				m.selectedVolumeIdx++
			}
		case 4:
			if m.selectedNetworkIdx < len(m.dockerNetworks)-1 {
				m.selectedNetworkIdx++
			}
		}
	case TabLogs:
		if m.focusOnConfig {
			if m.selectedLogContainerIdx < len(m.logContainers)-1 {
				m.selectedLogContainerIdx++
			}
		} else {
			if m.logScrollOffset > 0 {
				m.logScrollOffset--
			}
		}
	case TabEditor:
		if m.focusOnConfig {
			if m.selectedConfigIdx < len(m.configFiles)-1 {
				m.selectedConfigIdx++
				m.configScrollOffset = 0
				m.loadSelectedConfigContent()
			}
		} else {
			lines := strings.Split(m.selectedConfigContent, "\n")
			if m.configScrollOffset < len(lines)-1 {
				m.configScrollOffset++
			}
		}
	case TabSystem:
		switch m.systemTabFocus {
		case 0:
			if m.selectedServiceIdx < len(m.systemServices)-1 {
				m.selectedServiceIdx++
			}
		case 1:
			if m.selectedPartitionIdx < len(m.diskPartitions)-1 {
				m.selectedPartitionIdx++
			}
		case 2:
			if m.selectedCronIdx < len(m.cronJobs)-1 {
				m.selectedCronIdx++
			}
		case 3:
			if m.systemUpdates != nil && m.selectedUpdateIdx < len(m.systemUpdates.UpdatesList)-1 {
				m.selectedUpdateIdx++
			}
		}
	case TabTerminal:
		switch m.terminalTabFocus {
		case 1:
			if m.selectedSSHIdx < len(m.sshProfiles)-1 {
				m.selectedSSHIdx++
			}
		case 2:
			if m.selectedEnvVarIdx < len(m.envVars)-1 {
				m.selectedEnvVarIdx++
			}
		case 3:
			if m.selectedJobIdx < len(m.queue)-1 {
				m.selectedJobIdx++
			}
		case 4:
			limit := m.getPluginsListLen()
			if m.selectedPluginIdx < limit-1 {
				m.selectedPluginIdx++
			}
		}
	}
}

func (m *Model) postNavCmds() []tea.Cmd {
	var cmds []tea.Cmd
	if m.activeTab == TabContainers && m.containersTabFocus == 0 && len(m.stacks) > 0 {
		cmds = append(cmds, m.fetchContainersCmd(m.stacks[m.selectedStackIdx].Name))
	}
	if m.activeTab == TabLogs && m.focusOnConfig && len(m.logContainers) > 0 {
		selectedC := m.logContainers[m.selectedLogContainerIdx]
		cname := "unknown"
		if len(selectedC.Names) > 0 {
			cname = strings.TrimPrefix(selectedC.Names[0], "/")
		}
		m.logs = ""
		m.logScrollOffset = 0
		cmds = append(cmds, m.fetchLogsCmd(cname))
	}
	return cmds
}

func (m *Model) handleEnterAction() []tea.Cmd {
	var cmds []tea.Cmd
	if m.activeTab == TabDashboard && m.dashboardFocusIndex == 2 {
		// Quick Actions: 0 = Scan Stacks, 1 = Prune Docker, 2 = System Reconcile
		switch m.selectedQuickActionIdx {
		case 0:
			m.SetNotification("🔍 Scanning stacks...", 5*time.Second)
			cmds = append(cmds, m.scanStacksCmd())
		case 1:
			m.SetNotification("🧹 Pruning unused Docker resources...", 10*time.Second)
			cmds = append(cmds, func() tea.Msg {
				imagesSize, _ := m.client.PruneDockerImages()
				volSize, _ := m.client.PruneDockerVolumes()
				netCount, _ := m.client.PruneDockerNetworks()
				totalGB := float64(imagesSize+volSize) / (1024 * 1024 * 1024)
				return actionResultMsg{message: fmt.Sprintf("✅ Pruned Docker! Space saved: %.2f GB. Networks removed: %d", totalGB, netCount)}
			})
		case 2:
			m.SetNotification("⚙️ Reconciling M3TAL State...", 8*time.Second)
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.ReconcileSystem()
				return actionResultMsg{message: "✅ Daemon state reconciled!", err: err}
			})
		}
	} else if m.activeTab == TabTerminal {
		if m.terminalTabFocus == 0 {
			// Embedded terminal launcher
			m.SetNotification("🚀 Suspending TUI. Launching local shell...", 2*time.Second)
			cmds = append(cmds, m.execShellCmd())
		} else if m.terminalTabFocus == 1 && len(m.sshProfiles) > 0 {
			// SSH profile connection launcher
			host := m.sshProfiles[m.selectedSSHIdx]
			m.SetNotification(fmt.Sprintf("🔑 Suspending TUI. Launching SSH session to %s...", host), 2*time.Second)
			cmds = append(cmds, m.execSSHCmd(host))
		}
	}
	return cmds
}

// --- Asynchronous API Commands ---

func (m Model) fetchActiveTabCmd() tea.Cmd {
	switch m.activeTab {
	case TabDashboard:
		return tea.Batch(
			m.fetchDetailedMetricsCmd(),
			m.fetchStatusCmd(),
			m.fetchDoctorReportCmd(),
		)
	case TabContainers:
		return tea.Batch(
			func() tea.Msg {
				stacks, err := m.client.GetStacks()
				return stacksMsg{stacks: stacks, err: err}
			},
			func() tea.Msg {
				containers, err := m.client.GetContainers()
				return containersMsg{containers: containers, err: err}
			},
			m.fetchDockerImagesCmd(),
			m.fetchDockerVolumesCmd(),
			m.fetchDockerNetworksCmd(),
		)
	case TabLogs:
		return func() tea.Msg {
			containers, err := m.client.GetContainers()
			return containersMsg{containers: containers, err: err}
		}
	case TabEditor:
		return func() tea.Msg {
			cfg, err := m.client.GetConfig()
			cloudflared, err2 := m.client.GetCloudflaredConfig()
			envRaw, err3 := m.client.GetEnvConfigRaw()
			var finalErr error
			if err != nil {
				finalErr = err
			} else if err2 != nil {
				finalErr = err2
			} else if err3 != nil {
				finalErr = err3
			}
			return configMsg{config: cfg, cloudflared: cloudflared, envRaw: envRaw, err: finalErr}
		}
	case TabSystem:
		return tea.Batch(
			m.fetchSystemServicesCmd(),
			m.fetchSystemStorageCmd(),
			m.fetchCronJobsCmd(),
			m.fetchSystemUpdatesCmd(),
		)
	case TabTerminal:
		return tea.Batch(
			func() tea.Msg {
				queue, err := m.client.GetQueue()
				modelsList, err2 := m.client.GetAIModels()
				var finalErr error
				if err != nil {
					finalErr = err
				} else if err2 != nil {
					finalErr = err2
				}
				return aiMsg{queue: queue, models: modelsList, err: finalErr}
			},
			func() tea.Msg {
				plugins, err := m.client.GetPlugins()
				catalog, err2 := m.client.GetPluginCatalog()
				var finalErr error
				if err != nil {
					finalErr = err
				} else if err2 != nil {
					finalErr = err2
				}
				return pluginsMsg{plugins: plugins, catalog: catalog, err: finalErr}
			},
		)
	}
	return nil
}

func (m Model) fetchAllDataCmd() tea.Cmd {
	return tea.Batch(
		m.fetchActiveTabCmd(),
		m.fetchStatusCmd(),
	)
}

func (m Model) fetchDetailedMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		stats, err := m.client.GetTrayStats()
		return detailedMetricsMsg{stats: stats, err: err}
	}
}

func (m Model) fetchDoctorReportCmd() tea.Cmd {
	return func() tea.Msg {
		report, err := m.client.GetDoctorReport()
		return doctorReportMsg{report: report, err: err}
	}
}

func (m Model) fetchDockerImagesCmd() tea.Cmd {
	return func() tea.Msg {
		images, err := m.client.GetDockerImages()
		return dockerImagesMsg{images: images, err: err}
	}
}

func (m Model) fetchDockerVolumesCmd() tea.Cmd {
	return func() tea.Msg {
		volumes, err := m.client.GetDockerVolumes()
		return dockerVolumesMsg{volumes: volumes, err: err}
	}
}

func (m Model) fetchDockerNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		networks, err := m.client.GetDockerNetworks()
		return dockerNetworksMsg{networks: networks, err: err}
	}
}

func (m Model) fetchSystemServicesCmd() tea.Cmd {
	return func() tea.Msg {
		services, err := m.client.GetSystemServices()
		return systemServicesMsg{services: services, err: err}
	}
}

func (m Model) fetchSystemStorageCmd() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.GetSystemStorage()
		if err != nil {
			return systemStorageMsg{err: err}
		}
		return systemStorageMsg{partitions: resp.Partitions, samba: resp.Samba, err: nil}
	}
}

func (m Model) fetchCronJobsCmd() tea.Cmd {
	return func() tea.Msg {
		jobs, err := m.client.GetCronJobs()
		return cronJobsMsg{jobs: jobs, err: err}
	}
}

func (m Model) fetchSystemUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		updates, err := m.client.GetSystemUpdates()
		return systemUpdatesMsg{updates: updates, err: err}
	}
}

func (m Model) fetchStatusCmd() tea.Cmd {
	return func() tea.Msg {
		status, err := m.client.GetStatus()
		return statusMsg{status: status, err: err}
	}
}

func (m Model) fetchContainersCmd(_ string) tea.Cmd {
	return func() tea.Msg {
		containers, err := m.client.GetContainers()
		return containersMsg{containers: containers, err: err}
	}
}

func (m Model) fetchLogsCmd(cname string) tea.Cmd {
	return func() tea.Msg {
		logs, err := m.client.GetLogs(cname)
		return logsMsg{container: cname, logs: logs, err: err}
	}
}

func (m Model) scanStacksCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := m.client.ScanStacks()
		return actionResultMsg{message: "✅ Stacks scanned & synced!", err: err}
	}
}

func (m Model) deployStackCmd(name string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.client.StartStack(name)
		return actionResultMsg{message: fmt.Sprintf("✅ Stack %s deployed successfully!", name), err: err}
	}
}

func (m Model) stopStackCmd(name string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.client.StopStack(name)
		return actionResultMsg{message: fmt.Sprintf("✅ Stack %s stopped successfully!", name), err: err}
	}
}

func (m Model) pullStackCmd(name string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.client.PullStacks(name)
		return actionResultMsg{message: fmt.Sprintf("✅ Stack %s images pulled!", name), err: err}
	}
}

func (m Model) controlContainerCmd(name, action string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.ControlContainer(name, action)
		return actionResultMsg{message: fmt.Sprintf("✅ Container %s successfully %sed!", name, action), err: err}
	}
}

func (m Model) controlSystemServiceCmd(name, action string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.ControlSystemService(name, action)
		return actionResultMsg{message: fmt.Sprintf("✅ Service %s successfully %sed!", name, action), err: err}
	}
}

func (m Model) applySystemUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.client.ApplySystemUpdates()
		return actionResultMsg{message: "✅ System package updates applied!", err: err}
	}
}

func (m Model) exportLogsCmd(cname string, logs string) tea.Cmd {
	return func() tea.Msg {
		path := fmt.Sprintf("/docker/logs-%s.txt", cname)
		err := os.WriteFile(path, []byte(logs), 0644)
		return actionResultMsg{message: fmt.Sprintf("✅ Logs exported to %s", path), err: err}
	}
}

// --- Suspension Interactive Subprocess Launchers ---

func (m Model) execContainerCmd(cname string) tea.Cmd {
	c := exec.Command("docker", "exec", "-it", cname, "bash")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			// Fallback to sh if bash not found
			c2 := exec.Command("docker", "exec", "-it", cname, "sh")
			return tea.ExecProcess(c2, func(err error) tea.Msg {
				return actionResultMsg{message: fmt.Sprintf("🐳 Returned from container %s shell", cname), err: err}
			})()
		}
		return actionResultMsg{message: fmt.Sprintf("🐳 Returned from container %s shell", cname), err: err}
	})
}

func (m Model) execShellCmd() tea.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}
	c := exec.Command(shell)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return actionResultMsg{message: "💻 Returned from host shell", err: err}
	})
}

func (m Model) execSSHCmd(host string) tea.Cmd {
	c := exec.Command("ssh", host)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return actionResultMsg{message: fmt.Sprintf("🔑 Returned from SSH session with %s", host), err: err}
	})
}

// --- Config Explorer ---

func (m *Model) refreshConfigFiles() {
	m.configFiles = []ConfigFile{}
	seen := make(map[string]bool)

	addFile := func(path, name string) {
		if seen[path] {
			return
		}
		seen[path] = true

		var icon, tag string
		switch {
		case path == system.GetConfigPath():
			icon = "🔧"
			tag = "[/etc/m3tal]"
		case strings.HasSuffix(strings.ToLower(name), "-compose.yml"):
			lowerName := strings.ToLower(name)
			isSystem := false
			keywords := []string{"m3tal", "ai", "maintenance", "network", "routing", "traefik"}
			for _, kw := range keywords {
				if strings.Contains(lowerName, kw) {
					isSystem = true
					break
				}
			}
			if isSystem {
				icon = "🖥️"
			} else {
				icon = "🐳"
			}
			tag = "[/docker]"
		default:
			icon = "📄"
			tag = "[/docker]"
		}

		displayName := icon + " " + name
		if tag != "" {
			displayName += " " + tag
		}
		m.configFiles = append(m.configFiles, ConfigFile{Name: displayName, Path: path})
	}

	addFile(system.GetConfigPath(), "System Config (.env)")

	const dockerDir = "/docker"
	entries, err := os.ReadDir(dockerDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			lower := strings.ToLower(name)
			if strings.HasSuffix(lower, ".env") ||
				strings.HasSuffix(lower, "-compose.yml") ||
				strings.Contains(lower, "config.yml") ||
				strings.Contains(lower, "config.yaml") {
				addFile(filepath.Join(dockerDir, name), name)
			}
		}
	}

	if len(m.configFiles) > 1 {
		sort.Slice(m.configFiles[1:], func(i, j int) bool {
			return strings.ToLower(m.configFiles[1:][i].Name) < strings.ToLower(m.configFiles[1:][j].Name)
		})
	}

	if m.selectedConfigIdx >= len(m.configFiles) {
		m.selectedConfigIdx = 0
	}
}

func (m *Model) loadSelectedConfigContent() {
	if len(m.configFiles) == 0 {
		m.selectedConfigContent = ""
		return
	}

	bytes, err := os.ReadFile(m.configFiles[m.selectedConfigIdx].Path)
	if err != nil {
		m.selectedConfigContent = fmt.Sprintf("Error reading file: %v", err)
		return
	}
	m.selectedConfigContent = string(bytes)
}

func (m Model) editConfigCmd() tea.Cmd {
	if len(m.configFiles) == 0 {
		return func() tea.Msg {
			return editFinishedMsg{err: fmt.Errorf("no config file to edit")}
		}
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	filePath := m.configFiles[m.selectedConfigIdx].Path
	c := exec.Command(editor, filePath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return editFinishedMsg{err: fmt.Errorf("editor failed: %w", err)}
		}

		updatedBytes, err := os.ReadFile(filePath)
		if err != nil {
			return editFinishedMsg{err: fmt.Errorf("failed to read updated file: %w", err)}
		}

		isEnvEdit := (filePath == system.GetConfigPath())
		isCloudflaredEdit := (filepath.Base(filePath) == "cloudflared-config.yml")

		if isEnvEdit {
			_ = m.client.UpdateEnvConfig(string(updatedBytes))
		} else if isCloudflaredEdit {
			_ = m.client.UpdateCloudflaredConfig(string(updatedBytes))
		}

		return editFinishedMsg{content: string(updatedBytes), isEnvEdit: isEnvEdit}
	})
}

// --- Plugins List Length ---

func (m Model) getPluginsListLen() int {
	if m.showCatalog {
		return len(m.catalog)
	}
	if m.plugins == nil {
		return 0
	}
	return len(m.plugins.Routes) + len(m.plugins.Stacks) + len(m.plugins.Middleware)
}

func (m Model) togglePluginCmd() tea.Cmd {
	return func() tea.Msg {
		var err error
		var pName, pKind string
		var isEnabled bool

		if m.showCatalog {
			if len(m.catalog) == 0 || m.selectedPluginIdx >= len(m.catalog) {
				return actionResultMsg{err: fmt.Errorf("no plugin selected")}
			}
			item := m.catalog[m.selectedPluginIdx]
			pName = item.Name
			pKind = item.Kind
			if item.Status == "not_installed" {
				return actionResultMsg{err: fmt.Errorf("plugin %s must be installed first", pName)}
			}
			isEnabled = (item.Status == "enabled")
		} else {
			if m.plugins == nil {
				return actionResultMsg{err: fmt.Errorf("plugins not loaded")}
			}
			var list []models.Plugin
			list = append(list, m.plugins.Routes...)
			list = append(list, m.plugins.Stacks...)
			list = append(list, m.plugins.Middleware...)

			if len(list) == 0 || m.selectedPluginIdx >= len(list) {
				return actionResultMsg{err: fmt.Errorf("no plugin selected")}
			}
			p := list[m.selectedPluginIdx]
			pName = p.Name
			pKind = p.Kind
			isEnabled = p.Enabled
		}

		if isEnabled {
			err = m.client.DisablePlugin(pName, pKind)
			return actionResultMsg{message: fmt.Sprintf("✅ Plugin %s disabled!", pName), err: err}
		}

		err = m.client.EnablePlugin(pName, pKind)
		return actionResultMsg{message: fmt.Sprintf("✅ Plugin %s enabled!", pName), err: err}
	}
}

func (m Model) installPluginCmd() tea.Cmd {
	return func() tea.Msg {
		if !m.showCatalog || len(m.catalog) == 0 || m.selectedPluginIdx >= len(m.catalog) {
			return actionResultMsg{err: fmt.Errorf("no plugin selected")}
		}
		item := m.catalog[m.selectedPluginIdx]
		if item.Status != "not_installed" {
			return actionResultMsg{message: fmt.Sprintf("Plugin %s is already installed", item.Name), err: nil}
		}

		err := m.client.InstallPlugin(item.Name, item.Kind)
		return actionResultMsg{message: fmt.Sprintf("✅ Plugin %s installed successfully!", item.Name), err: err}
	}
}
