package theme

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/ui/style"
)

func List(args []string, flags *dispatchers.ParsedFlags) error {
	return list(args, flags, DefaultDeps())
}

func list(_ []string, _ *dispatchers.ParsedFlags, deps Deps) error {
	current, _ := deps.Get("color_theme")
	if current == "" {
		current = "default-dark"
	}

	for _, name := range deps.ThemeNames {
		marker := "  "
		if name == current {
			marker = style.Success("* ")
		}

		theme := deps.Themes[name]
		preview := renderColorPreview(theme)

		deps.Printf("%s%-14s  %s\n", marker, name, preview)
	}

	deps.Println("\nUse 'fp theme set <name>' to change, 'fp theme show' for details")

	return nil
}

// renderColorPreview returns colored text samples for a theme.
func renderColorPreview(cfg style.ColorConfig) string {
	colorize := func(text string, width int, color string) string {
		padded := fmt.Sprintf("%-*s", width, text)
		if color == "bold" {
			return lipgloss.NewStyle().Bold(true).Render(padded)
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(padded)
	}

	return colorize("success", 9, cfg.Success) +
		colorize("error", 7, cfg.Error) +
		colorize("info", 6, cfg.Info) +
		colorize("muted", 5, cfg.Muted)
}
