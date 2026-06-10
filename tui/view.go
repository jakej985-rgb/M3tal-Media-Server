package tui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
)

// Version is injected at build time via:
//
//	go build -ldflags "-X github.com/jakej985-rgb/m3tal-core/tui.Version=$(cat VERSION)"
var Version string

// Common UI Styling using Lip Gloss
var (
	// Colors
	colorPanelBg       = lipgloss.Color("#181818")
	colorGreen         = lipgloss.Color("#00E676")
	colorYellow        = lipgloss.Color("#FFD600")
	colorRed           = lipgloss.Color("#FF1744")
	colorCyan          = lipgloss.Color("#00B0FF")
	colorGray          = lipgloss.Color("#90A4AE")
	colorDarkGray      = lipgloss.Color("#2A2A2A")
	colorHighlightText = lipgloss.Color("#FFFFFF")

	// Styles
	styleAppTitle = lipgloss.NewStyle().
			Foreground(colorHighlightText).
			Background(colorCyan).
			Bold(true).
			Padding(0, 1)

	styleTabActive = lipgloss.NewStyle().
			Foreground(colorGreen).
			Background(colorDarkGray).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(colorGreen)

	styleTabInactive = lipgloss.NewStyle().
				Foreground(colorGray).
				Padding(0, 1)

	styleFocusedPanel = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorGreen).
				Background(colorPanelBg).
				Padding(0, 1)

	styleUnfocusedPanel = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorDarkGray).
				Background(colorPanelBg).
				Padding(0, 1)

	styleHeaderLabel = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true)

	styleStatusHealthy = lipgloss.NewStyle().
				Foreground(colorHighlightText).
				Background(colorGreen).
				Bold(true).
				Padding(0, 1)

	styleStatusUnhealthy = lipgloss.NewStyle().
				Foreground(colorHighlightText).
				Background(colorRed).
				Bold(true).
				Padding(0, 1)

	styleNotification = lipgloss.NewStyle().
				Foreground(colorHighlightText).
				Background(colorYellow).
				Bold(true).
				Padding(0, 2)

	styleFooter = lipgloss.NewStyle().
			Foreground(colorGray).
			Background(colorDarkGray).
			Padding(0, 1)
)

// View renders the visual state of the TUI.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing M3TAL TUI..."
	}

	// Calculate panel heights
	headerHeight := 4
	footerHeight := 2
	notificationHeight := 0
	if m.notification != "" && time.Now().Before(m.notificationTimeout) {
		notificationHeight = 1
	}
	contentHeight := m.height - headerHeight - footerHeight - notificationHeight - 2

	var content string

	if m.err != nil && m.activeTab == TabStacks && len(m.stacks) == 0 {
		// Connection error placeholder
		content = m.renderConnectionError(m.width-4, contentHeight)
	} else {
		// Render active tab view
		switch m.activeTab {
		case TabStacks:
			content = m.viewStacks(contentHeight)
		case TabLogs:
			content = m.viewLogs(contentHeight)
		case TabAI:
			content = m.viewAI(contentHeight)
		case TabPlugins:
			content = m.viewPlugins(contentHeight)
		case TabConfig:
			content = m.viewConfig(contentHeight)
		}
	}

	// Build parts
	header := m.renderHeader()
	tabs := m.renderTabs()
	metrics := m.renderMetricsBar()
	footer := m.renderFooter()

	var builder strings.Builder
	builder.WriteString(header + "\n")
	builder.WriteString(tabs + "\n")

	// Add toast notification if active
	if notificationHeight > 0 {
		notificationStr := styleNotification.Render(fmt.Sprintf("🔔 %s", m.notification))
		builder.WriteString(notificationStr + "\n")
	}

	builder.WriteString(content + "\n")
	builder.WriteString(metrics + "\n")
	builder.WriteString(footer)

	return builder.String()
}

