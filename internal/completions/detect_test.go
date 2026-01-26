package completions

import (
	"os"
	"testing"
)

func TestDetectShells(t *testing.T) {
	shells := DetectShells()
	// At minimum, we should find at least one shell on any Unix-like system
	// This test is environment-dependent but should pass on macOS/Linux
	if len(shells) == 0 {
		t.Log("Warning: no shells detected (this may be expected in some environments)")
	}

	// Verify the structure is correct
	for _, s := range shells {
		if s.Type == "" {
			t.Error("shell type should not be empty")
		}
		if s.Path == "" {
			t.Error("shell path should not be empty")
		}
	}
}

func TestCurrentShell(t *testing.T) {
	// Save and restore original SHELL
	origShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", origShell)

	tests := []struct {
		shellEnv string
		want     Shell
	}{
		{"/bin/bash", ShellBash},
		{"/usr/local/bin/bash", ShellBash},
		{"/bin/zsh", ShellZsh},
		{"/usr/bin/zsh", ShellZsh},
		{"/usr/local/bin/fish", ShellFish},
		{"/bin/fish", ShellFish},
		{"/bin/sh", ""},
		{"", ""},
	}

	for _, tc := range tests {
		os.Setenv("SHELL", tc.shellEnv)
		got := CurrentShell()
		if got != tc.want {
			t.Errorf("CurrentShell() with SHELL=%q: got %q, want %q", tc.shellEnv, got, tc.want)
		}
	}
}

func TestIsShellAvailable(t *testing.T) {
	// This test is environment-dependent
	// On most systems, at least bash should be available
	shells := DetectShells()
	if len(shells) > 0 {
		// If we detected any shell, IsShellAvailable should return true for it
		if !IsShellAvailable(shells[0].Type) {
			t.Errorf("IsShellAvailable(%q) should return true", shells[0].Type)
		}
	}

	// An unknown shell should not be available
	if IsShellAvailable("unknown-shell") {
		t.Error("IsShellAvailable should return false for unknown shell")
	}
}
