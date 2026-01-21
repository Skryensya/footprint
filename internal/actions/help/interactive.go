package help

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/help"
	"github.com/Skryensya/footprint/internal/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

//
// Public API
//

func Browser(args []string, flags *dispatchers.ParsedFlags) error {
	return browser(args, flags, DefaultDeps())
}

//
// Entrypoint
//

func browser(_ []string, _ *dispatchers.ParsedFlags, deps Deps) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return errors.New("interactive-help requires an interactive terminal")
	}

	root := deps.BuildTree()
	topics := deps.AllTopics()
	items := buildSidebarItems(root, topics)

	// Find first selectable item (skip category headers)
	cursor := 0
	for i, item := range items {
		if !item.IsCategory {
			cursor = i
			break
		}
	}

	m := model{
		allItems:      items,
		items:         items,
		cursor:        cursor,
		colors:        style.GetColors(),
		focusSidebar:  true,
		searchMode:    false,
		searchQuery:   "",
		totalCommands: countSelectableItems(items),
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}

//
// Sidebar Item
//

type sidebarItem struct {
	Name        string
	DisplayName string
	IsCategory  bool
	IsTopic     bool
	Node        *dispatchers.DispatchNode
	Topic       *help.Topic
}

//
// Model
//

type model struct {
	allItems       []sidebarItem // Original unfiltered items
	items          []sidebarItem // Filtered items to display
	cursor         int
	sidebarScroll  int
	contentScroll  int
	width          int
	height         int
	sidebarWidth   int  // Calculated sidebar width for mouse detection
	headerHeight   int  // Header height for mouse detection
	cancelled      bool
	colors         style.ColorConfig
	focusSidebar   bool   // true = sidebar focused, false = content focused
	searchMode     bool   // true = search input active
	searchQuery    string // Current search query
	totalCommands  int    // Total selectable items count
}

func countSelectableItems(items []sidebarItem) int {
	count := 0
	for _, item := range items {
		if !item.IsCategory {
			count++
		}
	}
	return count
}

//
// Build Sidebar
//

func buildSidebarItems(root *dispatchers.DispatchNode, topics []*help.Topic) []sidebarItem {
	var items []sidebarItem

	// Collect leaf commands
	var leaves []*dispatchers.DispatchNode
	for _, child := range root.Children {
		collectLeafCommands(child, &leaves)
	}

	// Group by category
	grouped := make(map[dispatchers.CommandCategory][]*dispatchers.DispatchNode)
	for _, cmd := range leaves {
		grouped[cmd.Category] = append(grouped[cmd.Category], cmd)
	}

	// Sort commands within each category
	for cat := range grouped {
		cmds := grouped[cat]
		sort.Slice(cmds, func(i, j int) bool {
			nameI := strings.Join(cmds[i].Path[1:], " ")
			nameJ := strings.Join(cmds[j].Path[1:], " ")
			return nameI < nameJ
		})
	}

	// Add categories in order
	for _, cat := range dispatchers.CategoryOrder() {
		cmds := grouped[cat]
		if len(cmds) == 0 {
			continue
		}

		// Category header (non-selectable)
		items = append(items, sidebarItem{
			Name:        cat.String(),
			DisplayName: strings.ToUpper(cat.String()),
			IsCategory:  true,
		})

		// Commands in this category
		for _, cmd := range cmds {
			displayName := strings.Join(cmd.Path[1:], " ")
			items = append(items, sidebarItem{
				Name:        displayName,
				DisplayName: displayName,
				IsCategory:  false,
				IsTopic:     false,
				Node:        cmd,
			})
		}
	}

	// Add conceptual guides section
	if len(topics) > 0 {
		items = append(items, sidebarItem{
			Name:        "conceptual",
			DisplayName: "CONCEPTUAL GUIDES",
			IsCategory:  true,
		})

		for _, topic := range topics {
			items = append(items, sidebarItem{
				Name:        topic.Name,
				DisplayName: topic.Name,
				IsCategory:  false,
				IsTopic:     true,
				Topic:       topic,
			})
		}
	}

	return items
}

