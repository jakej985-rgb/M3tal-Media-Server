package tui

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jakej985-rgb/m3tal-core/pkg/models"
)

// Version is injected at build time
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

	if m.width < 80 || m.height < 18 {
		return fmt.Sprintf("\n  ⚠️  Terminal size too small (%dx%d)\n  Minimum required size: 80x18\n  Please enlarge your terminal window.", m.width, m.height)
	}

	header := m.renderHeader()
	tabs := m.renderTabs()
	metrics := m.renderMetricsBar()
	footer := m.renderFooter()

	notificationHeight := 0
	var notificationStr string
	if m.notification != "" && time.Now().Before(m.notificationTimeout) {
		notificationStr = styleNotification.Render(fmt.Sprintf("🔔 %s", m.notification))
		notificationHeight = lipgloss.Height(notificationStr)
	}

	headerHeight := lipgloss.Height(header)
	tabsHeight := lipgloss.Height(tabs)
	metricsHeight := lipgloss.Height(metrics)
	footerHeight := lipgloss.Height(footer)

	// Spacing newlines count:
	// 1 after header, 1 after tabs, 1 after notification (if present), 1 after content, 1 after metrics.
	spacingLines := 4
	if notificationHeight > 0 {
		spacingLines = 5
	}

	contentHeight := m.height - (headerHeight + tabsHeight + notificationHeight + metricsHeight + footerHeight + spacingLines)
	if contentHeight < 4 {
		contentHeight = 4
	}

	var content string

	if m.err != nil && m.activeTab == TabDashboard && m.detailedStats == nil {
		content = m.renderConnectionError(m.width-4, contentHeight)
	} else {
		switch m.activeTab {
		case TabDashboard:
			content = m.viewDashboard(contentHeight)
		case TabContainers:
			content = m.viewContainers(contentHeight)
		case TabLogs:
			content = m.viewLogs(contentHeight)
		case TabEditor:
			content = m.viewEditor(contentHeight)
		case TabSystem:
			content = m.viewSystem(contentHeight)
		case TabTerminal:
			content = m.viewTerminal(contentHeight)
		}
	}

	var builder strings.Builder
	builder.WriteString(header + "\n")
	builder.WriteString(tabs + "\n")

	if notificationHeight > 0 {
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
	if m.detailedStats != nil {
		hostname = m.detailedStats.Hostname
		uptimeDur := time.Duration(m.detailedStats.Uptime) * time.Second
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

	spaces := m.width - lipgloss.Width(title) - lipgloss.Width(statusInfo) - 2
	if spaces < 0 {
		spaces = 0
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, title, strings.Repeat(" ", spaces), statusInfo)
}

func (m Model) renderTabs() string {
	var tabs []string
	tabNames := []string{
		"[1] Dashboard",
		"[2] Containers & Docker",
		"[3] Container Logs",
		"[4] Config Editor",
		"[5] System Admin",
		"[6] Shell & Term",
	}
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
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center)

	errMsgStr := fmt.Sprintf("\n\n🚨 [bold %s]API OFFLINE[/bold %s]\n\nConnection Error: %v\n\nChecking connection...", colorRed, colorRed, m.err)
	return fitHeight(styleBlock.Render(errMsgStr), height, true)
}

// --- Tab 1: Dashboard View ---

