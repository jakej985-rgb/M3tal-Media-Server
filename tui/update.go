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
		// Periodic refresh
		cmds = append(cmds, m.fetchAllDataCmd(), tick())

	case metricsMsg:
		if msg.err == nil {
			m.metrics = msg.metrics
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

	case stacksMsg:
		if msg.err == nil {
			m.stacks = msg.stacks
			m.err = nil
			// Clamp stack selection index
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
			// Update containers in Stacks view (filtered by selected stack project label)
			if m.activeTab == TabStacks && m.selectedStackIdx < len(m.stacks) {
				stackName := strings.ToLower(m.stacks[m.selectedStackIdx].Name)
				trimmedStackName := strings.TrimLeft(stackName, "0123456789-")
				var filtered []models.Container
				for _, c := range msg.containers {
					proj := strings.ToLower(c.Labels["com.docker.compose.project"])
					mStack := strings.ToLower(c.Labels["m3tal.stack"])
					configFiles := strings.ToLower(c.Labels["com.docker.compose.project.config_files"])

					// Match if:
					// 1. Project equals stack name or trimmed stack name
					// 2. m3tal.stack label matches
					// 3. com.docker.compose.project.config_files contains the compose filename
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
			if m.activeTab == TabLogs {
				m.logContainers = msg.containers
				if m.selectedLogContainerIdx >= len(m.logContainers) {
					m.selectedLogContainerIdx = 0
				}
				if len(m.logContainers) > 0 {
					selectedC := m.logContainers[m.selectedLogContainerIdx]
					cname := "unknown"
					if len(selectedC.Names) > 0 {
						cname = strings.TrimPrefix(selectedC.Names[0], "/")
					}
					cmds = append(cmds, m.fetchLogsCmd(cname))
				}
			}
			m.err = nil
		} else {
			m.err = msg.err
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
			m.configData = msg.config
			m.cloudflaredContent = msg.cloudflared
			m.envRawContent = msg.envRaw
			m.err = nil
			var keys []string
			for k := range m.configData {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			m.configKeys = keys

			m.refreshConfigFiles()
			m.loadSelectedConfigContent()
		} else {
			m.err = msg.err
		}

	case editFinishedMsg:
		if msg.err != nil {
			m.SetNotification(fmt.Sprintf("❌ Edit failed: %v", msg.err), 4*time.Second)
		} else {
			m.SetNotification("✅ Configuration file updated!", 3*time.Second)
			m.loadSelectedConfigContent()
		}
		cmds = append(cmds, m.fetchAllDataCmd())

	case actionResultMsg:
		if msg.err != nil {
			m.SetNotification(fmt.Sprintf("❌ Error: %v", msg.err), 4*time.Second)
		} else {
			m.SetNotification(msg.message, 3*time.Second)
		}
		// Refresh everything on action completion
		cmds = append(cmds, m.fetchAllDataCmd())

	case tea.KeyMsg:
		keyStr := strings.ToLower(msg.String())
		switch keyStr {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			m.SetNotification("Refreshing dashboard data...", 1*time.Second)
			cmds = append(cmds, m.fetchAllDataCmd())

		case "1":
			m.activeTab = TabStacks
			cmds = append(cmds, m.fetchAllDataCmd())
		case "2":
			m.activeTab = TabLogs
			m.logs = ""
			m.logScrollOffset = 0
			cmds = append(cmds, m.fetchAllDataCmd())
		case "3":
			m.activeTab = TabAI
			cmds = append(cmds, m.fetchAllDataCmd())
		case "4":
			m.activeTab = TabPlugins
			cmds = append(cmds, m.fetchAllDataCmd())
		case "5":
			m.activeTab = TabConfig
			m.refreshConfigFiles()
			m.loadSelectedConfigContent()
			cmds = append(cmds, m.fetchAllDataCmd())

		case "tab":
			switch m.activeTab {
			case TabStacks:
				m.focusOnStacks = !m.focusOnStacks
			case TabLogs:
				m.focusOnStacks = !m.focusOnStacks // repurposed for logs left/right
			case TabAI:
				m.focusOnQueue = !m.focusOnQueue
			case TabConfig:
				m.focusOnConfig = !m.focusOnConfig
			}

		case "up", "k":
			if keyStr == "k" && m.activeTab == TabAI && m.focusOnQueue && len(m.queue) > 0 {
				job := m.queue[m.selectedJobIdx]
				m.SetNotification(fmt.Sprintf("✕ Cancelling AI Job: %s...", job.ID), 5*time.Second)
				cmds = append(cmds, m.cancelJobCmd(job.ID))
			} else {
				m.handleNavigationUp()
				if m.activeTab == TabStacks && m.focusOnStacks && len(m.stacks) > 0 {
					cmds = append(cmds, m.fetchContainersCmd(m.stacks[m.selectedStackIdx].Name))
				}
				if m.activeTab == TabLogs && m.focusOnStacks && len(m.logContainers) > 0 {
					selectedC := m.logContainers[m.selectedLogContainerIdx]
					cname := "unknown"
					if len(selectedC.Names) > 0 {
						cname = strings.TrimPrefix(selectedC.Names[0], "/")
					}
					m.logs = ""
					m.logScrollOffset = 0
					cmds = append(cmds, m.fetchLogsCmd(cname))
				}
			}

		case "down", "j":
			m.handleNavigationDown()
			if m.activeTab == TabStacks && m.focusOnStacks && len(m.stacks) > 0 {
				cmds = append(cmds, m.fetchContainersCmd(m.stacks[m.selectedStackIdx].Name))
			}
			if m.activeTab == TabLogs && m.focusOnStacks && len(m.logContainers) > 0 {
				selectedC := m.logContainers[m.selectedLogContainerIdx]
				cname := "unknown"
				if len(selectedC.Names) > 0 {
					cname = strings.TrimPrefix(selectedC.Names[0], "/")
				}
				m.logs = ""
				m.logScrollOffset = 0
				cmds = append(cmds, m.fetchLogsCmd(cname))
			}

		// --- Stack Action Hotkeys ---
		case "a":
			if m.activeTab == TabStacks {
				m.SetNotification("🔍 Scanning for stack updates...", 5*time.Second)
				cmds = append(cmds, m.scanStacksCmd())
			}
		case "u":
			if m.activeTab == TabStacks && len(m.stacks) > 0 {
				stackName := m.stacks[m.selectedStackIdx].Name
				m.SetNotification(fmt.Sprintf("🚀 Deploying stack %s...", stackName), 10*time.Second)
				cmds = append(cmds, m.deployStackCmd(stackName))
			}
		case "d":
			if m.activeTab == TabStacks && len(m.stacks) > 0 {
				stackName := m.stacks[m.selectedStackIdx].Name
				m.SetNotification(fmt.Sprintf("🧹 Stopping stack %s...", stackName), 10*time.Second)
				cmds = append(cmds, m.stopStackCmd(stackName))
			}

		// --- Container Control Action Hotkeys ---
		case "s":
			if m.activeTab == TabStacks && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🐳 Starting container %s...", cname), 5*time.Second)
				cmds = append(cmds, m.controlContainerCmd(cname, "start"))
			}
		case "x":
			if m.activeTab == TabStacks && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🛑 Stopping container %s...", cname), 5*time.Second)
				cmds = append(cmds, m.controlContainerCmd(cname, "stop"))
			}
		case "t":
			if m.activeTab == TabStacks && len(m.containers) > 0 {
				c := m.containers[m.selectedContainerIdx]
				cname := strings.TrimPrefix(c.Names[0], "/")
				m.SetNotification(fmt.Sprintf("🔄 Restarting container %s...", cname), 5*time.Second)
				cmds = append(cmds, m.controlContainerCmd(cname, "restart"))
			}

		// --- Plugin Action Hotkeys ---
		case "c":
			switch m.activeTab {
			case TabPlugins:
				m.showCatalog = !m.showCatalog
				m.selectedPluginIdx = 0
				m.SetNotification("Toggled plugins view", 1*time.Second)
			}
		case "e":
			switch m.activeTab {
			case TabPlugins:
				cmds = append(cmds, m.togglePluginCmd())
			case TabConfig:
				cmds = append(cmds, m.editConfigCmd())
			}
		case "i":
			if m.activeTab == TabPlugins && m.showCatalog {
				cmds = append(cmds, m.installPluginCmd())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleNavigationUp() {
	switch m.activeTab {
	case TabStacks:
		if m.focusOnStacks {
			if m.selectedStackIdx > 0 {
				m.selectedStackIdx--
			}
		} else {
			if m.selectedContainerIdx > 0 {
				m.selectedContainerIdx--
			}
		}
	case TabLogs:
		if m.focusOnStacks { // Left list
			if m.selectedLogContainerIdx > 0 {
				m.selectedLogContainerIdx--
			}
		} else { // Scroll logs
			m.logScrollOffset++
		}
	case TabAI:
		if m.focusOnQueue {
			if m.selectedJobIdx > 0 {
				m.selectedJobIdx--
			}
		} else {
			if m.selectedModelIdx > 0 {
				m.selectedModelIdx--
			}
		}
	case TabPlugins:
		if m.selectedPluginIdx > 0 {
			m.selectedPluginIdx--
		}
	case TabConfig:
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
	}
}

func (m *Model) handleNavigationDown() {
	switch m.activeTab {
	case TabStacks:
		if m.focusOnStacks {
			if m.selectedStackIdx < len(m.stacks)-1 {
				m.selectedStackIdx++
			}
		} else {
			if m.selectedContainerIdx < len(m.containers)-1 {
				m.selectedContainerIdx++
			}
		}
	case TabLogs:
		if m.focusOnStacks { // Left list
			if m.selectedLogContainerIdx < len(m.logContainers)-1 {
				m.selectedLogContainerIdx++
			}
		} else { // Scroll logs
			if m.logScrollOffset > 0 {
				m.logScrollOffset--
			}
		}
	case TabAI:
		if m.focusOnQueue {
			if m.selectedJobIdx < len(m.queue)-1 {
				m.selectedJobIdx++
			}
		} else {
			if m.selectedModelIdx < len(m.aiModels)-1 {
				m.selectedModelIdx++
			}
		}
	case TabPlugins:
		limit := m.getPluginsListLen()
		if m.selectedPluginIdx < limit-1 {
			m.selectedPluginIdx++
		}
	case TabConfig:
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
	}
}

func (m Model) getPluginsListLen() int {
	if m.showCatalog {
		return len(m.catalog)
	}
	if m.plugins == nil {
		return 0
	}
	return len(m.plugins.Routes) + len(m.plugins.Stacks) + len(m.plugins.Middleware)
}

// --- Asynchronous API Commands ---

func (m Model) fetchAllDataCmd() tea.Cmd {
	return tea.Batch(
		m.fetchMetricsCmd(),
		m.fetchStatusCmd(),
		m.fetchActiveTabCmd(),
	)
}

func (m Model) fetchMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		metrics, err := m.client.GetStats()
		return metricsMsg{metrics: metrics, err: err}
	}
}

func (m Model) fetchStatusCmd() tea.Cmd {
	return func() tea.Msg {
		status, err := m.client.GetStatus()
		return statusMsg{status: status, err: err}
	}
}

func (m Model) fetchActiveTabCmd() tea.Cmd {
	switch m.activeTab {
	case TabStacks:
		return func() tea.Msg {
			stacks, err := m.client.GetStacks()
			return stacksMsg{stacks: stacks, err: err}
		}
	case TabLogs:
		return func() tea.Msg {
			containers, err := m.client.GetContainers()
			return containersMsg{containers: containers, err: err}
		}
	case TabAI:
		return func() tea.Msg {
			queue, err := m.client.GetQueue()
			modelsList, err2 := m.client.GetAIModels()

			// Consolidate error
			var finalErr error
			if err != nil {
				finalErr = err
			} else if err2 != nil {
				finalErr = err2
			}
			return aiMsg{queue: queue, models: modelsList, err: finalErr}
		}
	case TabPlugins:
		return func() tea.Msg {
			plugins, err := m.client.GetPlugins()
			catalog, err2 := m.client.GetPluginCatalog()

			var finalErr error
			if err != nil {
				finalErr = err
			} else if err2 != nil {
				finalErr = err2
			}
			return pluginsMsg{plugins: plugins, catalog: catalog, err: finalErr}
		}
	case TabConfig:
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
	}
	return nil
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

func (m Model) controlContainerCmd(name, action string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.ControlContainer(name, action)
		return actionResultMsg{message: fmt.Sprintf("✅ Container %s successfully %sed!", name, action), err: err}
	}
}

func (m Model) cancelJobCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.CancelQueueJob(id)
		return actionResultMsg{message: fmt.Sprintf("✅ Job %s cancelled successfully!", id), err: err}
	}
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

func (m *Model) refreshConfigFiles() {
	m.configFiles = []ConfigFile{}

	// 1. Add the primary system config .env
	envPath := system.GetConfigPath()
	m.configFiles = append(m.configFiles, ConfigFile{
		Name: "System Config (.env)",
		Path: envPath,
	})

	// 2. Scan GetStackDir() for other *.env and *config.yml / *config.yaml files
	stackDir := system.GetStackDir()
	entries, err := os.ReadDir(stackDir)
	if err == nil {
		var found []ConfigFile
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			lower := strings.ToLower(name)

			// Skip if it matches the main envPath to avoid duplication
			absPath := filepath.Clean(filepath.Join(stackDir, name))
			absEnvPath, _ := filepath.Abs(envPath)
			absFilePath, _ := filepath.Abs(absPath)
			if absEnvPath == absFilePath {
				continue
			}

			// Match .env files or config.yml / config.yaml
			if strings.HasSuffix(lower, ".env") || strings.Contains(lower, "config.yml") || strings.Contains(lower, "config.yaml") {
				found = append(found, ConfigFile{
					Name: name,
					Path: absPath,
				})
			}
		}

		// Sort found files alphabetically
		sort.Slice(found, func(i, j int) bool {
			return strings.ToLower(found[i].Name) < strings.ToLower(found[j].Name)
		})

		m.configFiles = append(m.configFiles, found...)
	}

	// Clamp selectedConfigIdx
	if m.selectedConfigIdx >= len(m.configFiles) {
		m.selectedConfigIdx = 0
	}
}

func (m *Model) loadSelectedConfigContent() {
	if len(m.configFiles) == 0 {
		m.selectedConfigContent = ""
		return
	}

	path := m.configFiles[m.selectedConfigIdx].Path
	bytes, err := os.ReadFile(path)
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

		// Sync with client API if it's the main .env or cloudflared-config.yml
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