func collectLeafCommands(node *dispatchers.DispatchNode, out *[]*dispatchers.DispatchNode) {
	if node.Action != nil {
		*out = append(*out, node)
		return
	}

	for _, child := range node.Children {
		collectLeafCommands(child, out)
	}
}

//
// Filter items based on search query
//

func (m *model) filterItems() {
	if m.searchQuery == "" {
		m.items = m.allItems
		return
	}

	query := strings.ToLower(m.searchQuery)
	var filtered []sidebarItem
	var currentCategory *sidebarItem
	hasItemsInCategory := false

	for i := range m.allItems {
		item := m.allItems[i]
		if item.IsCategory {
			// Store category, will add if it has matching items
			if currentCategory != nil && hasItemsInCategory {
				// Previous category had items, it's already added
			}
			currentCategory = &m.allItems[i]
			hasItemsInCategory = false
			continue
		}

		// Check if item matches search
		if strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.DisplayName), query) ||
			(item.Node != nil && strings.Contains(strings.ToLower(item.Node.Summary), query)) ||
			(item.Topic != nil && strings.Contains(strings.ToLower(item.Topic.Summary), query)) {

			// Add category header if this is first match in category
			if currentCategory != nil && !hasItemsInCategory {
				filtered = append(filtered, *currentCategory)
				hasItemsInCategory = true
			}
			filtered = append(filtered, item)
		}
	}

	m.items = filtered

	// Reset cursor to first selectable item
	m.cursor = 0
	for i, item := range m.items {
		if !item.IsCategory {
			m.cursor = i
			break
		}
	}
	m.sidebarScroll = 0
}

//
// Bubble Tea lifecycle
//

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Recalculate sidebar width for mouse detection
		m.sidebarWidth = m.width / 4
		if m.sidebarWidth < 24 {
			m.sidebarWidth = 24
		}
		if m.sidebarWidth > 36 {
			m.sidebarWidth = 36
		}
		m.headerHeight = 3
		return m, nil

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		// Handle search mode input
		if m.searchMode {
			switch msg.Type {
			case tea.KeyEsc:
				m.searchMode = false
				m.searchQuery = ""
				m.filterItems()
				return m, nil
			case tea.KeyEnter:
				m.searchMode = false
				return m, nil
			case tea.KeyBackspace:
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.filterItems()
				}
				return m, nil
			case tea.KeyRunes:
				m.searchQuery += string(msg.Runes)
				m.filterItems()
				return m, nil
			}
			return m, nil
		}

		switch msg.Type {

		case tea.KeyCtrlC:
			m.cancelled = true
			return m, tea.Quit

		case tea.KeyEsc:
			if m.searchQuery != "" {
				// Clear search first
				m.searchQuery = ""
				m.filterItems()
				return m, nil
			}
			m.cancelled = true
			return m, tea.Quit

		case tea.KeyTab:
			m.focusSidebar = !m.focusSidebar
			return m, nil

		case tea.KeyUp:
			if m.focusSidebar {
				m.moveCursor(-1)
				m.contentScroll = 0
			} else {
				m.contentScroll--
				if m.contentScroll < 0 {
					m.contentScroll = 0
				}
			}

		case tea.KeyDown:
			if m.focusSidebar {
				m.moveCursor(1)
				m.contentScroll = 0
			} else {
				m.contentScroll++
			}

		case tea.KeyPgUp:
			if m.focusSidebar {
				for i := 0; i < 5; i++ {
					m.moveCursor(-1)
				}
				m.contentScroll = 0
			} else {
				m.contentScroll -= 5
				if m.contentScroll < 0 {
					m.contentScroll = 0
				}
			}

		case tea.KeyPgDown:
			if m.focusSidebar {
				for i := 0; i < 5; i++ {
					m.moveCursor(1)
				}
				m.contentScroll = 0
			} else {
				m.contentScroll += 5
			}

		case tea.KeyHome:
			if m.focusSidebar {
				m.jumpToFirst()
				m.contentScroll = 0
			} else {
				m.contentScroll = 0
			}

		case tea.KeyEnd:
			if m.focusSidebar {
				m.jumpToLast()
				m.contentScroll = 0
			}

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				m.cancelled = true
				return m, tea.Quit
			case "/":
				m.searchMode = true
				return m, nil
			case "j":
				if m.focusSidebar {
					m.moveCursor(1)
					m.contentScroll = 0
				} else {
					m.contentScroll++
				}
			case "k":
				if m.focusSidebar {
					m.moveCursor(-1)
					m.contentScroll = 0
				} else {
					m.contentScroll--
					if m.contentScroll < 0 {
						m.contentScroll = 0
					}
				}
			case "g":
				if m.focusSidebar {
					m.jumpToFirst()
					m.contentScroll = 0
				} else {
					m.contentScroll = 0
				}
			case "G":
				if m.focusSidebar {
					m.jumpToLast()
					m.contentScroll = 0
				}
			case "u":
				m.contentScroll -= 5
				if m.contentScroll < 0 {
					m.contentScroll = 0
				}
			case "d":
				m.contentScroll += 5
			case "h":
				m.focusSidebar = true
			case "l":
				m.focusSidebar = false
			}
		}
	}

	return m, nil
}