func (m Model) viewDashboard(height int) string {
	leftWidth := 42
	if m.width < 90 {
		leftWidth = 38
	}
	rightWidth := m.width - leftWidth - 6

	// Left Side: Resource Telemetry
	leftStyle := styleUnfocusedPanel
	if m.dashboardFocusIndex == 0 {
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderDashboardTelemetry(leftWidth-4, height-2)
	leftBox := fitHeight(leftStyle.Width(leftWidth).Height(height-2).Render(leftContent), height, true)

	// Right Side: Docker & Doctor Reports & Quick Actions
	rightBox := m.renderDashboardDetails(rightWidth-2, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func (m Model) renderDashboardTelemetry(width, _ int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("⚡ SYSTEM METRICS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if m.detailedStats == nil {
		builder.WriteString("  Loading resource telemetry...\n")
		return builder.String()
	}

	coloredBar := func(pct float64, color lipgloss.Color) string {
		filled := int(math.Round(pct / 10.0))
		if filled < 0 {
			filled = 0
		}
		if filled > 10 {
			filled = 10
		}
		empty := 10 - filled
		filledStr := strings.Repeat("█", filled)
		emptyStr := strings.Repeat("░", empty)
		styledFilled := lipgloss.NewStyle().Foreground(color).Render(filledStr)
		return fmt.Sprintf("[%s%s]", styledFilled, emptyStr)
	}

	type telemetryEntry struct {
		label    string
		hasBar   bool
		barPct   float64
		barColor lipgloss.Color
		suffix   string
		isText   bool
		text     string
	}

	var entries []telemetryEntry

	addBar := func(label string, pct float64, color lipgloss.Color, suffix string) {
		entries = append(entries, telemetryEntry{
			label:    strings.ReplaceAll(label, "\ufe0f", ""),
			hasBar:   true,
			barPct:   pct,
			barColor: color,
			suffix:   suffix,
		})
	}

	addText := func(label string, text string) {
		entries = append(entries, telemetryEntry{
			label:  strings.ReplaceAll(label, "\ufe0f", ""),
			isText: true,
			text:   text,
		})
	}

	// 1. CPU Usage
	cpuVal := m.detailedStats.CPUUsage
	var cpuColor lipgloss.Color
	if cpuVal < 50.0 {
		cpuColor = colorGreen
	} else if cpuVal < 85.0 {
		cpuColor = colorYellow
	} else {
		cpuColor = colorRed
	}
	addBar("💻 CPU:", cpuVal, cpuColor, fmt.Sprintf("%.1f%%", cpuVal))

	// 2. CPU Temp
	if m.detailedStats.CPUTemp > 0 {
		tempC := m.detailedStats.CPUTemp
		tempF := tempC*1.8 + 32.0
		var tempColor lipgloss.Color
		if tempC < 60.0 {
			tempColor = colorGreen
		} else if tempC < 80.0 {
			tempColor = colorYellow
		} else {
			tempColor = colorRed
		}
		addBar("💻 CPU 🌡️:", tempC, tempColor, fmt.Sprintf("%.1f°F", tempF))
	}

	// 3. RAM Usage
	ramVal := m.detailedStats.MemoryUsage
	var ramColor lipgloss.Color
	if ramVal < 70.0 {
		ramColor = colorGreen
	} else if ramVal < 90.0 {
		ramColor = colorYellow
	} else {
		ramColor = colorRed
	}
	addBar("🧠 RAM:", ramVal, ramColor, fmt.Sprintf("%.1f%%", ramVal))

	// 4. RAM Stats
	var ramTextVal string
	if m.detailedStats.MemoryFrequency != "" && m.detailedStats.MemoryFrequency != "Unknown" {
		ramTextVal = fmt.Sprintf("[%.1f/%.1f GB]  [%s]", m.detailedStats.MemoryUsed, m.detailedStats.MemoryTotal, m.detailedStats.MemoryFrequency)
	} else {
		ramTextVal = fmt.Sprintf("[%.1f/%.1f GB]", m.detailedStats.MemoryUsed, m.detailedStats.MemoryTotal)
	}
	addText("🧠 RAM 📊:", ramTextVal)

	// 5. Disk Partitions
	if len(m.detailedStats.DiskPartitions) > 0 {
		for _, p := range m.detailedStats.DiskPartitions {
			freePct := 100.0 - p.UsedPercent
			var diskColor lipgloss.Color
			if freePct > 30.0 {
				diskColor = colorGreen
			} else if freePct >= 10.0 {
				diskColor = colorYellow
			} else {
				diskColor = colorRed
			}
			label := p.Label
			if label == "" {
				label = filepath.Base(p.Device)
			}
			addBar(fmt.Sprintf("💿 %s:", label), freePct, diskColor, fmt.Sprintf("%.1f%%", p.UsedPercent))
		}
	} else {
		freePct := 100.0 - m.detailedStats.DiskUsage
		var diskColor lipgloss.Color
		if freePct > 30.0 {
			diskColor = colorGreen
		} else if freePct >= 10.0 {
			diskColor = colorYellow
		} else {
			diskColor = colorRed
		}
		addBar("💿 root:", freePct, diskColor, fmt.Sprintf("%.1f%%", m.detailedStats.DiskUsage))
	}

	// 6. GPU Stats
	showGPU := m.detailedStats.GPUModel != "No GPU Detected" && m.detailedStats.GPUModel != ""
	if showGPU {
		addText(styleHeaderLabel.Render("📟 GPU: "+m.detailedStats.GPUModel), "")

		// GPU Core
		gpuVal := m.detailedStats.GPUUsage
		var gpuColor lipgloss.Color
		if gpuVal < 50.0 {
			gpuColor = colorGreen
		} else if gpuVal < 85.0 {
			gpuColor = colorYellow
		} else {
			gpuColor = colorRed
		}
		addBar("📟 Core:", gpuVal, gpuColor, fmt.Sprintf("%.1f%%", gpuVal))

		// GPU Temp
		if m.detailedStats.GPUTemp > 0 {
			gpuTempC := m.detailedStats.GPUTemp
			gpuTempF := gpuTempC*1.8 + 32.0
			var gpuTempColor lipgloss.Color
			if gpuTempC < 60.0 {
				gpuTempColor = colorGreen
			} else if gpuTempC < 80.0 {
				gpuTempColor = colorYellow
			} else {
				gpuTempColor = colorRed
			}
			addBar("📟 GPU 🌡️:", gpuTempC, gpuTempColor, fmt.Sprintf("%.1f°F", gpuTempF))
		}

		// GPU VRAM
		vramPct := 0.0
		if m.detailedStats.GPUMemTotal > 0 {
			vramPct = (m.detailedStats.GPUMemUsed / m.detailedStats.GPUMemTotal) * 100.0
		}
		var vramColor lipgloss.Color
		if vramPct < 70.0 {
			vramColor = colorGreen
		} else if vramPct < 90.0 {
			vramColor = colorYellow
		} else {
			vramColor = colorRed
		}
		addBar("📟 VRAM:", vramPct, vramColor, fmt.Sprintf("%.1f%%", vramPct))
	}

	// Calculate max width for entries with a bar
	maxW := 0
	for _, e := range entries {
		if e.hasBar {
			w := lipgloss.Width(e.label)
			if w > maxW {
				maxW = w
			}
		}
	}

	// Render entries
	for _, e := range entries {
		if e.hasBar {
			padding := strings.Repeat(" ", maxW-lipgloss.Width(e.label)+2)
			builder.WriteString(fmt.Sprintf("%s%s%s  %s\n", e.label, padding, coloredBar(e.barPct, e.barColor), e.suffix))
		} else {
			if e.text == "" {
				builder.WriteString(e.label + "\n")
			} else {
				builder.WriteString(fmt.Sprintf("%s %s\n", e.label, e.text))
			}
		}
	}

	return builder.String()
}

func (m Model) renderDashboardDetails(width, height int) string {
	if m.preflightExpanded {
		topHeight := height - 4
		docStyle := styleFocusedPanel
		docContent := m.renderDashboardDoctor(width-4, topHeight)
		topBox := fitHeight(docStyle.Width(width).Height(topHeight).Render(docContent), topHeight+2, true)
		return topBox
	}

	topHeight := (height - 4) / 2
	bottomHeight := height - 4 - topHeight

	// Doctor Report Panel
	docStyle := styleUnfocusedPanel
	if m.dashboardFocusIndex == 1 {
		docStyle = styleFocusedPanel
	}
	docContent := m.renderDashboardDoctor(width-4, topHeight)
	topBox := fitHeight(docStyle.Width(width).Height(topHeight).Render(docContent), topHeight+2, true)

	// Quick Actions Panel
	actStyle := styleUnfocusedPanel
	if m.dashboardFocusIndex == 2 {
		actStyle = styleFocusedPanel
	}
	actContent := m.renderDashboardQuickActions(width-4, bottomHeight)
	bottomBox := fitHeight(actStyle.Width(width).Height(bottomHeight).Render(actContent), bottomHeight+2, true)

	return lipgloss.JoinVertical(lipgloss.Left, topBox, bottomBox)
}

func (m Model) renderDashboardDoctor(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🩺 PRE-FLIGHT DIAGNOSTIC ALERTS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	hasContainers := len(m.doctorReport.Containers) > 0
	hasMounts := len(m.doctorReport.Mounts) > 0
	hasPorts := len(m.doctorReport.Ports) > 0

	if !hasContainers && !hasMounts && !hasPorts {
		builder.WriteString("\n  No active anomalies or Doctor report generated. Run Doctor scan below!\n")
		return builder.String()
	}

	visibleRows := height - 2
	count := 0

	// 1. Containers anomalies
	for _, c := range m.doctorReport.Containers {
		if c.Severity != models.SeverityPass {
			if count >= visibleRows {
				break
			}
			statusColor := colorYellow
			if c.Severity == models.SeverityFail {
				statusColor = colorRed
			}
			statusStyled := lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render(strings.ToUpper(string(c.Severity)))
			msg := c.Recommendation
			if msg == "" {
				msg = fmt.Sprintf("State: %s (Health: %s)", c.State, c.Health)
			}
			row := fmt.Sprintf("  [%s] Container %s: %s", statusStyled, c.Name, msg)
			builder.WriteString(truncateString(row, width) + "\n")
			count++
		}
	}

	// 2. Mounts anomalies
	for _, mt := range m.doctorReport.Mounts {
		if mt.Severity != models.SeverityPass {
			if count >= visibleRows {
				break
			}
			statusColor := colorYellow
			if mt.Severity == models.SeverityFail {
				statusColor = colorRed
			}
			statusStyled := lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render(strings.ToUpper(string(mt.Severity)))
			row := fmt.Sprintf("  [%s] Mount %s: %s", statusStyled, mt.Target, mt.Issue)
			builder.WriteString(truncateString(row, width) + "\n")
			count++
		}
	}

	// 3. Ports anomalies
	for _, pt := range m.doctorReport.Ports {
		if pt.Severity != models.SeverityPass {
			if count >= visibleRows {
				break
			}
			statusColor := colorYellow
			if pt.Severity == models.SeverityFail {
				statusColor = colorRed
			}
			statusStyled := lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render(strings.ToUpper(string(pt.Severity)))
			row := fmt.Sprintf("  [%s] Port conflict %d: %s", statusStyled, pt.Port, pt.Note)
			builder.WriteString(truncateString(row, width) + "\n")
			count++
		}
	}

	if count == 0 {
		builder.WriteString("\n  🟢 All diagnostics checks report PASS state!\n")
	}

	return builder.String()
}

func (m Model) renderDashboardQuickActions(width, _ int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("⚡ QUICK CONSOLE ACTIONS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	actions := []string{
		"🔍 Scan Compose Stacks",
		"🧹 Prune Unused Docker Resources",
		"🔄 Reconcile Server Daemon",
	}

	for idx, act := range actions {
		cursor := " "
		if m.dashboardFocusIndex == 2 && idx == m.selectedQuickActionIdx {
			cursor = ">"
		}
		labelStyle := lipgloss.NewStyle()
		if m.dashboardFocusIndex == 2 && idx == m.selectedQuickActionIdx {
			labelStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		}
		builder.WriteString(fmt.Sprintf("%s %s\n", cursor, labelStyle.Render(act)))
	}

	return builder.String()
}

// --- Tab 2: Containers & Docker Control View ---

func (m Model) viewContainers(height int) string {
	subTabHeight := 2
	innerContentHeight := height - subTabHeight

	// Sub tab header
	subTabs := []string{"Compose Stacks", "Docker Images", "Docker Volumes", "Docker Networks"}
	var renderedSub []string
	for i, name := range subTabs {
		isActive := false
		switch m.containersTabFocus {
		case 0, 1:
			isActive = (i == 0)
		case 2:
			isActive = (i == 1)
		case 3:
			isActive = (i == 2)
		case 4:
			isActive = (i == 3)
		}

		if isActive {
			renderedSub = append(renderedSub, lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Underline(true).Render(name))
		} else {
			renderedSub = append(renderedSub, lipgloss.NewStyle().Foreground(colorGray).Render(name))
		}
	}
	subTabsHeader := "📁 " + strings.Join(renderedSub, "  |  ") + "\n"

	var innerContent string
	switch m.containersTabFocus {
	case 0, 1: // Stacks
		leftWidth := 35
		rightWidth := m.width - leftWidth - 6
		leftStyle := styleUnfocusedPanel
		if m.containersTabFocus == 0 {
			leftStyle = styleFocusedPanel
		}
		leftContent := m.renderStacksList(leftWidth-4, innerContentHeight-2)
		leftBox := fitHeight(leftStyle.Width(leftWidth).Height(innerContentHeight-2).Render(leftContent), innerContentHeight, true)

		rightStyle := styleFocusedPanel
		if m.containersTabFocus == 0 {
			rightStyle = styleUnfocusedPanel
		}
		rightContent := m.renderServicesList(rightWidth-2, innerContentHeight-2)
		rightBox := fitHeight(rightStyle.Width(rightWidth).Height(innerContentHeight-2).Render(rightContent), innerContentHeight, true)
		innerContent = lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	case 2: // Images
		panelStyle := styleFocusedPanel
		content := m.renderDockerImagesList(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)

	case 3: // Volumes
		panelStyle := styleFocusedPanel
		content := m.renderDockerVolumesList(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)

	case 4: // Networks
		panelStyle := styleFocusedPanel
		content := m.renderDockerNetworksList(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	}

	return subTabsHeader + innerContent
}

func (m Model) renderDockerImagesList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📦 DOCKER IMAGES") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.dockerImages) == 0 {
		builder.WriteString("\n  No docker images found.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-35s %-12s %-15s %-15s\n", "REPOSITORY", "TAG", "IMAGE ID", "SIZE"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, img := range m.dockerImages {
		if idx >= visibleRows {
			break
		}
		cursor := " "
		if m.containersTabFocus == 2 && idx == m.selectedImageIdx {
			cursor = ">"
		}
		repo := truncateString(img.Repository, 33)
		tag := truncateString(img.Tag, 10)
		id := truncateString(strings.TrimPrefix(img.ID, "sha256:"), 12)
		sizeMB := float64(img.Size) / (1024 * 1024)

		row := fmt.Sprintf("%s %-35s %-12s %-15s %.1f MB\n", cursor, repo, tag, id, sizeMB)
		builder.WriteString(row)
	}
	return builder.String()
}

func (m Model) renderDockerVolumesList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🐳 DOCKER VOLUMES") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.dockerVolumes) == 0 {
		builder.WriteString("\n  No docker volumes found.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-45s %-15s %-20s\n", "VOLUME NAME", "DRIVER", "MOUNTPOINT"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, vol := range m.dockerVolumes {
		if idx >= visibleRows {
			break
		}
		cursor := " "
		if m.containersTabFocus == 3 && idx == m.selectedVolumeIdx {
			cursor = ">"
		}
		name := truncateString(vol.Name, 43)
		drv := truncateString(vol.Driver, 13)
		mp := truncateString(vol.Mountpoint, 28)

		row := fmt.Sprintf("%s %-45s %-15s %-20s\n", cursor, name, drv, mp)
		builder.WriteString(row)
	}
	return builder.String()
}

func (m Model) renderDockerNetworksList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🌐 DOCKER NETWORKS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.dockerNetworks) == 0 {
		builder.WriteString("\n  No docker networks found.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-15s %-30s %-15s %-15s\n", "NETWORK ID", "NAME", "DRIVER", "SCOPE"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, net := range m.dockerNetworks {
		if idx >= visibleRows {
			break
		}
		cursor := " "
		if m.containersTabFocus == 4 && idx == m.selectedNetworkIdx {
			cursor = ">"
		}
		id := truncateString(net.ID, 12)
		name := truncateString(net.Name, 28)
		drv := truncateString(net.Driver, 13)
		scp := truncateString(net.Scope, 13)

		row := fmt.Sprintf("%s %-15s %-30s %-15s %-15s\n", cursor, id, name, drv, scp)
		builder.WriteString(row)
	}
	return builder.String()
}

func (m Model) renderStacksList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📁 STACKS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.stacks) == 0 {
		builder.WriteString("\n  No stacks found.\n")
		return builder.String()
	}

	visibleRows := height - 2
	for idx, s := range m.stacks {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.containersTabFocus == 0 && idx == m.selectedStackIdx {
			cursor = ">"
		}

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
	builder.WriteString(styleHeaderLabel.Render("🐳 STACK SERVICES (Press 'e' on running container to Exec)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.stacks) == 0 {
		return builder.String()
	}
	selectedStackName := m.stacks[m.selectedStackIdx].Name

	if len(m.containers) == 0 {
		builder.WriteString(fmt.Sprintf("\n  No running services for stack: %s\n", selectedStackName))
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-25s %-25s %-20s\n", "NAME", "IMAGE", "STATE (STATUS)"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, c := range m.containers {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.containersTabFocus == 1 && idx == m.selectedContainerIdx {
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

// --- Tab 3: Container Logs View ---

func (m Model) viewLogs(height int) string {
	leftWidth := 30
	rightWidth := m.width - leftWidth - 6

	leftStyle := styleUnfocusedPanel
	if m.focusOnConfig { // in logs page we repurpose focus variable
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderLogsContainerList(leftWidth-4, height-2)
	leftBox := fitHeight(leftStyle.Width(leftWidth).Height(height-2).Render(leftContent), height, true)

	rightStyle := styleFocusedPanel
	if m.focusOnConfig {
		rightStyle = styleUnfocusedPanel
	}
	rightContent := m.renderLogsStream(rightWidth-2, height-2)
	rightBox := fitHeight(rightStyle.Width(rightWidth).Height(height-2).Render(rightContent), height, true)

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
		if m.focusOnConfig && idx == m.selectedLogContainerIdx {
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
	cname := "unknown"
	if len(m.logContainers) > 0 && m.selectedLogContainerIdx < len(m.logContainers) {
		c := m.logContainers[m.selectedLogContainerIdx]
		if len(c.Names) > 0 {
			cname = strings.TrimPrefix(c.Names[0], "/")
		}
	}

	// Logs title with filter level and follow badges
	followBadge := lipgloss.NewStyle().Foreground(colorRed).Render("[PAUSED]")
	if m.logFollow {
		followBadge = lipgloss.NewStyle().Foreground(colorGreen).Render("[FOLLOWING]")
	}
	searchBadge := ""
	if m.logSearchQuery != "" {
		searchBadge = lipgloss.NewStyle().Foreground(colorYellow).Render(" (Filtered: \"" + m.logSearchQuery + "\")")
	}

	builder.WriteString(styleHeaderLabel.Render(fmt.Sprintf("📄 LOG STREAM: %s  |  Level: %s  |  %s%s", cname, m.logLevelFilter, followBadge, searchBadge)) + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if m.showLogSearchPrompt {
		builder.WriteString(lipgloss.NewStyle().Foreground(colorYellow).Render(fmt.Sprintf("🔍 Log Search Filter: %s_", m.logSearchQuery)) + "\n")
		builder.WriteString(strings.Repeat("─", width) + "\n")
	}

	if m.logs == "" {
		builder.WriteString(fmt.Sprintf("\n  Loading logs for container: %s...\n", cname))
		return builder.String()
	}

	// Filter logs by level and search query
	lines := strings.Split(m.logs, "\n")
	var filteredLines []string
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Level Filter
		if m.logLevelFilter != "ALL" {
			upperLine := strings.ToUpper(line)
			if m.logLevelFilter == "INFO" && !strings.Contains(upperLine, "INFO") && !strings.Contains(upperLine, "INF") {
				continue
			}
			if m.logLevelFilter == "WARN" && !strings.Contains(upperLine, "WARN") && !strings.Contains(upperLine, "WRN") && !strings.Contains(upperLine, "WARNING") {
				continue
			}
			if m.logLevelFilter == "ERROR" && !strings.Contains(upperLine, "ERR") && !strings.Contains(upperLine, "ERROR") && !strings.Contains(upperLine, "FAIL") {
				continue
			}
		}

		// Search Query Filter
		if m.logSearchQuery != "" {
			if !strings.Contains(strings.ToLower(line), strings.ToLower(m.logSearchQuery)) {
				continue
			}
		}

		filteredLines = append(filteredLines, line)
	}

	totalLines := len(filteredLines)
	scrollAreaHeight := height - 4
	if m.showLogSearchPrompt {
		scrollAreaHeight = height - 5
	}
	if scrollAreaHeight < 1 {
		scrollAreaHeight = 1
	}

	m.logScrollHeight = scrollAreaHeight

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

	start := totalLines - m.logScrollHeight - m.logScrollOffset
	if start < 0 {
		start = 0
	}
	end := start + m.logScrollHeight
	if end > totalLines {
		end = totalLines
	}

	for i := start; i < end; i++ {
		builder.WriteString(truncateString(filteredLines[i], width) + "\n")
	}

	scrollHelp := fmt.Sprintf("─ [Line %d-%d of %d] (Arrow Up/Down: Scroll | /: Search | l: Level | f: Follow | o: Export) ─", start+1, end, totalLines)
	if totalLines == 0 {
		scrollHelp = "─ No matching logs found ─"
	}
	builder.WriteString("\n" + lipgloss.NewStyle().Foreground(colorGray).Render(scrollHelp))

	return builder.String()
}

// --- Tab 4: Config Editor View ---

func (m Model) viewEditor(height int) string {
	leftWidth := 35
	rightWidth := m.width - leftWidth - 6

	leftStyle := styleUnfocusedPanel
	if m.focusOnConfig {
		leftStyle = styleFocusedPanel
	}
	leftContent := m.renderConfigsList(leftWidth-4, height-2)
	leftBox := fitHeight(leftStyle.Width(leftWidth).Height(height-2).Render(leftContent), height, true)

	rightStyle := styleFocusedPanel
	if m.focusOnConfig {
		rightStyle = styleUnfocusedPanel
	}
	rightContent := m.renderSelectedConfigViewer(rightWidth-4, height-2)
	rightBox := fitHeight(rightStyle.Width(rightWidth).Height(height-2).Render(rightContent), height, true)

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
	visibleRows := height - 4
	if visibleRows < 1 {
		visibleRows = 1
	}

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

	scrollHelp := fmt.Sprintf("─ [Line %d-%d of %d] (Arrow Up/Down to scroll, 'e' to edit) ─", m.configScrollOffset+1, endIdx, totalLines)
	builder.WriteString("\n" + lipgloss.NewStyle().Foreground(colorGray).Render(scrollHelp))

	return builder.String()
}

// --- Tab 5: System Admin View ---

func (m Model) viewSystem(height int) string {
	subTabHeight := 2
	innerContentHeight := height - subTabHeight

	// Sub tab titles
	subTabs := []string{"systemd Services", "Partitions / Mounts", "Scheduled Tasks", "System Updates"}
	var renderedSub []string
	for i, name := range subTabs {
		if i == m.systemTabFocus {
			renderedSub = append(renderedSub, lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Underline(true).Render(name))
		} else {
			renderedSub = append(renderedSub, lipgloss.NewStyle().Foreground(colorGray).Render(name))
		}
	}
	subTabsHeader := "⚙️  " + strings.Join(renderedSub, "  |  ") + "\n"

	var innerContent string
	panelStyle := styleFocusedPanel

	switch m.systemTabFocus {
	case 0: // Services
		content := m.renderSystemServices(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	case 1: // Partitions
		content := m.renderSystemPartitions(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	case 2: // Cron
		content := m.renderSystemCron(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	case 3: // Updates
		content := m.renderSystemUpdates(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	}

	return subTabsHeader + innerContent
}

func (m Model) renderSystemServices(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("⚙️ SYSTEMD SERVICES (s: Start | x: Stop | t: Restart | v: Toggle Enable)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.systemServices) == 0 {
		builder.WriteString("\n  No services found or unable to fetch systemd unit status.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-30s %-12s %-12s %-12s %-30s\n", "SERVICE UNIT", "LOADED", "ACTIVE", "SUB-STATE", "DESCRIPTION"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, s := range m.systemServices {
		if idx >= visibleRows {
			break
		}
		cursor := " "
		if m.systemTabFocus == 0 && idx == m.selectedServiceIdx {
			cursor = ">"
		}

		var activeStyle lipgloss.Style
		switch s.ActiveState {
		case "active":
			activeStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		case "inactive":
			activeStyle = lipgloss.NewStyle().Foreground(colorGray)
		default:
			activeStyle = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
		}

		nameT := truncateString(s.Name, 28)
		loadT := truncateString(s.LoadState, 10)
		actT := activeStyle.Render(strings.ToUpper(s.ActiveState))
		subT := truncateString(s.SubState, 10)
		descT := truncateString(s.Description, width-72)

		row := fmt.Sprintf("%s %-30s %-12s %-12s %-12s %-30s\n", cursor, nameT, loadT, actT, subT, descT)
		builder.WriteString(row)
	}

	return builder.String()
}

func (m Model) renderSystemPartitions(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("💿 DISK PARTITIONS & SAMBA NETWORK SHARES") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.diskPartitions) == 0 {
		builder.WriteString("\n  Loading partitions telemetry...\n")
		return builder.String()
	}

	// Partitions List
	builder.WriteString(lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("  Disk Partitions:") + "\n")
	builder.WriteString(fmt.Sprintf("  %-25s %-20s %-10s %-12s %-12s %-12s\n", "DEVICE", "MOUNTPOINT", "FS TYPE", "TOTAL (GB)", "USED (GB)", "USED %"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	for _, p := range m.diskPartitions {
		dev := truncateString(p.Device, 23)
		mp := truncateString(p.Mountpoint, 18)
		fs := truncateString(p.FSType, 8)
		row := fmt.Sprintf("  %-25s %-20s %-10s %-12.1f %-12.1f %.1f%%\n", dev, mp, fs, p.Total, p.Used, p.UsedPercent)
		builder.WriteString(row)
	}

	builder.WriteString("\n")

	// Samba Shares List
	builder.WriteString(lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("  Samba Network Shares:") + "\n")
	builder.WriteString(fmt.Sprintf("  %-25s %-30s %-12s %-12s\n", "SHARE NAME", "PATH", "READ ONLY", "GUEST OK"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.sambaShares) == 0 {
		builder.WriteString("  No Samba shares found.\n")
	} else {
		for _, s := range m.sambaShares {
			name := truncateString(s.Name, 23)
			path := truncateString(s.Path, 28)
			ro := "no"
			if s.ReadOnly {
				ro = "yes"
			}
			goOk := "no"
			if s.GuestOk {
				goOk = "yes"
			}

			row := fmt.Sprintf("  %-25s %-30s %-12s %-12s\n", name, path, ro, goOk)
			builder.WriteString(row)
		}
	}

	return builder.String()
}

func (m Model) renderSystemCron(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📅 SCHEDULED TASKS (CRONTAB & SYSTEMD TIMERS)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.cronJobs) == 0 {
		builder.WriteString("\n  No scheduled tasks found.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-30s %-15s %-10s %-40s\n", "SCHEDULE / TIMER", "SOURCE", "USER", "COMMAND"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, job := range m.cronJobs {
		if idx >= visibleRows {
			break
		}
		cursor := " "
		if m.systemTabFocus == 2 && idx == m.selectedCronIdx {
			cursor = ">"
		}
		sched := truncateString(job.Schedule, 28)
		src := truncateString(job.Source, 13)
		user := truncateString(job.User, 8)
		cmd := truncateString(job.Command, width-60)

		row := fmt.Sprintf("%s %-30s %-15s %-10s %-40s\n", cursor, sched, src, user, cmd)
		builder.WriteString(row)
	}

	return builder.String()
}

func (m Model) renderSystemUpdates(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("📥 PENDING SYSTEM UPDATES (Press 'u' to upgrade)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	if m.systemUpdates == nil {
		builder.WriteString("  Checking for updates...\n")
		return builder.String()
	}

	if !m.systemUpdates.HasUpdates {
		builder.WriteString("  🟢 System is completely up to date! No updates available.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  🔔 %d packages have pending updates available:\n\n", m.systemUpdates.Count))

	visibleRows := height - 6
	for idx, pkg := range m.systemUpdates.UpdatesList {
		if idx >= visibleRows {
			builder.WriteString("  ... and more packages ...\n")
			break
		}
		builder.WriteString(fmt.Sprintf("  • %s\n", pkg))
	}

	return builder.String()
}

// --- Tab 6: Terminal & Shell & Integrations View ---

func (m Model) viewTerminal(height int) string {
	subTabHeight := 2
	innerContentHeight := height - subTabHeight

	// Sub tab titles
	subTabs := []string{"Shell Launcher", "Saved SSH Targets", "Ollama AI Queue", "Plugins Catalog"}
	var renderedSub []string
	for i, name := range subTabs {
		isActive := false
		switch m.terminalTabFocus {
		case 0:
			isActive = (i == 0)
		case 1:
			isActive = (i == 1)
		case 2, 3:
			isActive = (i == 2)
		case 4:
			isActive = (i == 3)
		}

		if isActive {
			renderedSub = append(renderedSub, lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Underline(true).Render(name))
		} else {
			renderedSub = append(renderedSub, lipgloss.NewStyle().Foreground(colorGray).Render(name))
		}
	}
	subTabsHeader := "💻 " + strings.Join(renderedSub, "  |  ") + "\n"

	var innerContent string
	panelStyle := styleFocusedPanel

	switch m.terminalTabFocus {
	case 0: // Shell Launcher
		content := m.renderLocalShellLauncher(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	case 1: // SSH profiles
		content := m.renderSSHProfiles(m.width-6, innerContentHeight-2)
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	case 2, 3: // AI Queue & Models
		leftWidth := m.width * 7 / 10
		if leftWidth < 40 {
			leftWidth = 40
		}
		rightWidth := m.width - leftWidth - 6

		leftStyle := styleUnfocusedPanel
		if m.terminalTabFocus == 2 {
			leftStyle = styleFocusedPanel
		}
		leftContent := m.renderAIQueue(leftWidth-4, innerContentHeight-2)
		leftBox := fitHeight(leftStyle.Width(leftWidth).Height(innerContentHeight-2).Render(leftContent), innerContentHeight, true)

		rightStyle := styleFocusedPanel
		if m.terminalTabFocus == 2 {
			rightStyle = styleUnfocusedPanel
		}
		rightContent := m.renderAIModels(rightWidth-2, innerContentHeight-2)
		rightBox := fitHeight(rightStyle.Width(rightWidth).Height(innerContentHeight-2).Render(rightContent), innerContentHeight, true)
		innerContent = lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	case 4: // Plugins
		width := m.width - 6
		var content string
		if m.showCatalog {
			content = m.renderPluginCatalog(width-4, innerContentHeight-2)
		} else {
			content = m.renderPluginsList(width-4, innerContentHeight-2)
		}
		innerContent = fitHeight(panelStyle.Width(m.width-2).Height(innerContentHeight-2).Render(content), innerContentHeight, true)
	}

	return subTabsHeader + innerContent
}

func (m Model) renderLocalShellLauncher(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("💻 HOST TERMINAL LAUNCHER") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	builder.WriteString("  This launcher lets you open a full, interactive terminal session on the host.\n")
	builder.WriteString("  When launched, the TUI suspends, giving you direct command prompt access.\n")
	builder.WriteString("  Upon exiting the shell session (type 'exit' or Ctrl+D), the TUI resumes automatically.\n\n")

	builder.WriteString(lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("  👉 Press ENTER to open local shell session") + "\n\n")

	builder.WriteString(lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render("  Active Host Environment Variables:") + "\n")
	builder.WriteString(strings.Repeat("─", 40) + "\n")
	for _, env := range m.envVars {
		builder.WriteString(fmt.Sprintf("  • %s\n", env))
	}

	return builder.String()
}

func (m Model) renderSSHProfiles(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🔑 SAVED SSH REMOTE TARGETS") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n\n")

	builder.WriteString("  Select a remote server profile below to launch an interactive SSH session.\n")
	builder.WriteString("  TUI will suspend and hand terminal controls over to SSH until you disconnect.\n\n")

	builder.WriteString(fmt.Sprintf("  %-30s %-30s\n", "CONNECTION NAME", "SSH LAUNCH COMMAND"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	for idx, prof := range m.sshProfiles {
		cursor := " "
		if m.terminalTabFocus == 1 && idx == m.selectedSSHIdx {
			cursor = ">"
		}
		labelStyle := lipgloss.NewStyle()
		if m.terminalTabFocus == 1 && idx == m.selectedSSHIdx {
			labelStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		}

		builder.WriteString(fmt.Sprintf("%s %-30s %-30s\n", cursor, labelStyle.Render(prof), "ssh "+prof))
	}

	builder.WriteString("\n  (Press Up/Down to navigate, Enter to launch session)")
	return builder.String()
}

func (m Model) renderAIQueue(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🤖 AI JOB QUEUE") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.queue) == 0 {
		builder.WriteString("\n  No active or historical AI tasks in queue.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-12s %-16s %-10s %-20s %-12s\n", "JOB ID", "TYPE", "PRIORITY", "DETAILS", "STATUS"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, job := range m.queue {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.terminalTabFocus == 2 && idx == m.selectedJobIdx {
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
		if m.terminalTabFocus == 3 && idx == m.selectedModelIdx {
			cursor = ">"
		}

		modelT := truncateString(model, width-3)
		builder.WriteString(fmt.Sprintf("%s %s\n", cursor, modelT))
	}

	return builder.String()
}

func (m Model) renderPluginsList(width, height int) string {
	var builder strings.Builder
	builder.WriteString(styleHeaderLabel.Render("🧩 LOADED PLUGINS (Press 'c' to view Catalog | 'g' to Toggle Enable)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if m.plugins == nil {
		builder.WriteString("\n  Loading loaded plugins list...\n")
		return builder.String()
	}

	var list []models.Plugin
	list = append(list, m.plugins.Routes...)
	list = append(list, m.plugins.Stacks...)
	list = append(list, m.plugins.Middleware...)

	if len(list) == 0 {
		builder.WriteString("\n  No plugins loaded.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-12s %-20s %-10s %-20s %-12s %-20s\n", "KIND", "NAME", "VERSION", "AUTHOR", "STATUS", "WARNINGS"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, p := range list {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.terminalTabFocus == 4 && idx == m.selectedPluginIdx {
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
	builder.WriteString(styleHeaderLabel.Render("🧩 PLUGINS CATALOG (Press 'c' to view Installed | 'i' to Install)") + "\n")
	builder.WriteString(strings.Repeat("─", width) + "\n")

	if len(m.catalog) == 0 {
		builder.WriteString("\n  No catalog items found.\n")
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("  %-12s %-20s %-10s %-20s %-15s %-30s\n", "KIND", "NAME", "VERSION", "AUTHOR", "INSTALL STATUS", "DESCRIPTION"))
	builder.WriteString(strings.Repeat("─", width) + "\n")

	visibleRows := height - 4
	for idx, item := range m.catalog {
		if idx >= visibleRows {
			break
		}

		cursor := " "
		if m.terminalTabFocus == 4 && idx == m.selectedPluginIdx {
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

// --- Metrics Bar (Bottom) ---

func (m Model) renderMetricsBar() string {
	if m.detailedStats == nil {
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

	cpuStr := progBar(m.detailedStats.CPUUsage)
	ramStr := progBar(m.detailedStats.MemoryUsage)
	diskStr := progBar(m.detailedStats.DiskUsage)

	return fmt.Sprintf("🔥 CPU: %s  |  🧠 RAM: %s  |  💿 Disk: %s", cpuStr, ramStr, diskStr)
}

// --- Footer Hotkeys ---

func (m Model) renderFooter() string {
	var keys []string
	switch m.activeTab {
	case TabDashboard:
		actionKey := "Enter: Trigger Action"
		if m.dashboardFocusIndex == 1 {
			if m.preflightExpanded {
				actionKey = "Enter: Collapse Alerts"
			} else {
				actionKey = "Enter: Expand Alerts"
			}
		}
		keys = []string{"1-6: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate", actionKey, "q: Quit"}
	case TabContainers:
		keys = []string{"1-6: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate", "s/x/t: Start/Stop/Restart", "u/d: Stack Up/Down", "p: Pull images", "e: Exec in container", "q: Quit"}
	case TabLogs:
		keys = []string{"1-6: Switch Tabs", "Tab: Select Panel", "Arrows: Scroll/Navigate", "/: Search", "l: Level Filter", "f: Follow Mode", "o: Export logs", "q: Quit"}
	case TabEditor:
		keys = []string{"1-6: Switch Tabs", "Tab: Swap Panels", "Arrows: Scroll/Navigate", "e: Edit File (nano)", "q: Quit"}
	case TabSystem:
		keys = []string{"1-6: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate", "s/x/t: Start/Stop/Restart Service", "v: Toggle Enable Service", "u: Apply pending upgrades", "q: Quit"}
	case TabTerminal:
		keys = []string{"1-6: Switch Tabs", "Tab: Swap Panels", "Arrows: Navigate", "Enter: Launch Connection/Shell", "c: Catalog Catalog/Installed", "i/g: Plugin Install/Toggle Enable", "q: Quit"}
	}
	return styleFooter.Width(m.width).Render("🔑 Keys: " + strings.Join(keys, " | "))
}

func truncateString(s string, limit int) string {
	if limit <= 3 {
		return "..."
	}
	if lipgloss.Width(s) <= limit {
		return s
	}

	var builder strings.Builder
	visibleWidth := 0
	inEscape := false

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '\x1b' {
			inEscape = true
			builder.WriteRune(r)
			continue
		}
		if inEscape {
			builder.WriteRune(r)
			if r == 'm' {
				inEscape = false
			}
			continue
		}

		w := lipgloss.Width(string(r))
		if visibleWidth+w > limit-3 {
			builder.WriteString("...")
			builder.WriteString("\x1b[0m")
			break
		}
		builder.WriteRune(r)
		visibleWidth += w
	}
	return builder.String()
}

func fitHeight(rendered string, targetHeight int, hasBorder bool) string {
	lines := strings.Split(rendered, "\n")
	if len(lines) <= targetHeight {
		return rendered
	}
	if hasBorder {
		result := make([]string, targetHeight)
		copy(result, lines[:targetHeight-1])
		result[targetHeight-1] = lines[len(lines)-1]
		return strings.Join(result, "\n")
	}
	return strings.Join(lines[:targetHeight], "\n")
}