func (m Model) renderHeader() string {
	healthStr := styleStatusHealthy.Render("🟢 HEALTHY")
	if m.status != nil && m.status.Status == "unhealthy" {
		healthStr = styleStatusUnhealthy.Render("🔴 FAULT")
	}

	hostname := "unknown"
	uptimeStr := "0m"
	if m.metrics != nil {
		hostname = m.metrics.Hostname
		uptimeDur := time.Duration(m.metrics.Uptime) * time.Second
		days := int(uptimeDur.Hours()) / 24
		hours := int(uptimeDur.Hours()) % 24
		minutes := int(uptimeDur.Minutes()) % 60
		if days > 0 {
			uptimeStr = fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
		} else {
			uptimeStr = fmt.Sprintf("%dh %dm", hours, minutes)
		}
	}

	version := strings.TrimSpace(Version)
	if version == "" {
		if bytes, err := os.ReadFile("/etc/m3tal/VERSION"); err == nil {
			if v := strings.TrimSpace(string(bytes)); v != "" {
				version = v
			}
		}
	}
	if version == "" {
		version = "dev"
	}

	title := styleAppTitle.Render(fmt.Sprintf("M3TAL CONTROL CENTER (%s)", version))
	statusInfo := fmt.Sprintf("💻 Host: %s | Uptime: %s | Health: %s", hostname, uptimeStr, healthStr)

	// Calculate spaces to right-align the status info
	spaces := m.width - lipgloss.Width(title) - lipgloss.Width(statusInfo) - 2
	if spaces < 0 {
		spaces = 0
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, title, strings.Repeat(" ", spaces), statusInfo)
}

