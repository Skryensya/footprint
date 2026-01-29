// Package style provides semantic terminal styling using lipgloss.
//
// This package is the only place where lipgloss is imported. All styling
// is semantic (Success, Warning, Error, etc.) rather than visual (RedBold, etc.).
//
// When disabled, all helpers return the input string unchanged with no ANSI codes.
package style

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	enabled bool
	colors  ColorConfig

	// Pre-created styles for performance.
	// These are only used when enabled is true.
	successStyle lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
	infoStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	mutedStyle   lipgloss.Style
	borderStyle  lipgloss.Style
	color1Style  lipgloss.Style
	color2Style  lipgloss.Style
	color3Style  lipgloss.Style
	color4Style  lipgloss.Style
	color5Style  lipgloss.Style
	color6Style  lipgloss.Style
	color7Style  lipgloss.Style
)

// Init initializes the style package with the given enabled state and config.
// It also respects NO_COLOR and FP_NO_COLOR environment variables;
// if either is set (to any non-empty value), styling is disabled
// regardless of the enabled parameter.
//
// The cfg parameter is used to load color theme and individual color overrides.
// If cfg is nil, default colors are used.
//
// This function should be called once from main before any output.
func Init(enable bool, cfg map[string]string) {
	// Respect standard NO_COLOR convention and FP-specific override
	if os.Getenv("NO_COLOR") != "" || os.Getenv("FP_NO_COLOR") != "" {
		enabled = false
		return
	}

	enabled = enable

	if enabled {
		colors = LoadColorConfig(cfg)
		initStyles(colors)
	}
}

// GetColors returns the current color configuration.
// Returns empty config if styling is not enabled.
func GetColors() ColorConfig {
	return colors
}

// initStyles creates the lipgloss styles from the given color configuration.
// Uses ANSI 256-color palette to support both basic and extended colors.
func initStyles(colors ColorConfig) {
	// Force lipgloss to use ANSI256 colors regardless of TTY detection.
	// This supports both basic ANSI colors (0-15) and extended 256 colors.
	lipgloss.SetColorProfile(termenv.ANSI256)

	successStyle = makeStyle(colors.Success)
	warningStyle = makeStyle(colors.Warning)
	errorStyle = makeStyle(colors.Error)
	infoStyle = makeStyle(colors.Info)
	mutedStyle = makeStyle(colors.Muted)
	headerStyle = makeStyle(colors.Header)
	borderStyle = makeStyle(colors.Border)
	color1Style = makeStyle(colors.Color1)
	color2Style = makeStyle(colors.Color2)
	color3Style = makeStyle(colors.Color3)
	color4Style = makeStyle(colors.Color4)
	color5Style = makeStyle(colors.Color5)
	color6Style = makeStyle(colors.Color6)
	color7Style = makeStyle(colors.Color7)
}

// makeStyle creates a lipgloss style from a color value.
// The value can be "bold" for bold styling, or an ANSI color number (0-255).
func makeStyle(value string) lipgloss.Style {
	if value == "bold" {
		return lipgloss.NewStyle().Bold(true)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(value))
}

// Enabled returns whether styling is currently enabled.
func Enabled() bool {
	return enabled
}

// render applies a style if styling is enabled, otherwise returns text unchanged.
func render(style lipgloss.Style, text string) string {
	if !enabled {
		return text
	}
	return style.Render(text)
}

// Success styles text for successful operations.
func Success(text string) string { return render(successStyle, text) }

// Warning styles text for warning messages.
func Warning(text string) string { return render(warningStyle, text) }

// Error styles text for error messages.
func Error(text string) string { return render(errorStyle, text) }

// Info styles text for informational messages.
func Info(text string) string { return render(infoStyle, text) }

// Header styles text for section headers or titles.
func Header(text string) string { return render(headerStyle, text) }

// Muted styles text for less important or secondary information.
func Muted(text string) string { return render(mutedStyle, text) }

// Border styles text for interactive delimiters (scrollbars, card borders, etc.).
func Border(text string) string { return render(borderStyle, text) }

// Color1 through Color7 are neutral colors for visual distinction only.
func Color1(text string) string { return render(color1Style, text) }
func Color2(text string) string { return render(color2Style, text) }
func Color3(text string) string { return render(color3Style, text) }
func Color4(text string) string { return render(color4Style, text) }
func Color5(text string) string { return render(color5Style, text) }
func Color6(text string) string { return render(color6Style, text) }
func Color7(text string) string { return render(color7Style, text) }