func (m *model) moveCursor(delta int) {
	if len(m.items) == 0 {
		return
	}

	newCursor := m.cursor + delta

	// Stop at boundaries (no wrap around)
	if newCursor < 0 {
		newCursor = 0
	} else if newCursor >= len(m.items) {
		newCursor = len(m.items) - 1
	}

	// Skip category headers
	iterations := 0
	for m.items[newCursor].IsCategory && iterations < len(m.items) {
		newCursor += delta
		if delta == 0 {
			delta = 1
		}
		// Stop at boundaries when skipping categories
		if newCursor < 0 {
			// Can't go further up, stay at current position
			return
		} else if newCursor >= len(m.items) {
			// Can't go further down, stay at current position
			return
		}
		iterations++
	}

	m.cursor = newCursor
}

func (m *model) jumpToFirst() {
	for i, item := range m.items {
		if !item.IsCategory {
			m.cursor = i
			return
		}
	}
}

func (m *model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Calculate regions
	footerHeight := 2
	mainHeight := m.height - m.headerHeight - footerHeight

	// Check if click is in the main area (between header and footer)
	inMainArea := msg.Y >= m.headerHeight && msg.Y < m.height-footerHeight

	switch msg.Type {
	case tea.MouseLeft:
		if !inMainArea {
			return *m, nil
		}

		// Determine if click is in sidebar or content
		if msg.X < m.sidebarWidth {
			// Click in sidebar - focus it and optionally select item
			m.focusSidebar = true

			// Calculate which item was clicked
			clickedLine := msg.Y - m.headerHeight
			clickedItem := m.sidebarScroll + clickedLine

			if clickedItem >= 0 && clickedItem < len(m.items) {
				// Skip category headers - find nearest selectable
				if !m.items[clickedItem].IsCategory {
					m.cursor = clickedItem
					m.contentScroll = 0
				}
			}
		} else {
			// Click in content - focus it
			m.focusSidebar = false
		}

	case tea.MouseWheelUp:
		if !inMainArea {
			return *m, nil
		}

		// Scroll based on mouse position (sidebar or content)
		// Use small increments for smooth, intentional scrolling
		if msg.X < m.sidebarWidth {
			// Scroll sidebar - move 1 item at a time
			m.focusSidebar = true
			m.moveCursor(-1)
			m.contentScroll = 0
		} else {
			// Scroll content - move 1 line at a time
			m.focusSidebar = false
			m.contentScroll--
			if m.contentScroll < 0 {
				m.contentScroll = 0
			}
		}

	case tea.MouseWheelDown:
		if !inMainArea {
			return *m, nil
		}

		// Scroll based on mouse position (sidebar or content)
		// Use small increments for smooth, intentional scrolling
		if msg.X < m.sidebarWidth {
			// Scroll sidebar - move 1 item at a time
			m.focusSidebar = true
			m.moveCursor(1)
			m.contentScroll = 0
		} else {
			// Scroll content - move 1 line at a time
			m.focusSidebar = false
			m.contentScroll++
			// Clamp will happen in render
		}
	}

	_ = mainHeight // Used for bounds checking
	return *m, nil
}

func (m *model) jumpToLast() {
	for i := len(m.items) - 1; i >= 0; i-- {
		if !m.items[i].IsCategory {
			m.cursor = i
			return
		}
	}
}

