package text

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

// Wrap wraps text to the specified width, breaking on word boundaries when possible.
// It preserves ANSI codes and handles multi-line input.
func Wrap(text string, width int) string {
	if width <= 0 {
		return ""
	}

	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		wrapped := wrapLine(line, width)
		result = append(result, wrapped...)
	}

	return strings.Join(result, "\n")
}

// WrapPreserveIndent wraps text while preserving the indentation of each line.
// Continuation lines maintain the same indentation as their source line.
func WrapPreserveIndent(text string, width int) string {
	if width <= 0 {
		return ""
	}

	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		indent := extractIndent(line)
		indentWidth := len(indent)
		contentWidth := width - indentWidth

		if contentWidth <= 0 {
			// Not enough room, just include the line as-is
			result = append(result, line)
			continue
		}

		content := strings.TrimPrefix(line, indent)
		wrapped := wrapLine(content, contentWidth)

		for i, w := range wrapped {
			if i == 0 {
				result = append(result, indent+w)
			} else {
				result = append(result, indent+w)
			}
		}
	}

	return strings.Join(result, "\n")
}

// wrapLine wraps a single line to the specified width.
func wrapLine(line string, width int) []string {
	if lipgloss.Width(line) <= width {
		return []string{line}
	}

	var result []string
	var currentLine strings.Builder
	var currentWidth int
	var activeStyles []string

	words := splitIntoWords(line)

	for _, word := range words {
		wordWidth := lipgloss.Width(word.text)

		// Check if word fits on current line
		if currentWidth > 0 && currentWidth+wordWidth > width {
			// Finalize current line
			lineStr := currentLine.String()
			if len(activeStyles) > 0 {
				lineStr += "\x1b[0m"
			}
			result = append(result, lineStr)

			// Start new line with active styles
			currentLine.Reset()
			for _, style := range activeStyles {
				currentLine.WriteString(style)
			}
			currentWidth = 0
		}

		// Handle words that are longer than the width
		if wordWidth > width && currentWidth == 0 {
			chunks := splitLongWord(word.text, width)
			for i, chunk := range chunks {
				if i > 0 {
					lineStr := currentLine.String()
					if len(activeStyles) > 0 {
						lineStr += "\x1b[0m"
					}
					result = append(result, lineStr)
					currentLine.Reset()
					for _, style := range activeStyles {
						currentLine.WriteString(style)
					}
				}
				currentLine.WriteString(chunk)
				currentWidth = lipgloss.Width(chunk)
			}
		} else {
			currentLine.WriteString(word.text)
			currentWidth += wordWidth
		}

		// Track ANSI styles in this word
		for _, style := range word.styles {
			if style == "\x1b[0m" || style == "\x1b[m" {
				activeStyles = nil
			} else {
				activeStyles = append(activeStyles, style)
			}
		}
	}

	// Add final line if non-empty
	if currentLine.Len() > 0 {
		lineStr := currentLine.String()
		if len(activeStyles) > 0 {
			lineStr += "\x1b[0m"
		}
		result = append(result, lineStr)
	}

	if len(result) == 0 {
		return []string{""}
	}

	return result
}

// word represents a word with its ANSI styles
type word struct {
	text   string
	styles []string
}

// splitIntoWords splits text into words, preserving ANSI codes.
func splitIntoWords(text string) []word {
	var words []word
	var current strings.Builder
	var currentStyles []string
	inWord := false

	matches := ansiPattern.FindAllStringIndex(text, -1)
	ansiRanges := make(map[int]int)
	for _, m := range matches {
		ansiRanges[m[0]] = m[1]
	}

	i := 0
	for i < len(text) {
		// Check if we're at an ANSI sequence
		if end, ok := ansiRanges[i]; ok {
			ansi := text[i:end]
			current.WriteString(ansi)
			if ansi == "\x1b[0m" || ansi == "\x1b[m" {
				currentStyles = nil
			} else {
				currentStyles = append(currentStyles, ansi)
			}
			i = end
			continue
		}

		r := rune(text[i])
		if unicode.IsSpace(r) {
			if inWord {
				words = append(words, word{text: current.String(), styles: copyStyles(currentStyles)})
				current.Reset()
				// Preserve styles for next word
				for _, s := range currentStyles {
					current.WriteString(s)
				}
				inWord = false
			}
			// Include space as its own "word"
			current.WriteRune(r)
			words = append(words, word{text: current.String(), styles: copyStyles(currentStyles)})
			current.Reset()
			for _, s := range currentStyles {
				current.WriteString(s)
			}
		} else {
			current.WriteRune(r)
			inWord = true
		}
		i++
	}

	// Add final word
	if current.Len() > 0 {
		words = append(words, word{text: current.String(), styles: copyStyles(currentStyles)})
	}

	return words
}

// copyStyles creates a copy of the styles slice
func copyStyles(styles []string) []string {
	if len(styles) == 0 {
		return nil
	}
	result := make([]string, len(styles))
	copy(result, styles)
	return result
}

// splitLongWord splits a word that's longer than width into chunks.
func splitLongWord(text string, width int) []string {
	var chunks []string
	var current strings.Builder
	var activeStyles []string
	visibleWidth := 0

	matches := ansiPattern.FindAllStringIndex(text, -1)
	ansiRanges := make(map[int]int)
	for _, m := range matches {
		ansiRanges[m[0]] = m[1]
	}

	i := 0
	for i < len(text) {
		// Check if we're at an ANSI sequence
		if end, ok := ansiRanges[i]; ok {
			ansi := text[i:end]
			current.WriteString(ansi)
			if ansi == "\x1b[0m" || ansi == "\x1b[m" {
				activeStyles = nil
			} else {
				activeStyles = append(activeStyles, ansi)
			}
			i = end
			continue
		}

		r := rune(text[i])
		charWidth := runeWidth(r)

		if visibleWidth+charWidth > width && visibleWidth > 0 {
			// End current chunk
			if len(activeStyles) > 0 {
				current.WriteString("\x1b[0m")
			}
			chunks = append(chunks, current.String())

			// Start new chunk with active styles
			current.Reset()
			for _, s := range activeStyles {
				current.WriteString(s)
			}
			visibleWidth = 0
		}

		current.WriteRune(r)
		visibleWidth += charWidth
		i++
	}

	// Add final chunk
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}

// extractIndent returns the leading whitespace of a string.
func extractIndent(s string) string {
	var indent strings.Builder
	for _, r := range s {
		if r == ' ' || r == '\t' {
			indent.WriteRune(r)
		} else {
			break
		}
	}
	return indent.String()
}
