package paths

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	return configDir()
}

func ConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".fprc"), nil
}
