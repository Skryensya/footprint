package logs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/footprint-tools/footprint-cli/internal/ui/splitpanel"
	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model
func (m logsModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate dimensions
	headerHeight := 3
	footerHeight := 2
	mainHeight := m.height - headerHeight - footerHeight
	if mainHeight < 1 {
		mainHeight = 1
	}

	// Create layout with drawer support
	cfg := splitpanel.Config{
		SidebarWidthPercent: 0.18,
		SidebarMinWidth:     16,
		SidebarMaxWidth:     22,
		HasDrawer:           true,
		DrawerWidthPercent:  0.35,
	}
	layout := splitpanel.NewLayout(m.width, cfg, m.colors)
	layout.SetFocus(false)
	layout.SetDrawerOpen(m.drawerOpen)

	// Build panels
	statsPanel := m.buildStatsPanel(layout, mainHeight)
	logsPanel := m.buildLogsPanel(layout, mainHeight)

	// Render components
	header := m.renderHeader()

	var main string
	if m.drawerOpen {
		drawerPanel := m.buildDrawerPanel(layout, mainHeight)
		main = layout.RenderWithDrawer(statsPanel, logsPanel, &drawerPanel, mainHeight)
	} else {
		main = layout.Render(statsPanel, logsPanel, mainHeight)
	}

	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func (m logsModel) renderHeader() string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	warnColor := lipgloss.Color(colors.Warning)
	successColor := lipgloss.Color(colors.Success)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(infoColor)
	mutedStyle := lipgloss.NewStyle().Foreground(mutedColor)
	warnStyle := lipgloss.NewStyle().Foreground(warnColor)
	successStyle := lipgloss.NewStyle().Foreground(successColor)

	// Title
	title := titleStyle.Render("fp logs")

	// Session duration
	duration := m.sessionDuration()
	hours := int(duration.Hours())
	mins := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)

	// Status indicator
	status := ""
	if m.paused {
		status = warnStyle.Render(" [PAUSED]")
	}
	if m.autoScroll {
		status += successStyle.Render(" [AUTO]")
	}

	// Filter indicators
	filterStr := ""
	if m.filterLevel != "" {
		filterStr = mutedStyle.Render(" | Level: ") + m.levelStyle(m.filterLevel).Render(m.filterLevel)
	}
	if m.filterQuery != "" {
		filterStr += mutedStyle.Render(" | Search: ") + mutedStyle.Render(m.filterQuery)
	}

	headerContent := title + mutedStyle.Render(" | ") +
		mutedStyle.Render("Session: ") + timeStr +
		status + filterStr

	headerStyle := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1)

	return headerStyle.Render(headerContent)
}

