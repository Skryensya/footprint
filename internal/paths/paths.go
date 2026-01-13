package paths

import (
	"os"
	"path/filepath"
)

func AppDataDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "."
	}

	path := filepath.Join(dir, "footprint")

	_ = os.MkdirAll(path, 0755)

	return path
}

func ConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".fprc"), nil
}
