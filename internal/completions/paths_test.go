package completions

import (
	"strings"
	"testing"
)

func TestSourceInstructions(t *testing.T) {
	tests := []struct {
		shell    Shell
		contains string
	}{
		{ShellBash, "fp completions --script"},
		{ShellZsh, "fp completions --script"},
		{ShellFish, "fp completions --script"},
		{"unknown", ""},
	}

	for _, tc := range tests {
		got := SourceInstructions(tc.shell)
		if tc.contains != "" && !strings.Contains(got, tc.contains) {
			t.Errorf("SourceInstructions(%q) should contain %q, got %q", tc.shell, tc.contains, got)
		}
		if tc.contains == "" && got != "" {
			t.Errorf("SourceInstructions(%q) should be empty, got %q", tc.shell, got)
		}
	}
}

func TestRcFile(t *testing.T) {
	tests := []struct {
		shell Shell
		want  string
	}{
		{ShellBash, "~/.bashrc"},
		{ShellZsh, "~/.zshrc"},
		{ShellFish, "~/.config/fish/config.fish"},
		{"unknown", ""},
	}

	for _, tc := range tests {
		got := RcFile(tc.shell)
		if got != tc.want {
			t.Errorf("RcFile(%q) = %q, want %q", tc.shell, got, tc.want)
		}
	}
}
