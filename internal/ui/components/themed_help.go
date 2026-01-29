package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ThemedHelp wraps bubbles/help with theme-aware styling.
type ThemedHelp struct {
	Model help.Model
}

// NewThemedHelp creates a help component styled with the current theme.
func NewThemedHelp() ThemedHelp {
	colors := style.GetColors()

	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.Info)).
		Padding(0, 1)
	h.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Muted))
	h.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.UIDim))
	h.Styles.FullKey = h.Styles.ShortKey
	h.Styles.FullDesc = h.Styles.ShortDesc
	h.Styles.FullSeparator = h.Styles.ShortSeparator

	return ThemedHelp{Model: h}
}

// ShortHelpView renders a single-line help view.
func (t ThemedHelp) ShortHelpView(bindings []key.Binding) string {
	return t.Model.ShortHelpView(bindings)
}

// FullHelpView renders a multi-line help view.
func (t ThemedHelp) FullHelpView(groups [][]key.Binding) string {
	return t.Model.FullHelpView(groups)
}

// View renders the help based on the current show state.
func (t ThemedHelp) View(km help.KeyMap) string {
	return t.Model.View(km)
}

// SetWidth sets the maximum width of the help view.
func (t *ThemedHelp) SetWidth(w int) {
	t.Model.Width = w
}

// RefreshColors updates the help styling from the current theme.
func (t *ThemedHelp) RefreshColors() {
	colors := style.GetColors()
	t.Model.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.Info)).
		Padding(0, 1)
	t.Model.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Muted))
	t.Model.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.UIDim))
	t.Model.Styles.FullKey = t.Model.Styles.ShortKey
	t.Model.Styles.FullDesc = t.Model.Styles.ShortDesc
	t.Model.Styles.FullSeparator = t.Model.Styles.ShortSeparator
}