//
// View
//

func (m model) View() string {
	// Default dimensions for initial render
	width := m.width
	height := m.height
	if width == 0 {
		width = 100
	}
	if height == 0 {
		height = 30
	}

	// Reserve space for header and footer
	headerHeight := 3
	footerHeight := 2
	mainHeight := height - headerHeight - footerHeight

	// Calculate sidebar and content widths (sidebar ~25% but min 24, max 36)
	sidebarWidth := width / 4
	if sidebarWidth < 24 {
		sidebarWidth = 24
	}
	if sidebarWidth > 36 {
		sidebarWidth = 36
	}
	contentWidth := width - sidebarWidth - 1 // 1 for border

	// Render header
	header := m.renderHeader(width)

	// Render sidebar
	sidebar := m.renderSidebar(sidebarWidth, mainHeight)

	// Render content
	content := m.renderContent(contentWidth, mainHeight)

	// Join sidebar and content
	main := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		content,
	)

	// Footer
	footer := m.renderFooter(width)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func (m model) renderHeader(width int) string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(infoColor)

	countStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	searchStyle := lipgloss.NewStyle().
		Foreground(infoColor)

	// Title and count
	title := titleStyle.Render("fp interactive-help")

	filteredCount := countSelectableItems(m.items)
	countText := ""
	if m.searchQuery != "" {
		countText = countStyle.Render(fmt.Sprintf(" (%d/%d items)", filteredCount, m.totalCommands))
	} else {
		countText = countStyle.Render(fmt.Sprintf(" (%d items)", m.totalCommands))
	}

	// Search indicator
	searchText := ""
	if m.searchMode {
		searchText = searchStyle.Render(fmt.Sprintf("  Search: %s_", m.searchQuery))
	} else if m.searchQuery != "" {
		searchText = countStyle.Render(fmt.Sprintf("  Filter: %s", m.searchQuery))
	}

	headerContent := title + countText + searchText

	borderColor := lipgloss.Color(colors.Border)

	headerStyle := lipgloss.NewStyle().
		Width(width).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		MarginBottom(0)

	return headerStyle.Render(headerContent)
}

