package theme

import (
	"strings"

	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//
// Public API
//

func Pick(args []string, flags *dispatchers.ParsedFlags) error {
	return pick(args, flags, DefaultDeps())
}

//
// Entrypoint
//

func pick(_ []string, _ *dispatchers.ParsedFlags, deps Deps) error {
	current, _ := deps.Get("color_theme")
	if current == "" {
		current = "default-dark"
	}

	currentIdx := 0
	for i, name := range deps.ThemeNames {
		if name == current {
			currentIdx = i
			break
		}
	}

	m := model{
		themes:   deps.ThemeNames,
		configs:  deps.Themes,
		cursor:   currentIdx,
		selected: current,
		deps:     deps,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	fm := finalModel.(model)

	if fm.chosen != "" && fm.chosen != current {
		lines, err := deps.ReadLines()
		if err != nil {
			return err
		}

		lines, _ = deps.Set(lines, "color_theme", fm.chosen)
		if err := deps.WriteLines(lines); err != nil {
			return err
		}

		deps.Printf("\nTheme set to %s\n", style.Success(fm.chosen))
		return nil
	}

	if fm.chosen == "" {
		deps.Println("\nCancelled")
	}

	return nil
}

//
// Model
//

type model struct {
	themes   []string
	configs  map[string]style.ColorConfig
	cursor   int
	selected string
	chosen   string
	deps     Deps
}

//
// Bubble Tea lifecycle
//

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.chosen = m.themes[m.cursor]
			return m, tea.Quit
		}
	}

	return m, nil
}

//
// View
//

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Select a theme:\n\n")

	for i, name := range m.themes {
		cfg := m.configs[name]

		// Cursor
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		// Current marker
		marker := "  "
		if name == m.selected {
			marker = "* "
		}

		// Theme name: fixed visible width (Lipgloss handles ANSI safely)
		nameStyle := lipgloss.NewStyle().Width(14)
		if i == m.cursor {
			nameStyle = nameStyle.Bold(true)
		}
		themeName := nameStyle.Render(name)

		// Color preview (fixed width segments)
		preview := renderPickerPreview(cfg)

		b.WriteString(cursor)
		b.WriteString(marker)
		b.WriteString(themeName)
		b.WriteString("  ")
		b.WriteString(preview)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	cfg := m.configs[m.themes[m.cursor]]
	b.WriteString(renderThemeDetails(m.themes[m.cursor], cfg))

	b.WriteString("\n")
	b.WriteString(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("↑/↓ navigate  enter select  q quit"),
	)

	return b.String()
}

//
// Rendering helpers
//

func renderPickerPreview(cfg style.ColorConfig) string {
	// Nota: dejamos 1 espacio explícito entre columnas para legibilidad.
	return colorize("success", 8, cfg.Success) + " " +
		colorize("warning", 8, cfg.Warning) + " " +
		colorize("error", 6, cfg.Error) + " " +
		colorize("info", 5, cfg.Info) + " " +
		colorize("muted", 5, cfg.Muted)
}

func renderThemeDetails(name string, cfg style.ColorConfig) string {
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	label := lipgloss.NewStyle().Width(14).Align(lipgloss.Left)

	source := func(text string, color string) string {
		return lipgloss.NewStyle().
			Width(18).
			Align(lipgloss.Left).
			Foreground(lipgloss.Color(color)).
			Render(text)
	}

	var b strings.Builder

	b.WriteString(muted.Render("Preview: "))
	b.WriteString(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Info)).
			Render(name),
	)
	b.WriteString("\n\n")

	// UI roles
	b.WriteString("  ")
	b.WriteString(label.Render("success"))
	b.WriteString(colorize("confirmations", 18, cfg.Success))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(label.Render("warning"))
	b.WriteString(colorize("cautionary messages", 19, cfg.Warning))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(label.Render("error"))
	b.WriteString(colorize("error messages", 18, cfg.Error))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(label.Render("info"))
	b.WriteString(colorize("highlights", 18, cfg.Info))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(label.Render("muted"))
	b.WriteString(colorize("secondary text", 18, cfg.Muted))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(label.Render("header"))
	b.WriteString(colorize("commit hashes", 18, cfg.Header))
	b.WriteString("\n\n")

	// Sources (canonical, aligned list)
	b.WriteString("  ")
	b.WriteString(muted.Render("Sources:"))
	b.WriteString("\n")

	// IMPORTANT: same prefix for every row
	const srcIndent = "    "

	b.WriteString(srcIndent)
	b.WriteString(source("POST-COMMIT", cfg.Color1))
	b.WriteString("\n")

	b.WriteString(srcIndent)
	b.WriteString(source("POST-REWRITE", cfg.Color2))
	b.WriteString("\n")

	b.WriteString(srcIndent)
	b.WriteString(source("POST-CHECKOUT", cfg.Color3))
	b.WriteString("\n")

	b.WriteString(srcIndent)
	b.WriteString(source("POST-MERGE", cfg.Color4))
	b.WriteString("\n")

	b.WriteString(srcIndent)
	b.WriteString(source("PRE-PUSH", cfg.Color5))
	b.WriteString("\n")

	b.WriteString(srcIndent)
	b.WriteString(source("MANUAL", cfg.Color6))
	b.WriteString("\n")

	b.WriteString(srcIndent)
	b.WriteString(source("BACKFILL", cfg.Color6))

	return b.String()
}

//
// Tiny helpers
//

func colorize(text string, width int, color string) string {
	st := lipgloss.NewStyle().Width(width)

	if color == "bold" {
		st = st.Bold(true)
	} else {
		st = st.Foreground(lipgloss.Color(color))
	}

	return st.Render(text)
}