func (m Model) renderTabs() string {
	var tabs []string
	tabNames := []string{"[1] Stacks & Services", "[2] Container Logs", "[3] AI Queue & Models", "[4] Plugins Manager", "[5] System Config"}
	for i, name := range tabNames {
		if Tab(i) == m.activeTab {
			tabs = append(tabs, styleTabActive.Render(name))
		} else {
			tabs = append(tabs, styleTabInactive.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderConnectionError(width, height int) string {
	styleBlock := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorRed).
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	errMsgStr := fmt.Sprintf("\n\n🚨 [bold %s]API OFFLINE[/bold %s]\n\nConnection Error: %v\n\nChecking connection...", colorRed, colorRed, m.err)
	return styleBlock.Render(errMsgStr)
}

func (m Model) viewStacks(height int) string {
	leftWidth := 35
	rightWidth := m.width - leftWidth - 6 // Border padding

	// 1. Render Left panel: Stacks List
	leftStyle := styleUnfocusedPanel
	if m.focusOnStacks {
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderStacksList(leftWidth-4, height-2)
	leftBox := leftStyle.Width(leftWidth).Height(height).Render(leftContent)

	// 2. Render Right panel: Services List
	rightStyle := styleFocusedPanel
	if m.focusOnStacks {
		rightStyle = styleUnfocusedPanel
	}
	rightContent := m.renderServicesList(rightWidth-2, height-2)
	rightBox := rightStyle.Width(rightWidth).Height(height).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func (m Model) renderStacksList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📁 STACKS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.stacks) == 0 {
		builder.WriteString("\n  No stacks found.\n")
		return builder.String()
	}

	// Truncate to height limit
	visibleRows := height - 2
	for idx, s := range m.stacks {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.focusOnStacks && idx == m.selectedStackIdx {
			cursor = ">"
		}

		// Status coloring
		statusStr := strings.ToUpper(s.Status)
		var statusStyle lipgloss.Style
		switch s.Status {
		case "running", "success":
			statusStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		case "failed":
			statusStyle = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
		default:
			statusStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		}

		nameTrunc := truncateString(s.Name, width-15)
		row := fmt.Sprintf("%s %-15s [%s]", cursor, nameTrunc, statusStyle.Render(statusStr))
		builder.WriteString(row + "\n")
	}

	return builder.String()
}

func (m Model) renderServicesList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🐳 STACK SERVICES") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.stacks) == 0 {
		return builder.String()
	}
	selectedStackName := m.stacks[m.selectedStackIdx].Name

	if len(m.containers) == 0 {
		builder.WriteString(fmt.Sprintf("\n  No running services for stack: %s\n", selectedStackName))
		return builder.String()
	}

	// Print column headers
	builder.WriteString(fmt.Sprintf("  %-25s %-25s %-20s\n", "NAME", "IMAGE", "STATE (STATUS)"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, c := range m.containers {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if !m.focusOnStacks && idx == m.selectedContainerIdx {
			cursor = ">"
		}

		cname := "unknown"
		if len(c.Names) > 0 {
			cname = strings.TrimPrefix(c.Names[0], "/")
		}

		var stateStyle lipgloss.Style
		switch c.State {
		case "running":
			stateStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		case "exited", "dead":
			stateStyle = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
		default:
			stateStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		}

		nameT := truncateString(cname, 23)
		imgT := truncateString(c.Image, 23)
		stateT := fmt.Sprintf("%s (%s)", stateStyle.Render(strings.ToUpper(c.State)), c.Status)

		row := fmt.Sprintf("%s %-25s %-25s %-20s\n", cursor, nameT, imgT, stateT)
		builder.WriteString(row)
	}

	return builder.String()
}

func (m Model) viewLogs(height int) string {
	leftWidth := 30
	rightWidth := m.width - leftWidth - 6

	// 1. Render Left panel: Container list
	leftStyle := styleUnfocusedPanel
	if m.focusOnStacks { // in logs page we repurpose focus variable
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderLogsContainerList(leftWidth-4, height-2)
	leftBox := leftStyle.Width(leftWidth).Height(height).Render(leftContent)

	// 2. Render Right panel: Logs viewer
	rightStyle := styleFocusedPanel
	if m.focusOnStacks {
		rightStyle = styleUnfocusedPanel
	}
	rightContent := m.renderLogsStream(rightWidth-2, height-2)
	rightBox := rightStyle.Width(rightWidth).Height(height).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func (m Model) renderLogsContainerList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🔍 CONTAINERS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.logContainers) == 0 {
		builder.WriteString("\n  No containers found.\n")
		return builder.String()
	}

	visibleRows := height - 2
	for idx, c := range m.logContainers {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.focusOnStacks && idx == m.selectedLogContainerIdx {
			cursor = ">"
		}

		cname := "unknown"
		if len(c.Names) > 0 {
			cname = strings.TrimPrefix(c.Names[0], "/")
		}

		nameT := truncateString(cname, width-3)
		builder.WriteString(fmt.Sprintf("%s %s\n", cursor, nameT))
	}

	return builder.String()
}

func (m Model) renderLogsStream(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📄 LOG STREAM") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.logContainers) == 0 {
		return builder.String()
	}
	cname := "unknown"
	c := m.logContainers[m.selectedLogContainerIdx]
	if len(c.Names) > 0 {
		cname = strings.TrimPrefix(c.Names[0], "/")
	}

	// If connection offline or logs empty
	if m.logs == "" {
		builder.WriteString(fmt.Sprintf("\n  Loading logs for container: %s...\n", cname))
		return builder.String()
	}

	lines := strings.Split(m.logs, "\n")
	totalLines := len(lines)
	m.logScrollHeight = height - 3

	// Clamp scroll offset
	if m.logScrollOffset < 0 {
		m.logScrollOffset = 0
	}
	maxOffset := totalLines - m.logScrollHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.logScrollOffset > maxOffset {
		m.logScrollOffset = maxOffset
	}

	// Determine visible slice
	start := totalLines - m.logScrollHeight - m.logScrollOffset
	if start < 0 {
		start = 0
	}
	end := start + m.logScrollHeight
	if end > totalLines {
		end = totalLines
	}

	for i := start; i < end; i++ {
		line := lines[i]
		// Clean line wrap / truncate
		builder.WriteString(truncateString(line, width) + "\n")
	}

	// Print scrolling help
	scrollHelp := fmt.Sprintf("─ [Line %d-%d of %d] (Arrow Up/Down to scroll) ─", start+1, end, totalLines)
	if totalLines == 0 {
		scrollHelp = "─ No logs available ─"
	}
	builder.WriteString("\n" + lipgloss.NewStyle().Foreground(colorGray).Render(scrollHelp))

	return builder.String()
}

func (m Model) viewAI(height int) string {
	leftWidth := m.width * 7 / 10
	if leftWidth < 40 {
		leftWidth = 40
	}
	rightWidth := m.width - leftWidth - 6

	// 1. Render Left panel: Queue List
	leftStyle := styleUnfocusedPanel
	if m.focusOnQueue {
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderAIQueue(leftWidth-4, height-2)
	leftBox := leftStyle.Width(leftWidth).Height(height).Render(leftContent)

	// 2. Render Right panel: Ollama Models
	rightStyle := styleFocusedPanel
	if m.focusOnQueue {
		rightStyle = styleUnfocusedPanel
	}
	rightContent := m.renderAIModels(rightWidth-2, height-2)
	rightBox := rightStyle.Width(rightWidth).Height(height).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func (m Model) renderAIQueue(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🤖 AI JOB QUEUE") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.queue) == 0 {
		builder.WriteString("\n  No active or historical AI tasks in queue.\n")
		return builder.String()
	}

	// Print headers
	builder.WriteString(fmt.Sprintf("  %-12s %-16s %-10s %-20s %-12s\n", "JOB ID", "TYPE", "PRIORITY", "DETAILS", "STATUS"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, job := range m.queue {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.focusOnQueue && idx == m.selectedJobIdx {
			cursor = ">"
		}

		prioStr := "Normal"
		switch job.Priority {
		case 3:
			prioStr = "High"
		case 1:
			prioStr = "Low"
		}

		details := ""
		if job.Type == "ai_generation" {
			if prompt, ok := job.Payload["prompt"].(string); ok {
				details = prompt
			}
		} else {
			details = fmt.Sprintf("%v", job.Payload)
		}

		var statusStyle lipgloss.Style
		switch job.Status {
		case "completed":
			statusStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		case "failed":
			statusStyle = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
		default:
			statusStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		}

		idT := truncateString(job.ID, 12)
		typeT := truncateString(job.Type, 15)
		detailsT := truncateString(details, width-56)
		statusT := statusStyle.Render(strings.ToUpper(string(job.Status)))

		row := fmt.Sprintf("%s %-12s %-16s %-10s %-20s %-12s\n", cursor, idT, typeT, prioStr, detailsT, statusT)
		builder.WriteString(row)
	}

	return builder.String()
}

func (m Model) renderAIModels(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📦 OLLAMA MODELS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.aiModels) == 0 {
		builder.WriteString("\n  No Ollama models found.\n")
		return builder.String()
	}

	visibleRows := height - 2
	for idx, model := range m.aiModels {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if !m.focusOnQueue && idx == m.selectedModelIdx {
			cursor = ">"
		}

		modelT := truncateString(model, width-3)
		builder.WriteString(fmt.Sprintf("%s %s\n", cursor, modelT))
	}

	return builder.String()
}

func (m Model) viewPlugins(height int) string {
	width := m.width - 6
	panelStyle := styleFocusedPanel

	var content string
	if m.showCatalog {
		content = m.renderPluginCatalog(width-4, height-2)
	} else {
		content = m.renderPluginsList(width-4, height-2)
	}

	return panelStyle.Width(width).Height(height).Render(content)
}

func (m Model) renderPluginsList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🧩 LOADED PLUGINS (Press 'c' to view Catalog)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if m.plugins == nil {
		builder.WriteString("\n  Loading loaded plugins list...\n")
		return builder.String()
	}

	// Consolidate plugins list
	var list []models.Plugin
	list = append(list, m.plugins.Routes...)
	list = append(list, m.plugins.Stacks...)
	list = append(list, m.plugins.Middleware...)

	if len(list) == 0 {
		builder.WriteString("\n  No plugins currently loaded.\n")
		return builder.String()
	}

	// Columns headers
	builder.WriteString(fmt.Sprintf("  %-12s %-20s %-10s %-20s %-12s %-20s\n", "KIND", "NAME", "VERSION", "AUTHOR", "STATUS", "WARNINGS"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, p := range list {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if idx == m.selectedPluginIdx {
			cursor = ">"
		}

		statusStr := "[bold #FFD600]DISABLED[/bold #FFD600]"
		if p.Enabled {
			statusStr = "[bold #00E676]ENABLED[/bold #00E676]"
		}

		statusStyled := lipgloss.NewStyle().Render(statusStr)
		warningsStr := strings.Join(p.Warnings, ", ")
		if warningsStr != "" {
			warningsStr = lipgloss.NewStyle().Foreground(colorRed).Render("⚠️  " + warningsStr)
		}

		kindT := truncateString(p.Kind, 11)
		nameT := truncateString(p.Name, 19)
		versionT := truncateString(p.Version, 9)
		authorT := truncateString(p.Author, 19)
		warningsT := truncateString(warningsStr, width-80)

		row := fmt.Sprintf("%s %-12s %-20s %-10s %-20s %-12s %-20s\n", cursor, kindT, nameT, versionT, authorT, statusStyled, warningsT)
		builder.WriteString(row)
	}

	return builder.String()
}

func (m Model) renderPluginCatalog(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🧩 PLUGINS CATALOG (Press 'c' to view Installed)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.catalog) == 0 {
		builder.WriteString("\n  No catalog items found.\n")
		return builder.String()
	}

	// Column headers
	builder.WriteString(fmt.Sprintf("  %-12s %-20s %-10s %-20s %-15s %-30s\n", "KIND", "NAME", "VERSION", "AUTHOR", "INSTALL STATUS", "DESCRIPTION"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, item := range m.catalog {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if idx == m.selectedPluginIdx {
			cursor = ">"
		}

		var statusStyled string
		switch item.Status {
		case "enabled":
			statusStyled = lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("ENABLED")
		case "disabled":
			statusStyled = lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("DISABLED")
		default:
			statusStyled = lipgloss.NewStyle().Foreground(colorGray).Render("NOT INSTALLED")
		}

		kindT := truncateString(item.Kind, 11)
		nameT := truncateString(item.Name, 19)
		versionT := truncateString(item.Version, 9)
		authorT := truncateString(item.Author, 19)
		descT := truncateString(item.Description, width-83)

		row := fmt.Sprintf("%s %-12s %-20s %-10s %-20s %-15s %-30s\n", cursor, kindT, nameT, versionT, authorT, statusStyled, descT)
		builder.WriteString(row)
	}

	return builder.String()
}

func (m Model) viewConfig(height int) string {
	leftWidth := 35
	rightWidth := m.width - leftWidth - 6 // Border padding

	// 1. Render Left panel: Configurations List
	leftStyle := styleUnfocusedPanel
	if m.focusOnConfig {
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderConfigsList(leftWidth-4, height-2)
	leftBox := leftStyle.Width(leftWidth).Height(height).Render(leftContent)

	// 2. Render Right panel: Selected Config Viewer/Editor
	rightStyle := styleFocusedPanel
	if m.focusOnConfig {
		rightStyle = styleUnfocusedPanel
	}
	rightContent := m.renderSelectedConfigViewer(rightWidth-4, height-2)
	rightBox := rightStyle.Width(rightWidth).Height(height).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func (m Model) renderConfigsList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📁 CONFIGS/ENVS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.configFiles) == 0 {
		builder.WriteString("\n  No configs found.\n")
		return builder.String()
	}

	visibleRows := height - 2
	for idx, f := range m.configFiles {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if idx == m.selectedConfigIdx {
			cursor = ">"
		}

		nameTrunc := truncateString(f.Name, width-4)
		row := fmt.Sprintf("%s %s", cursor, nameTrunc)
		builder.WriteString(row + "\n")
	}

	return builder.String()
}

func (m Model) renderSelectedConfigViewer(width, height int) string {
	var builder strings.Builder

	if len(m.configFiles) == 0 {
		builder.WriteString(styleHeaderLabel.Render("⚙️ CONFIG VIEWER") + "\n")
		builder.WriteString(strings.Repeat("─", width) + "\n\n  No file selected.\n")
		return builder.String()
	}

	selectedFile := m.configFiles[m.selectedConfigIdx]
	builder.WriteString(styleHeaderLabel.Render(fmt.Sprintf("⚙️ %s", strings.ToUpper(selectedFile.Name))) + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if m.selectedConfigContent == "" {
		if m.err != nil {
			builder.WriteString(fmt.Sprintf("\n  Error loading config: %v\n", m.err))
		} else {
			builder.WriteString("\n  Loading/Reading configuration...\n")
		}
		return builder.String()
	}

	lines := strings.Split(m.selectedConfigContent, "\n")
	totalLines := len(lines)
	visibleRows := height - 3
	if visibleRows < 1 {
		visibleRows = 1
	}

	// Clamp scroll offset
	if m.configScrollOffset < 0 {
		m.configScrollOffset = 0
	}
	maxOffset := totalLines - visibleRows
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.configScrollOffset > maxOffset {
		m.configScrollOffset = maxOffset
	}

	endIdx := m.configScrollOffset + visibleRows
	if endIdx > totalLines {
		endIdx = totalLines
	}

	for idx := m.configScrollOffset; idx < endIdx; idx++ {
		builder.WriteString(truncateString(lines[idx], width) + "\n")
	}

	// Print scrolling/edit help
	scrollHelp := fmt.Sprintf("─ [Line %d-%d of %d] (Arrow Up/Down to scroll, 'e' to edit) ─", m.configScrollOffset+1, endIdx, totalLines)
	builder.WriteString("\n" + lipgloss.NewStyle().Foreground(colorGray).Render(scrollHelp))

	return builder.String()
}

func (m Model) renderMetricsBar() string {
	if m.metrics == nil {
		return "Metrics: loading..."
	}

	progBar := func(val float64) string {
		filled := int(math.Round(val / 10.0))
		if filled < 0 {
			filled = 0
		}
		if filled > 10 {
			filled = 10
		}
		empty := 10 - filled
		return fmt.Sprintf("[%s%s] %.1f%%", strings.Repeat("█", filled), strings.Repeat("░", empty), val)
	}

	cpuStr := progBar(m.metrics.CPUUsage)
	ramStr := progBar(m.metrics.MemoryUsage)
	diskStr := progBar(m.metrics.DiskUsage)

	return fmt.Sprintf("🔥 CPU: %s  |  🧠 RAM: %s  |  💿 Disk: %s", cpuStr, ramStr, diskStr)
}

func (m Model) renderFooter() string {
	var keys []string
	switch m.activeTab {
	case TabStacks:
		keys = []string{"1-5: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate", "a: Scan Stacks", "u: Deploy Stack", "d: Stop Stack", "s: Start Svc", "x: Stop Svc", "t: Restart Svc", "q: Quit"}
	case TabLogs:
		keys = []string{"1-5: Switch Tabs", "Tab: Select Panel", "Arrows: Navigate/Scroll Logs", "r: Force Refresh", "q: Quit"}
	case TabAI:
		keys = []string{"1-5: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate", "k: Cancel Job", "q: Quit"}
	case TabPlugins:
		keys = []string{"1-5: Switch Tabs", "Arrows: Navigate", "c: Toggle Catalog", "i: Install", "e: Toggle Enable", "q: Quit"}
	case TabConfig:
		keys = []string{"1-5: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate/Scroll", "e: Edit Config", "r: Force Refresh", "q: Quit"}
	}
	return styleFooter.Width(m.width).Render("🔑 Keys: " + strings.Join(keys, " | "))
}

func truncateString(s string, limit int) string {
	if limit <= 3 {
		return "..."
	}
	if len(s) > limit {
		return s[:limit-3] + "..."
	}
	return s
}