func (m model) renderSidebar(width, height int) string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	warnColor := lipgloss.Color(colors.Warning)
	themeBorderColor := lipgloss.Color(colors.Border)

	// Topic colors from theme
	topicColors := []lipgloss.Color{
		lipgloss.Color(colors.Color1),
		lipgloss.Color(colors.Color2),
		lipgloss.Color(colors.Color3),
		lipgloss.Color(colors.Color4),
		lipgloss.Color(colors.Color5),
		lipgloss.Color(colors.Color6),
		lipgloss.Color(colors.Color7),
	}

	// Reserve space for scrollbar (2 chars: space + bar)
	scrollbarWidth := 2
	contentWidth := width - scrollbarWidth - 2 // -2 for padding

	visibleHeight := height - 2 // Account for padding

	// Calculate scroll offset to keep cursor visible
	scrollOffset := m.sidebarScroll
	if m.cursor < scrollOffset {
		scrollOffset = m.cursor
	}
	if m.cursor >= scrollOffset+visibleHeight {
		scrollOffset = m.cursor - visibleHeight + 1
	}
	m.sidebarScroll = scrollOffset

	// Handle empty filtered results
	var lines []string
	if len(m.items) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(mutedColor).Italic(true)
		lines = append(lines, emptyStyle.Render("No matches found"))
	}

	// Track topic index for coloring
	topicIndex := 0

	// Build visible lines - one line per item, no extra spacing
	for i, item := range m.items {
		// Skip items before scroll offset
		if i < scrollOffset {
			if item.IsTopic {
				topicIndex++
			}
			continue
		}
		// Stop when we've filled visible area
		if len(lines) >= visibleHeight {
			break
		}

		var line string
		itemWidth := contentWidth - 2 // Account for prefix

		if item.IsCategory {
			// Category header - compact, no extra spacing
			categoryStyle := lipgloss.NewStyle().
				Foreground(mutedColor).
				Bold(true).
				Width(itemWidth + 2)

			line = categoryStyle.Render(item.DisplayName)
		} else {
			// Regular item
			prefix := "  "
			if i == m.cursor {
				if m.focusSidebar {
					prefix = "> "
				} else {
					prefix = "* "
				}
			}

			nameStyle := lipgloss.NewStyle().Width(itemWidth)

			if i == m.cursor {
				// Selected item
				if m.focusSidebar {
					// Focused: bold with background
					nameStyle = nameStyle.
						Bold(true).
						Foreground(lipgloss.Color("0")).
						Background(infoColor)
				} else {
					// Not focused: just highlighted with border color
					nameStyle = nameStyle.
						Bold(true).
						Foreground(warnColor)
				}
				line = prefix + nameStyle.Render(item.DisplayName)
			} else if item.IsTopic {
				// Topic items with theme colors
				colorIdx := topicIndex % len(topicColors)
				nameStyle = nameStyle.Foreground(topicColors[colorIdx])
				line = prefix + nameStyle.Render(item.DisplayName)
			} else {
				// Regular command
				line = prefix + nameStyle.Render(item.DisplayName)
			}

			if item.IsTopic {
				topicIndex++
			}
		}

		lines = append(lines, line)
	}

	// Pad to fill height
	for len(lines) < visibleHeight {
		lines = append(lines, "")
	}

	// Build scrollbar - thumb color depends on focus
	// Focused: bright (infoColor), Unfocused: dim (borderColor)
	thumbColor := themeBorderColor
	if m.focusSidebar {
		thumbColor = infoColor
	}
	scrollbar := m.buildScrollbar(visibleHeight, len(m.items), scrollOffset, thumbColor, themeBorderColor)

	// Combine content with scrollbar
	var combinedLines []string
	for i, line := range lines {
		scrollChar := " "
		if i < len(scrollbar) {
			scrollChar = scrollbar[i]
		}
		// Pad line to fixed width and add scrollbar
		paddedLine := line
		lineLen := lipgloss.Width(line)
		if lineLen < contentWidth {
			paddedLine = line + strings.Repeat(" ", contentWidth-lineLen)
		}
		combinedLines = append(combinedLines, paddedLine+" "+scrollChar)
	}

	sidebarContent := strings.Join(combinedLines, "\n")

	// Border style depends on focus
	borderColor := themeBorderColor
	if m.focusSidebar {
		borderColor = infoColor
	}

	sidebarStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return sidebarStyle.Render(sidebarContent)
}

// Scrollbar characters
const (
	scrollThumbChar = "█" // Full block for thumb (solid)
	scrollTrackChar = "│" // Box drawing vertical for track (hollow/border only)
)

// buildScrollbar creates a visual scrollbar for the given parameters.
// viewHeight: the visible height of the scrollbar track
// totalItems: total number of items/lines in the content
// scrollOffset: current scroll position (0-based)
func (m model) buildScrollbar(viewHeight, totalItems, scrollOffset int, activeColor, trackColor lipgloss.Color) []string {
	scrollbar := make([]string, viewHeight)
	trackStyle := lipgloss.NewStyle().Foreground(trackColor)

	// If all items fit, show only track (no thumb needed)
	if totalItems <= viewHeight {
		for i := range scrollbar {
			scrollbar[i] = trackStyle.Render(scrollTrackChar)
		}
		return scrollbar
	}

	// Calculate thumb size proportional to visible content
	// thumbSize = (visible / total) * trackHeight
	// Make it slightly smaller for better visual appearance
	thumbSize := (viewHeight * viewHeight) / totalItems

	// Ensure minimum size of 1, maximum of viewHeight-2 (leave room for position indication)
	if thumbSize < 1 {
		thumbSize = 1
	}
	maxThumbSize := viewHeight - 2
	if maxThumbSize < 1 {
		maxThumbSize = 1
	}
	if thumbSize > maxThumbSize {
		thumbSize = maxThumbSize
	}

	// Calculate thumb position
	// Position is proportional to scroll offset within scrollable range
	maxScroll := totalItems - viewHeight
	if maxScroll < 1 {
		maxScroll = 1
	}

	// Available track space for thumb movement
	trackSpace := viewHeight - thumbSize
	if trackSpace < 0 {
		trackSpace = 0
	}

	// Calculate position: (scrollOffset / maxScroll) * trackSpace
	thumbPos := 0
	if maxScroll > 0 && trackSpace > 0 {
		thumbPos = (scrollOffset * trackSpace) / maxScroll
	}

	// Clamp thumb position
	if thumbPos < 0 {
		thumbPos = 0
	}
	if thumbPos > trackSpace {
		thumbPos = trackSpace
	}

	// Build scrollbar
	thumbStyle := lipgloss.NewStyle().Foreground(activeColor)

	for i := 0; i < viewHeight; i++ {
		if i >= thumbPos && i < thumbPos+thumbSize {
			scrollbar[i] = thumbStyle.Render(scrollThumbChar)
		} else {
			scrollbar[i] = trackStyle.Render(scrollTrackChar)
		}
	}

	return scrollbar
}

