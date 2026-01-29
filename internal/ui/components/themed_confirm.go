package components

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// ConfirmResult is sent when the user responds to a confirmation dialog.
type ConfirmResult struct {
	Confirmed bool
	ID        string // Optional ID to identify which confirmation this is
}

// ThemedConfirm is a confirmation dialog with theme-aware styling.
type ThemedConfirm struct {
	title    string
	message  string
	id       string
	width    int
	height   int
	selected int // 0 = Yes, 1 = No
	colors   style.ColorConfig
}

// NewThemedConfirm creates a new confirmation dialog.
func NewThemedConfirm(title, message string) ThemedConfirm {
	return ThemedConfirm{
		title:    title,
		message:  message,
		selected: 1, // Default to "No" for safety
		colors:   style.GetColors(),
	}
}

// WithID sets an identifier for this confirmation.
func (c ThemedConfirm) WithID(id string) ThemedConfirm {
	c.id = id
	return c
}

// Init implements tea.Model.
func (c ThemedConfirm) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (c ThemedConfirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		return c, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return c, func() tea.Msg {
				return ConfirmResult{Confirmed: false, ID: c.id}
			}

		case tea.KeyEnter:
			return c, func() tea.Msg {
				return ConfirmResult{Confirmed: c.selected == 0, ID: c.id}
			}

		case tea.KeyLeft, tea.KeyRight, tea.KeyTab:
			c.selected = 1 - c.selected
			return c, nil

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "y", "Y":
				return c, func() tea.Msg {
					return ConfirmResult{Confirmed: true, ID: c.id}
				}
			case "n", "N":
				return c, func() tea.Msg {
					return ConfirmResult{Confirmed: false, ID: c.id}
				}
			case "h", "l":
				c.selected = 1 - c.selected
				return c, nil
			}
		}
	}

	return c, nil
}

// View implements tea.Model.
func (c ThemedConfirm) View() string {
	colors := c.colors

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Warning)).
		Padding(1, 2).
		Width(40)

	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.Warning)).
		MarginBottom(1)

	// Message style
	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Muted)).
		MarginBottom(1)

	// Button styles
	activeBtn := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.UIActive)).
		Padding(0, 2)

	inactiveBtn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Muted)).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(colors.Border)).
		Padding(0, 1)

	// Build buttons
	var yesBtn, noBtn string
	if c.selected == 0 {
		yesBtn = activeBtn.Render(" Yes ")
		noBtn = inactiveBtn.Render(" No ")
	} else {
		yesBtn = inactiveBtn.Render(" Yes ")
		noBtn = activeBtn.Render(" No ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yesBtn, "  ", noBtn)
	buttonsRow := lipgloss.NewStyle().Width(36).Align(lipgloss.Center).Render(buttons)

	// Compose content
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(c.title),
		msgStyle.Render(c.message),
		"",
		buttonsRow,
	)

	return boxStyle.Render(content)
}

// RefreshColors updates the dialog styling from the current theme.
func (c *ThemedConfirm) RefreshColors() {
	c.colors = style.GetColors()
}

// ConfirmOverlay creates an overlay with a confirmation dialog centered on the background.
func ConfirmOverlay(confirm ThemedConfirm, background tea.Model) *overlay.Model {
	return overlay.New(
		confirm,
		background,
		overlay.Center,
		overlay.Center,
		0, 0,
	)
}

// RenderConfirmOverlay composites a confirmation dialog over a background string.
func RenderConfirmOverlay(confirm ThemedConfirm, background string) string {
	return overlay.Composite(
		confirm.View(),
		background,
		overlay.Center,
		overlay.Center,
		0, 0,
	)
}

// ConfirmKeyBindings returns keybindings for the confirmation dialog.
func ConfirmKeyBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yes")),
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "no")),
		key.NewBinding(key.WithKeys("←", "→"), key.WithHelp("←/→", "select")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "confirm")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc", "cancel")),
	}
}