func (m *logsModel) buildStatsPanel(layout *splitpanel.Layout, height int) splitpanel.Panel {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	successColor := lipgloss.Color(colors.Success)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(successColor)
	labelStyle := lipgloss.NewStyle().Foreground(mutedColor)
	valueStyle := lipgloss.NewStyle().Foreground(infoColor)

	var lines []string

	// Stats header
	lines = append(lines, headerStyle.Render("STATS"))
	lines = append(lines, "")

	// Total lines
	lines = append(lines, labelStyle.Render("Total: ")+valueStyle.Render(fmt.Sprintf("%d", len(m.lines))))

	// Filtered count
	filtered := m.filteredLines()
	if m.filterLevel != "" || m.filterQuery != "" {
		lines = append(lines, labelStyle.Render("Shown: ")+valueStyle.Render(fmt.Sprintf("%d", len(filtered))))
	}

	lines = append(lines, "")

	// By level
	if len(m.byLevel) > 0 {
		lines = append(lines, headerStyle.Render("BY LEVEL"))
		lines = append(lines, "")

		// Sort levels for consistent display
		levelOrder := []string{"ERROR", "WARN", "INFO", "DEBUG"}
		for _, level := range levelOrder {
			if count, ok := m.byLevel[level]; ok {
				style := m.levelStyle(level)
				indicator := "  "
				if m.filterLevel == level {
					indicator = "> "
				}
				lines = append(lines, indicator+style.Render(fmt.Sprintf("%-6s", level))+labelStyle.Render(fmt.Sprintf(" %d", count)))
			}
		}

		// Any other levels
		var otherLevels []string
		for level := range m.byLevel {
			found := false
			for _, l := range levelOrder {
				if level == l {
					found = true
					break
				}
			}
			if !found {
				otherLevels = append(otherLevels, level)
			}
		}
		sort.Strings(otherLevels)
		for _, level := range otherLevels {
			count := m.byLevel[level]
			lines = append(lines, labelStyle.Render(fmt.Sprintf("  %-6s %d", level, count)))
		}
	}

	lines = append(lines, "")

	// Key hints
	lines = append(lines, headerStyle.Render("FILTERS"))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("e ERROR"))
	lines = append(lines, labelStyle.Render("w WARN"))
	lines = append(lines, labelStyle.Render("i INFO"))
	lines = append(lines, labelStyle.Render("d DEBUG"))
	lines = append(lines, labelStyle.Render("c clear"))

	return splitpanel.Panel{
		Lines:      lines,
		ScrollPos:  0,
		TotalItems: len(lines),
	}
}

func (m *logsModel) buildLogsPanel(layout *splitpanel.Layout, height int) splitpanel.Panel {
	colors := m.colors
	mutedColor := lipgloss.Color(colors.Muted)

	filtered := m.filteredLines()
	visibleHeight := height - 2

	// Adjust scroll to keep cursor visible
	scrollOffset := m.scrollPos
	if m.cursor < scrollOffset {
		scrollOffset = m.cursor
	}
	if m.cursor >= scrollOffset+visibleHeight {
		scrollOffset = m.cursor - visibleHeight + 1
	}
	m.scrollPos = scrollOffset

	var lines []string
	width := layout.MainContentWidth()

	if len(filtered) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(mutedColor).Italic(true)
		if m.filterQuery != "" || m.filterLevel != "" {
			lines = append(lines, emptyStyle.Render("No matching log lines"))
		} else {
			lines = append(lines, emptyStyle.Render("No log lines yet..."))
		}
	} else {
		for i := scrollOffset; i < len(filtered) && len(lines) < visibleHeight; i++ {
			logLine := filtered[i]
			line := m.formatLogLine(logLine, width, i == m.cursor)
			lines = append(lines, line)
		}
	}

	return splitpanel.Panel{
		Lines:      lines,
		ScrollPos:  scrollOffset,
		TotalItems: len(filtered),
	}
}

func (m logsModel) formatLogLine(logLine LogLine, width int, selected bool) string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)

	// Format: [timestamp] LEVEL message
	var line string

	if logLine.Timestamp != "" {
		// Shorten timestamp (just time part)
		ts := logLine.Timestamp
		if len(ts) > 11 {
			ts = ts[11:] // Skip date, keep time
		}
		line = "[" + ts + "] "
	}

	if logLine.Level != "" {
		line += logLine.Level + ": "
	}

	if logLine.Message != "" {
		line += logLine.Message
	} else if logLine.Raw != "" && logLine.Timestamp == "" {
		line = logLine.Raw
	}

	// Truncate if too long
	if len(line) > width-4 {
		line = line[:width-7] + "..."
	}

	prefix := "  "
	if selected {
		prefix = "> "
	}

	if selected {
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(infoColor)
		return style.Render(prefix + line)
	}

	// Color by level
	if logLine.Level != "" {
		levelStyle := m.levelStyle(logLine.Level)
		tsStyle := lipgloss.NewStyle().Foreground(mutedColor)

		if logLine.Timestamp != "" {
			ts := logLine.Timestamp
			if len(ts) > 11 {
				ts = ts[11:]
			}
			return prefix + tsStyle.Render("["+ts+"] ") + levelStyle.Render(logLine.Level+": "+logLine.Message)
		}
		return prefix + levelStyle.Render(line)
	}

	return prefix + lipgloss.NewStyle().Foreground(mutedColor).Render(line)
}

