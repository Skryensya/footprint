package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/ui/style"
)

func Show(args []string, flags *dispatchers.ParsedFlags) error {
	return show(args, flags, DefaultDeps())
}

func show(_ []string, _ *dispatchers.ParsedFlags, deps Deps) error {
	cfg, _ := deps.GetAll()
	colors := style.LoadColorConfig(cfg)

	themeName, _ := deps.Get("color_theme")
	if themeName == "" {
		themeName = style.ResolveThemeName("default")
	}

	// Helper to colorize text
	colorize := func(text, color string) string {
		if color == "bold" {
			return lipgloss.NewStyle().Bold(true).Render(text)
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
	}

	deps.Printf("Theme: %s\n\n", colorize(themeName, colors.Info))

	deps.Println("UI colors:")
	deps.Printf("  %s - confirmations, current theme marker\n", colorize("success", colors.Success))
	deps.Printf("  %s - error messages\n", colorize("error", colors.Error))
	deps.Printf("  %s - command names, highlighted text\n", colorize("info", colors.Info))
	deps.Printf("  %s - secondary info (timestamps, repo IDs)\n", colorize("muted", colors.Muted))
	deps.Printf("  %s - commit hashes\n", colorize("header", colors.Header))

	deps.Println("\nGit event sources:")
	deps.Printf("  %s, %s, %s, %s, %s, %s, %s\n",
		colorize("POST-COMMIT", colors.Color1),
		colorize("POST-REWRITE", colors.Color2),
		colorize("POST-CHECKOUT", colors.Color3),
		colorize("POST-MERGE", colors.Color4),
		colorize("PRE-PUSH", colors.Color5),
		colorize("BACKFILL", colors.Color6),
		colorize("MANUAL", colors.Color7))

	deps.Println("\nOverride: fp config set color_<name> <value>")

	return nil
}