func (m model) renderContent(width, height int) string {
	infoColor := lipgloss.Color(m.colors.Info)
	mutedColor := lipgloss.Color(m.colors.Muted)
	themeBorderColor := lipgloss.Color(m.colors.Border)

	if len(m.items) == 0 || m.cursor >= len(m.items) {
		// No content to show
		emptyStyle := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Padding(1, 2).
			Foreground(mutedColor)
		return emptyStyle.Render("No item selected")
	}

	item := m.items[m.cursor]

	// Reserve space for scrollbar
	scrollbarWidth := 2
	contentWidth := width - scrollbarWidth - 4 // -4 for padding

	var content string
	if item.IsTopic {
		content = m.renderTopicContent(item.Topic, contentWidth)
	} else if item.Node != nil {
		content = m.renderCommandContent(item.Node, contentWidth)
	}

	// Apply scrolling
	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	// Clamp scroll to valid range
	maxScroll := totalLines - (height - 2)
	if maxScroll < 0 {
		maxScroll = 0
	}
	scrollOffset := m.contentScroll
	if scrollOffset > maxScroll {
		scrollOffset = maxScroll
	}
	if scrollOffset < 0 {
		scrollOffset = 0
	}

	if scrollOffset > 0 && scrollOffset < len(lines) {
		lines = lines[scrollOffset:]
	}

	visibleHeight := height - 2

	// Truncate to fit height
	if len(lines) > visibleHeight {
		lines = lines[:visibleHeight]
	}

	// Pad lines to fill visible height
	for len(lines) < visibleHeight {
		lines = append(lines, "")
	}

	// Build scrollbar for content - thumb color depends on focus
	// Focused (content): bright (infoColor), Unfocused: dim (borderColor)
	thumbColor := themeBorderColor
	if !m.focusSidebar {
		thumbColor = infoColor
	}
	scrollbar := m.buildScrollbar(visibleHeight, totalLines, scrollOffset, thumbColor, themeBorderColor)

	// Combine content with scrollbar
	var combinedLines []string
	for i, line := range lines {
		scrollChar := " "
		if i < len(scrollbar) {
			scrollChar = scrollbar[i]
		}
		// Pad line to fixed width and add scrollbar
		paddedLine := line
		lineLen := lipgloss.Width(line)
		if lineLen < contentWidth {
			paddedLine = line + strings.Repeat(" ", contentWidth-lineLen)
		}
		combinedLines = append(combinedLines, paddedLine+" "+scrollChar)
	}

	rendered := strings.Join(combinedLines, "\n")

	// Border style depends on focus
	borderColor := themeBorderColor
	if !m.focusSidebar {
		borderColor = infoColor
	}

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(0, 2)

	return contentStyle.Render(rendered)
}

