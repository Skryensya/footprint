package config

import (
	"bufio"
	"os"

	"github.com/Skryensya/footprint/internal/paths"
)

func WriteLines(lines []string) error {
	configPath, err := paths.ConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		configPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
