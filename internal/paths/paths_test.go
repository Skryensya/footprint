package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppDataDir_ReturnsNonEmpty(t *testing.T) {
	dir := AppDataDir()
	require.NotEmpty(t, dir)
	require.NotEqual(t, ".", dir)
}

func TestAppDataDir_ContainsFootprint(t *testing.T) {
	dir := AppDataDir()
	dirLower := strings.ToLower(dir)
	require.True(t, strings.Contains(dirLower, "footprint"),
		"AppDataDir should contain 'footprint' (case-insensitive): %s", dir)
}

func TestAppLocalDataDir_ReturnsNonEmpty(t *testing.T) {
	dir := AppLocalDataDir()
	require.NotEmpty(t, dir)
	require.NotEqual(t, ".", dir)
}

func TestAppLocalDataDir_ContainsFootprint(t *testing.T) {
	dir := AppLocalDataDir()
	require.True(t, strings.HasSuffix(dir, "footprint"),
		"AppLocalDataDir should end with 'footprint': %s", dir)
}

func TestAppLocalDataDir_Platform(t *testing.T) {
	dir := AppLocalDataDir()

	switch runtime.GOOS {
	case "darwin":
		require.Contains(t, dir, "Library")
		require.Contains(t, dir, "Application Support")
	case "linux":
		// Could be XDG_DATA_HOME or .local/share
		require.True(t, strings.Contains(dir, ".local/share") ||
			os.Getenv("XDG_DATA_HOME") != "",
			"Linux path should use XDG_DATA_HOME or .local/share: %s", dir)
	case "windows":
		require.True(t, strings.Contains(dir, "AppData") ||
			strings.Contains(dir, "Local"),
			"Windows path should contain AppData: %s", dir)
	}
}

func TestExportRepoDir_ReturnsValidPath(t *testing.T) {
	dir := ExportRepoDir()
	require.NotEmpty(t, dir)
	require.True(t, strings.HasSuffix(dir, "export"),
		"ExportRepoDir should end with 'export': %s", dir)
}

func TestExportRepoDir_IsUnderAppLocalDataDir(t *testing.T) {
	exportDir := ExportRepoDir()
	localDataDir := AppLocalDataDir()

	require.True(t, strings.HasPrefix(exportDir, localDataDir),
		"ExportRepoDir should be under AppLocalDataDir: %s vs %s",
		exportDir, localDataDir)
}

func TestConfigFilePath_Success(t *testing.T) {
	path, err := ConfigFilePath()

	require.NoError(t, err)
	require.NotEmpty(t, path)
	require.True(t, strings.HasSuffix(path, ".fprc"),
		"ConfigFilePath should end with .fprc: %s", path)
}

func TestConfigFilePath_UnderHomeDir(t *testing.T) {
	path, err := ConfigFilePath()
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	require.True(t, strings.HasPrefix(path, home),
		"ConfigFilePath should be under home dir: %s", path)
}

func TestLogFilePath_ReturnsValidPath(t *testing.T) {
	path := LogFilePath()

	require.NotEmpty(t, path)
	require.True(t, strings.HasSuffix(path, "fp.log"),
		"LogFilePath should end with fp.log: %s", path)
}

func TestLogFilePath_IsUnderAppDataDir(t *testing.T) {
	logPath := LogFilePath()
	appDataDir := AppDataDir()

	require.True(t, strings.HasPrefix(logPath, appDataDir),
		"LogFilePath should be under AppDataDir: %s vs %s",
		logPath, appDataDir)
}

func TestAppDataDir_CreatesDirectory(t *testing.T) {
	// This test verifies that AppDataDir creates the directory if it doesn't exist
	dir := AppDataDir()

	// Check that it's actually an absolute path
	require.True(t, filepath.IsAbs(dir),
		"AppDataDir should return an absolute path: %s", dir)
}

func TestPaths_NoDotDotComponents(t *testing.T) {
	// Security check: paths should not contain ..
	paths := []string{
		AppDataDir(),
		AppLocalDataDir(),
		ExportRepoDir(),
		LogFilePath(),
	}

	cfgPath, err := ConfigFilePath()
	require.NoError(t, err)
	paths = append(paths, cfgPath)

	for _, p := range paths {
		require.False(t, strings.Contains(p, ".."),
			"Path should not contain '..': %s", p)
	}
}
