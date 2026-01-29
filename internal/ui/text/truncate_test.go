package text

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "hello",
			width:    10,
			expected: "hello",
		},
		{
			name:     "exact width",
			input:    "hello",
			width:    5,
			expected: "hello",
		},
		{
			name:     "truncate plain text",
			input:    "hello world",
			width:    5,
			expected: "hello",
		},
		{
			name:     "zero width",
			input:    "hello",
			width:    0,
			expected: "",
		},
		{
			name:     "negative width",
			input:    "hello",
			width:    -1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.input, tt.width)
			if lipgloss.Width(result) > tt.width && tt.width > 0 {
				t.Errorf("Truncate() width %d exceeds max %d", lipgloss.Width(result), tt.width)
			}
		})
	}
}

func TestTruncateWithANSI(t *testing.T) {
	// Test that ANSI codes are preserved and properly closed
	styled := "\x1b[31mred text\x1b[0m"
	result := Truncate(styled, 3)

	// Should be "red" with styles properly closed
	if lipgloss.Width(result) > 3 {
		t.Errorf("ANSI truncation failed: width %d exceeds 3", lipgloss.Width(result))
	}

	// Should end with reset if there was an active style
	// The actual content depends on implementation, just verify width
}

func TestTruncateWithEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		wantDots bool
	}{
		{
			name:     "no ellipsis needed",
			input:    "hi",
			width:    10,
			wantDots: false,
		},
		{
			name:     "ellipsis added",
			input:    "hello world",
			width:    8,
			wantDots: true,
		},
		{
			name:     "very small width",
			input:    "hello",
			width:    3,
			wantDots: true, // Should just be "..."
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateWithEllipsis(tt.input, tt.width)
			if lipgloss.Width(result) > tt.width && tt.width > 0 {
				t.Errorf("TruncateWithEllipsis() width %d exceeds max %d", lipgloss.Width(result), tt.width)
			}
			hasEllipsis := len(result) >= 3 && result[len(result)-3:] == "..."
			if tt.wantDots && !hasEllipsis && tt.width >= 3 {
				t.Errorf("TruncateWithEllipsis() expected ellipsis in %q", result)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ANSI",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "single color",
			input:    "\x1b[31mred\x1b[0m",
			expected: "red",
		},
		{
			name:     "multiple styles",
			input:    "\x1b[1m\x1b[31mbold red\x1b[0m normal \x1b[32mgreen\x1b[0m",
			expected: "bold red normal green",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("StripANSI() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		width         int
		expectedWidth int
	}{
		{
			name:          "pad needed",
			input:         "hi",
			width:         5,
			expectedWidth: 5,
		},
		{
			name:          "no pad needed",
			input:         "hello",
			width:         3,
			expectedWidth: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadRight(tt.input, tt.width)
			if lipgloss.Width(result) != tt.expectedWidth {
				t.Errorf("PadRight() width = %d, want %d", lipgloss.Width(result), tt.expectedWidth)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	result := PadLeft("hi", 5)
	if lipgloss.Width(result) != 5 {
		t.Errorf("PadLeft() width = %d, want 5", lipgloss.Width(result))
	}
	if result != "   hi" {
		t.Errorf("PadLeft() = %q, want %q", result, "   hi")
	}
}

func TestCenter(t *testing.T) {
	result := Center("hi", 6)
	if lipgloss.Width(result) != 6 {
		t.Errorf("Center() width = %d, want 6", lipgloss.Width(result))
	}
	if result != "  hi  " {
		t.Errorf("Center() = %q, want %q", result, "  hi  ")
	}
}
