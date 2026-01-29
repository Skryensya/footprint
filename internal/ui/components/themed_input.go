package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ThemedInput wraps bubbles/textinput with theme-aware styling.
type ThemedInput struct {
	Model textinput.Model
}

// NewThemedInput creates a new text input styled with the current theme.
func NewThemedInput(placeholder string) ThemedInput {
	colors := style.GetColors()

	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Muted))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Info))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.UIActive))
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.UIActive))

	return ThemedInput{Model: ti}
}

// NewThemedInputWithPrompt creates a text input with a custom prompt.
func NewThemedInputWithPrompt(placeholder, prompt string) ThemedInput {
	ti := NewThemedInput(placeholder)
	ti.Model.Prompt = prompt
	return ti
}

// Focus sets the input to focused state.
func (t *ThemedInput) Focus() tea.Cmd {
	return t.Model.Focus()
}

// Blur removes focus from the input.
func (t *ThemedInput) Blur() {
	t.Model.Blur()
}

// Focused returns whether the input is focused.
func (t ThemedInput) Focused() bool {
	return t.Model.Focused()
}

// SetValue sets the input value.
func (t *ThemedInput) SetValue(s string) {
	t.Model.SetValue(s)
}

// Value returns the current input value.
func (t ThemedInput) Value() string {
	return t.Model.Value()
}

// SetCursor sets the cursor position.
func (t *ThemedInput) SetCursor(pos int) {
	t.Model.SetCursor(pos)
}

// CursorEnd moves the cursor to the end.
func (t *ThemedInput) CursorEnd() {
	t.Model.CursorEnd()
}

// Reset clears the input value.
func (t *ThemedInput) Reset() {
	t.Model.Reset()
}

// Update handles a tea.Msg and returns updated model and command.
func (t ThemedInput) Update(msg tea.Msg) (ThemedInput, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

// View renders the input.
func (t ThemedInput) View() string {
	return t.Model.View()
}

// SetWidth sets the width of the input.
func (t *ThemedInput) SetWidth(w int) {
	t.Model.Width = w
}

// RefreshColors updates the input styling from the current theme.
// Call this if the theme changes during the session.
func (t *ThemedInput) RefreshColors() {
	colors := style.GetColors()
	t.Model.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Muted))
	t.Model.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Info))
	t.Model.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.UIActive))
	t.Model.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.UIActive))
}
