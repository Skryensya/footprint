package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobalHooksDir(t *testing.T) {
	dir, err := GlobalHooksDir()
	require.NoError(t, err)
	require.NotEmpty(t, dir)
	require.Contains(t, dir, ".config")
	require.Contains(t, dir, "git")
	require.Contains(t, dir, "hooks")
}

func TestCheckGlobalHooksStatus_NotSet(t *testing.T) {
	// This test assumes core.hooksPath is not set globally
	// Skip if it is set (don't want to mess with user's config)
	current := GetCurrentGlobalHooksPath()
	if current != "" {
		t.Skip("core.hooksPath is already set, skipping test")
	}

	status := CheckGlobalHooksStatus()
	require.False(t, status.IsSet)
	require.Empty(t, status.Path)
}

func TestInstallGlobal_CreatesDirectory(t *testing.T) {
	// Use a temp directory to avoid modifying user's actual global hooks
	tmpDir := t.TempDir()
	hooksDir := filepath.Join(tmpDir, "hooks")

	// We can't actually set core.hooksPath in tests (would affect the system)
	// So we just test the Install part
	err := Install(hooksDir)
	require.NoError(t, err)

	// Verify hooks were created
	for _, hook := range ManagedHooks {
		hookPath := filepath.Join(hooksDir, hook)
		_, err := os.Stat(hookPath)
		require.NoError(t, err, "hook %s should exist", hook)
	}
}

func TestGlobalHooksStatus_String(t *testing.T) {
	tests := []struct {
		status   RepoHookStatus
		expected string
	}{
		{StatusClean, "Clean"},
		{StatusManagedPreCommit, "Managed: pre-commit"},
		{StatusManagedHusky, "Managed: husky"},
		{StatusManagedLefthook, "Managed: lefthook"},
		{StatusUnmanagedHooks, "Unmanaged hooks"},
		{StatusHooksPathOverride, "core.hooksPath set"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestCanInstall(t *testing.T) {
	require.True(t, StatusClean.CanInstall())
	require.False(t, StatusManagedPreCommit.CanInstall())
	require.False(t, StatusManagedHusky.CanInstall())
	require.False(t, StatusManagedLefthook.CanInstall())
	require.False(t, StatusUnmanagedHooks.CanInstall())
	require.False(t, StatusHooksPathOverride.CanInstall())
}
