package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TempConfigFile creates a temporary config file with the given lines.
// Returns the path to the config file. The file is automatically removed when the test finishes.
func TempConfigFile(t *testing.T, lines []string) string {
	t.Helper()

	dir := t.TempDir()
	configPath := filepath.Join(dir, ".fprc")

	content := ""
	for _, line := range lines {
		content += line + "\n"
	}

	err := os.WriteFile(configPath, []byte(content), 0600)
	require.NoError(t, err, "failed to write temp config file")

	return configPath
}

// MockConfigLines creates config line strings from key-value pairs.
// Example: MockConfigLines(map[string]string{"key": "value"}) returns []string{"key=value"}
func MockConfigLines(kvPairs map[string]string) []string {
	lines := make([]string, 0, len(kvPairs))
	for k, v := range kvPairs {
		lines = append(lines, k+"="+v)
	}
	return lines
}
