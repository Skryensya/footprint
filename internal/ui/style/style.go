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

	// Pre-created styles for performance.
	// These are only used when enabled is true.
	successStyle lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
	infoStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	mutedStyle   lipgloss.Style
	color1Style  lipgloss.Style
	color2Style  lipgloss.Style
	color3Style  lipgloss.Style
	color4Style  lipgloss.Style
	color5Style  lipgloss.Style
	color6Style  lipgloss.Style
)

// Init initializes the style package with the given enabled state.
// It also respects NO_COLOR and FP_NO_COLOR environment variables;
// if either is set (to any non-empty value), styling is disabled
// regardless of the enabled parameter.
//
// This function should be called once from main before any output.
func Init(enable bool) {
	// Respect standard NO_COLOR convention and FP-specific override
	if os.Getenv("NO_COLOR") != "" || os.Getenv("FP_NO_COLOR") != "" {
		enabled = false
		return
	}

	enabled = enable

	if enabled {
		initStyles()
	}
}

// initStyles creates the lipgloss styles.
// Colors are intentionally subtle and readable.
// Uses ANSI 16-color palette for broad terminal compatibility.
func initStyles() {
	// Force lipgloss to use ANSI colors regardless of TTY detection.
	// This is safe because the caller (main) has already decided colors are appropriate.
	lipgloss.SetColorProfile(termenv.ANSI)

	// Green for success - not too bright
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	// Yellow for warnings
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	// Red for errors
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	// Cyan for informational messages
	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	// Bold for headers - emphasis without color
	headerStyle = lipgloss.NewStyle().Bold(true)

	// Dim/gray for muted text
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	// Color styles for visual distinction
	color1Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan
	color2Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // magenta
	color3Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4")) // blue
	color4Style = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	color5Style = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	color6Style = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
}

// Enabled returns whether styling is currently enabled.
func Enabled() bool {
	return enabled
}

// Success styles text for successful operations.
func Success(text string) string {
	if !enabled {
		return text
	}
	return successStyle.Render(text)
}

// Warning styles text for warning messages.
func Warning(text string) string {
	if !enabled {
		return text
	}
	return warningStyle.Render(text)
}

// Error styles text for error messages.
func Error(text string) string {
	if !enabled {
		return text
	}
	return errorStyle.Render(text)
}

// Info styles text for informational messages.
func Info(text string) string {
	if !enabled {
		return text
	}
	return infoStyle.Render(text)
}

// Header styles text for section headers or titles.
func Header(text string) string {
	if !enabled {
		return text
	}
	return headerStyle.Render(text)
}

// Muted styles text for less important or secondary information.
func Muted(text string) string {
	if !enabled {
		return text
	}
	return mutedStyle.Render(text)
}

// Color1 through Color6 are neutral colors for visual distinction only.
// They have no semantic meaning.

func Color1(text string) string {
	if !enabled {
		return text
	}
	return color1Style.Render(text)
}

func Color2(text string) string {
	if !enabled {
		return text
	}
	return color2Style.Render(text)
}

func Color3(text string) string {
	if !enabled {
		return text
	}
	return color3Style.Render(text)
}

func Color4(text string) string {
	if !enabled {
		return text
	}
	return color4Style.Render(text)
}

func Color5(text string) string {
	if !enabled {
		return text
	}
	return color5Style.Render(text)
}

func Color6(text string) string {
	if !enabled {
		return text
	}
	return color6Style.Render(text)
}
