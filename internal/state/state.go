package state

import (
	"bufio"
	"os"
	"strings"

	"github.com/footprint-tools/footprint-cli/internal/paths"
)

// Get retrieves a value from the state file.
func Get(key string) (string, error) {
	data, err := readAll()
	if err != nil {
		return "", err
	}
	return data[key], nil
}

// Set stores a value in the state file.
func Set(key, value string) error {
	data, err := readAll()
	if err != nil {
		// If file doesn't exist, start with empty map
		data = make(map[string]string)
	}

	data[key] = value
	return writeAll(data)
}

// readAll reads the state file into a map.
func readAll() (map[string]string, error) {
	statePath := paths.StateFilePath()
	data := make(map[string]string)

	file, err := os.Open(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			data[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return data, scanner.Err()
}

// writeAll writes the map to the state file.
func writeAll(data map[string]string) error {
	statePath := paths.StateFilePath()

	// Ensure directory exists
	dir := paths.AppDataDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	file, err := os.OpenFile(statePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	for key, value := range data {
		if _, err := file.WriteString(key + "=" + value + "\n"); err != nil {
			return err
		}
	}

	return nil
}
