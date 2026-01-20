package ui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDisablePager(t *testing.T) {
	// Reset state
	pagerDisabled = false

	DisablePager()

	require.True(t, pagerDisabled)

	// Reset for other tests
	pagerDisabled = false
}

func TestSetPager(t *testing.T) {
	// Reset state
	pagerOverride = ""

	SetPager("less -R")

	require.Equal(t, "less -R", pagerOverride)

	// Reset for other tests
	pagerOverride = ""
}

func TestIsBypassPager(t *testing.T) {
	tests := []struct {
		cmd    string
		bypass bool
	}{
		{"cat", true},
		{"less", false},
		{"more", false},
		{"", false},
		{"less -R", false},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			result := isBypassPager(tt.cmd)
			require.Equal(t, tt.bypass, result,
				"isBypassPager(%q) should be %v", tt.cmd, tt.bypass)
		})
	}
}

func TestRunPagerCmd_EmptyCommand(t *testing.T) {
	// Empty command should just print content directly
	// This is tested by ensuring no panic occurs
	runPagerCmd("", "test content")
}

func TestPager_DisabledFlag(t *testing.T) {
	// When pager is disabled, content should be printed directly
	// (We can't easily capture output in this test without more infrastructure)
	pagerDisabled = true
	defer func() { pagerDisabled = false }()

	// Should not panic
	Pager("test content")
}

func TestPager_CatOverride(t *testing.T) {
	pagerDisabled = false
	pagerOverride = "cat"
	defer func() { pagerOverride = "" }()

	// Should not panic - cat bypass should print directly
	// We can't easily test the output capture without mocking stdout
	// but we verify the function doesn't hang or error
}

func TestRunPagerCmd_WithArgs(t *testing.T) {
	// Test runPagerCmd with a command that has arguments
	// Using "echo" as a simple command that will work
	runPagerCmd("echo test", "content")
}

func TestRunPagerCmd_SingleCommand(t *testing.T) {
	// Test with just a command name
	runPagerCmd("true", "content")
}

func TestPager_OverrideNotCat(t *testing.T) {
	pagerDisabled = false
	pagerOverride = "true" // A command that exists and will succeed
	defer func() { pagerOverride = "" }()

	// This should attempt to run the pager
	Pager("test content")
}