func (m *logsModel) buildDrawerPanel(layout *splitpanel.Layout, height int) splitpanel.Panel {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	successColor := lipgloss.Color(colors.Success)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(successColor)
	labelStyle := lipgloss.NewStyle().Foreground(mutedColor)
	valueStyle := lipgloss.NewStyle().Foreground(infoColor)

	var lines []string
	width := layout.DrawerContentWidth()

	if m.drawerDetail == nil {
		lines = append(lines, labelStyle.Render("No line selected"))
	} else {
		logLine := m.drawerDetail

		// Header
		lines = append(lines, headerStyle.Render("LOG DETAIL"))
		lines = append(lines, "")

		// Level
		if logLine.Level != "" {
			levelStyle := m.levelStyle(logLine.Level)
			lines = append(lines, labelStyle.Render("Level: ")+levelStyle.Render(logLine.Level))
			lines = append(lines, "")
		}

		// Timestamp
		if logLine.Timestamp != "" {
			lines = append(lines, labelStyle.Render("Timestamp:"))
			lines = append(lines, valueStyle.Render("  "+logLine.Timestamp))
			lines = append(lines, "")
		}

		// Message
		if logLine.Message != "" {
			lines = append(lines, headerStyle.Render("MESSAGE"))
			lines = append(lines, "")
			// Wrap message
			wrapped := wrapText(logLine.Message, width-4)
			for _, l := range strings.Split(wrapped, "\n") {
				lines = append(lines, valueStyle.Render("  "+l))
			}
			lines = append(lines, "")
		}

		// Raw line
		lines = append(lines, headerStyle.Render("RAW"))
		lines = append(lines, "")
		wrapped := wrapText(logLine.Raw, width-4)
		for _, l := range strings.Split(wrapped, "\n") {
			lines = append(lines, labelStyle.Render("  "+l))
		}
	}

	return splitpanel.Panel{
		Lines:      lines,
		ScrollPos:  0,
		TotalItems: len(lines),
	}
}

func (m logsModel) renderFooter() string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	borderColor := lipgloss.Color(colors.Border)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(infoColor).
		Padding(0, 1)

	sepStyle := lipgloss.NewStyle().Foreground(borderColor)
	labelStyle := lipgloss.NewStyle().Foreground(mutedColor)

	sep := sepStyle.Render(" | ")

	var footer string
	if m.drawerOpen {
		footer = keyStyle.Render("Esc") + labelStyle.Render(" close") + sep +
			keyStyle.Render("jk") + labelStyle.Render(" navigate")
	} else {
		footer = keyStyle.Render("q") + labelStyle.Render(" quit") + sep +
			keyStyle.Render("p") + labelStyle.Render(" pause") + sep +
			keyStyle.Render("a") + labelStyle.Render(" auto") + sep +
			keyStyle.Render("jk") + labelStyle.Render(" nav") + sep +
			keyStyle.Render("Enter") + labelStyle.Render(" detail") + sep +
			labelStyle.Render("type to search")
	}

	footerStyle := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1)

	return footerStyle.Render(footer)
}

func (m logsModel) levelStyle(level string) lipgloss.Style {
	colors := m.colors
	switch level {
	case "ERROR":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Error))
	case "WARN":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Warning))
	case "INFO":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Info))
	case "DEBUG":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Muted))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Muted))
	}
}

func wrapText(text string, width int) string {
	if width <= 0 || len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		if i > 0 {
			if lineLen+1+len(word) > width {
				result.WriteString("\n")
				lineLen = 0
			} else {
				result.WriteString(" ")
				lineLen++
			}
		}
		result.WriteString(word)
		lineLen += len(word)
	}

	return result.String()
}