func (m model) renderCommandContent(node *dispatchers.DispatchNode, width int) string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	successColor := lipgloss.Color(colors.Success)

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(infoColor)

	displayName := strings.Join(node.Path[1:], " ")
	b.WriteString(titleStyle.Render(displayName))
	b.WriteString("\n")

	// Summary
	if node.Summary != "" {
		summaryStyle := lipgloss.NewStyle().Foreground(mutedColor)
		b.WriteString(summaryStyle.Render(node.Summary))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Section header style
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor)

	// Usage
	b.WriteString(headerStyle.Render("USAGE"))
	b.WriteString("\n")

	usageStyle := lipgloss.NewStyle().Foreground(infoColor)
	b.WriteString("   ")
	b.WriteString(usageStyle.Render(node.Usage))
	b.WriteString("\n\n")

	// Description
	if node.Description != "" {
		b.WriteString(headerStyle.Render("DESCRIPTION"))
		b.WriteString("\n")
		b.WriteString(wrapText(node.Description, width-6))
		b.WriteString("\n\n")
	}

	// Flags
	if len(node.Flags) > 0 {
		b.WriteString(headerStyle.Render("FLAGS"))
		b.WriteString("\n")

		flagStyle := lipgloss.NewStyle().Foreground(infoColor)
		descStyle := lipgloss.NewStyle().Foreground(mutedColor)

		for _, f := range node.Flags {
			name := strings.Join(f.Names, ", ")
			if f.ValueHint != "" {
				name = name + " " + f.ValueHint
			}
			b.WriteString("   ")
			b.WriteString(flagStyle.Render(fmt.Sprintf("%-24s", name)))
			b.WriteString("  ")
			b.WriteString(descStyle.Render(f.Description))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Arguments
	if len(node.Args) > 0 {
		b.WriteString(headerStyle.Render("ARGUMENTS"))
		b.WriteString("\n")

		argStyle := lipgloss.NewStyle().Foreground(infoColor)
		descStyle := lipgloss.NewStyle().Foreground(mutedColor)

		for _, a := range node.Args {
			required := ""
			if a.Required {
				required = " (required)"
			}
			b.WriteString("   ")
			b.WriteString(argStyle.Render(fmt.Sprintf("%-16s", a.Name)))
			b.WriteString("  ")
			b.WriteString(descStyle.Render(a.Description + required))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) renderTopicContent(topic *help.Topic, width int) string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	successColor := lipgloss.Color(colors.Success)

	var b strings.Builder

	// Title with topic color
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(infoColor)

	b.WriteString(titleStyle.Render(topic.Name))
	b.WriteString("\n")

	// Summary
	summaryStyle := lipgloss.NewStyle().Foreground(mutedColor)
	b.WriteString(summaryStyle.Render(topic.Summary))
	b.WriteString("\n\n")

	// Section header style (same as commands for consistency)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor)

	b.WriteString(headerStyle.Render("CONTENT"))
	b.WriteString("\n\n")

	// Wrap content to width
	content := topic.Content()
	if width > 0 {
		content = wrapText(content, width-4)
	}
	b.WriteString(content)

	return b.String()
}

func wrapText(text string, width int) string {
	if width <= 0 {
		width = 72
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(line) <= width {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		// Simple word wrap
		words := strings.Fields(line)
		current := ""
		for _, word := range words {
			if current == "" {
				current = word
			} else if len(current)+1+len(word) <= width {
				current += " " + word
			} else {
				result.WriteString(current)
				result.WriteString("\n")
				current = word
			}
		}
		if current != "" {
			result.WriteString(current)
			result.WriteString("\n")
		}
	}

	return strings.TrimSuffix(result.String(), "\n")
}

func (m model) renderFooter(width int) string {
	colors := m.colors
	infoColor := lipgloss.Color(colors.Info)
	mutedColor := lipgloss.Color(colors.Muted)
	borderColor := lipgloss.Color(colors.Border)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(infoColor).
		Padding(0, 1)

	sepStyle := lipgloss.NewStyle().
		Foreground(borderColor)

	labelStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	sep := sepStyle.Render(" | ")

	// Different footer based on mode
	var footer string
	if m.searchMode {
		footer = keyStyle.Render("Enter") + labelStyle.Render(" confirm") + sep +
			keyStyle.Render("Esc") + labelStyle.Render(" cancel")
	} else {
		footer = keyStyle.Render("/") + labelStyle.Render(" search") + sep +
			keyStyle.Render("Tab") + labelStyle.Render(" switch") + sep +
			keyStyle.Render("jk") + labelStyle.Render(" nav") + sep +
			keyStyle.Render("ud") + labelStyle.Render(" scroll") + sep +
			keyStyle.Render("q") + labelStyle.Render(" quit")
	}

	footerStyle := lipgloss.NewStyle().
		Width(width).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return footerStyle.Render(footer)
}
