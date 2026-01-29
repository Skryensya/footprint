// Package text provides ANSI-safe text manipulation utilities.
package text

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ansiPattern matches ANSI escape sequences
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Truncate truncates a string to maxWidth while preserving ANSI escape codes.
// It ensures that styling is properly terminated and does not break mid-sequence.
func Truncate(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	currentWidth := lipgloss.Width(text)
	if currentWidth <= maxWidth {
		return text
	}

	return truncatePreservingANSI(text, maxWidth)
}

// TruncateWithEllipsis truncates a string with "..." suffix when needed.
// The ellipsis is included in the maxWidth calculation.
func TruncateWithEllipsis(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	currentWidth := lipgloss.Width(text)
	if currentWidth <= maxWidth {
		return text
	}

	if maxWidth <= 3 {
		return strings.Repeat(".", maxWidth)
	}

	truncated := truncatePreservingANSI(text, maxWidth-3)
	return truncated + "..."
}

// truncatePreservingANSI performs the actual truncation while preserving ANSI sequences.
func truncatePreservingANSI(text string, maxWidth int) string {
	var result strings.Builder
	var activeStyles []string
	visibleWidth := 0

	// Find all ANSI sequences and their positions
	matches := ansiPattern.FindAllStringIndex(text, -1)
	lastEnd := 0

	for _, match := range matches {
		// Process text before this ANSI sequence
		if match[0] > lastEnd {
			segment := text[lastEnd:match[0]]
			for _, r := range segment {
				charWidth := runeWidth(r)
				if visibleWidth+charWidth > maxWidth {
					// Close any open styles
					if len(activeStyles) > 0 {
						result.WriteString("\x1b[0m")
					}
					return result.String()
				}
				result.WriteRune(r)
				visibleWidth += charWidth
			}
		}

		// Process the ANSI sequence
		ansi := text[match[0]:match[1]]
		result.WriteString(ansi)

		// Track style state
		if ansi == "\x1b[0m" || ansi == "\x1b[m" {
			activeStyles = nil
		} else {
			activeStyles = append(activeStyles, ansi)
		}

		lastEnd = match[1]
	}

	// Process remaining text after last ANSI sequence
	if lastEnd < len(text) {
		segment := text[lastEnd:]
		for _, r := range segment {
			charWidth := runeWidth(r)
			if visibleWidth+charWidth > maxWidth {
				break
			}
			result.WriteRune(r)
			visibleWidth += charWidth
		}
	}

	// Close any open styles
	if len(activeStyles) > 0 {
		result.WriteString("\x1b[0m")
	}

	return result.String()
}

// runeWidth returns the display width of a rune.
// Most characters are width 1, but some CJK characters are width 2.
func runeWidth(r rune) int {
	// Common CJK ranges that are typically double-width
	if r >= 0x1100 && r <= 0x115F { // Hangul Jamo
		return 2
	}
	if r >= 0x2E80 && r <= 0x9FFF { // CJK blocks
		return 2
	}
	if r >= 0xAC00 && r <= 0xD7A3 { // Hangul Syllables
		return 2
	}
	if r >= 0xF900 && r <= 0xFAFF { // CJK Compatibility Ideographs
		return 2
	}
	if r >= 0xFE10 && r <= 0xFE1F { // Vertical Forms
		return 2
	}
	if r >= 0xFE30 && r <= 0xFE6F { // CJK Compatibility Forms
		return 2
	}
	if r >= 0xFF00 && r <= 0xFF60 { // Fullwidth Forms
		return 2
	}
	if r >= 0xFFE0 && r <= 0xFFE6 { // Fullwidth symbols
		return 2
	}
	return 1
}

// StripANSI removes all ANSI escape codes from a string.
func StripANSI(text string) string {
	return ansiPattern.ReplaceAllString(text, "")
}

// VisualWidth returns the visual width of a string, ignoring ANSI codes.
// This is a convenience wrapper around lipgloss.Width.
func VisualWidth(text string) int {
	return lipgloss.Width(text)
}

// PadRight pads a string to the specified width with spaces on the right.
// It preserves ANSI codes and uses visual width for calculations.
func PadRight(text string, width int) string {
	currentWidth := lipgloss.Width(text)
	if currentWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-currentWidth)
}

// PadLeft pads a string to the specified width with spaces on the left.
// It preserves ANSI codes and uses visual width for calculations.
func PadLeft(text string, width int) string {
	currentWidth := lipgloss.Width(text)
	if currentWidth >= width {
		return text
	}
	return strings.Repeat(" ", width-currentWidth) + text
}

// Center centers a string within the specified width.
// It preserves ANSI codes and uses visual width for calculations.
func Center(text string, width int) string {
	currentWidth := lipgloss.Width(text)
	if currentWidth >= width {
		return text
	}
	padding := width - currentWidth
	leftPad := padding / 2
	rightPad := padding - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}
