package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ThemedSpinner wraps bubbles/spinner with theme-aware styling.
type ThemedSpinner struct {
	Model spinner.Model
}

// NewThemedSpinner creates a spinner styled with the current theme.
func NewThemedSpinner() ThemedSpinner {
	colors := style.GetColors()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Info))

	return ThemedSpinner{Model: s}
}

// NewThemedSpinnerWithType creates a spinner with a specific spinner type.
func NewThemedSpinnerWithType(spinnerType spinner.Spinner) ThemedSpinner {
	colors := style.GetColors()

	s := spinner.New()
	s.Spinner = spinnerType
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Info))

	return ThemedSpinner{Model: s}
}

// Tick returns the spinner tick command to start animation.
func (t ThemedSpinner) Tick() tea.Msg {
	return t.Model.Tick()
}

// Update handles a tea.Msg and returns updated model and command.
func (t ThemedSpinner) Update(msg tea.Msg) (ThemedSpinner, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

// View renders the spinner.
func (t ThemedSpinner) View() string {
	return t.Model.View()
}

// RefreshColors updates the spinner styling from the current theme.
func (t *ThemedSpinner) RefreshColors() {
	colors := style.GetColors()
	t.Model.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Info))
}
