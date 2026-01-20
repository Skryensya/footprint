package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/Skryensya/footprint/internal/paths"
)

func ReadLines() ([]string, error) {
	configPath, err := paths.ConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Ensure correct permissions if file already existed
	_ = os.Chmod(configPath, 0600)

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\r") // Windows CRLF
		lines = append(lines, line)
	}

	return lines, scanner.Err()
}
