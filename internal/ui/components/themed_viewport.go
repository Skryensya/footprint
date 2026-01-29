package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ThemedViewport wraps bubbles/viewport with theme-aware styling.
type ThemedViewport struct {
	Model       viewport.Model
	borderStyle lipgloss.Style
	showBorder  bool
}

// NewThemedViewport creates a viewport with theme-based border styling.
func NewThemedViewport(width, height int) ThemedViewport {
	colors := style.GetColors()

	vp := viewport.New(width, height)
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Border))

	return ThemedViewport{
		Model:       vp,
		borderStyle: borderStyle,
		showBorder:  false,
	}
}

// NewThemedViewportWithBorder creates a viewport with a visible border.
func NewThemedViewportWithBorder(width, height int) ThemedViewport {
	tv := NewThemedViewport(width, height)
	tv.showBorder = true
	return tv
}

// SetContent sets the viewport content.
func (t *ThemedViewport) SetContent(content string) {
	t.Model.SetContent(content)
}

// GotoTop scrolls to the top.
func (t *ThemedViewport) GotoTop() {
	t.Model.GotoTop()
}

// GotoBottom scrolls to the bottom.
func (t *ThemedViewport) GotoBottom() {
	t.Model.GotoBottom()
}

// LineDown scrolls down n lines.
func (t *ThemedViewport) LineDown(n int) {
	t.Model.SetYOffset(t.Model.YOffset + n)
}

// LineUp scrolls up n lines.
func (t *ThemedViewport) LineUp(n int) {
	t.Model.SetYOffset(max(0, t.Model.YOffset-n))
}

// YOffset returns the current scroll position.
func (t ThemedViewport) YOffset() int {
	return t.Model.YOffset
}

// SetYOffset sets the scroll position.
func (t *ThemedViewport) SetYOffset(n int) {
	t.Model.SetYOffset(n)
}

// TotalLineCount returns the total number of lines.
func (t ThemedViewport) TotalLineCount() int {
	return t.Model.TotalLineCount()
}

// VisibleLineCount returns the number of visible lines.
func (t ThemedViewport) VisibleLineCount() int {
	return t.Model.VisibleLineCount()
}

// AtTop returns whether the viewport is scrolled to the top.
func (t ThemedViewport) AtTop() bool {
	return t.Model.AtTop()
}

// AtBottom returns whether the viewport is scrolled to the bottom.
func (t ThemedViewport) AtBottom() bool {
	return t.Model.AtBottom()
}

// SetSize sets the viewport dimensions.
func (t *ThemedViewport) SetSize(width, height int) {
	t.Model.Width = width
	t.Model.Height = height
}

// Update handles a tea.Msg and returns updated model and command.
func (t ThemedViewport) Update(msg tea.Msg) (ThemedViewport, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

// View renders the viewport.
func (t ThemedViewport) View() string {
	if t.showBorder {
		return t.borderStyle.Render(t.Model.View())
	}
	return t.Model.View()
}

// SetBorderFocused updates the border color for focus state.
func (t *ThemedViewport) SetBorderFocused(focused bool) {
	colors := style.GetColors()
	if focused {
		t.borderStyle = t.borderStyle.BorderForeground(lipgloss.Color(colors.UIActive))
	} else {
		t.borderStyle = t.borderStyle.BorderForeground(lipgloss.Color(colors.Border))
	}
}

// RefreshColors updates the viewport styling from the current theme.
func (t *ThemedViewport) RefreshColors() {
	colors := style.GetColors()
	t.borderStyle = t.borderStyle.BorderForeground(lipgloss.Color(colors.Border))
}
