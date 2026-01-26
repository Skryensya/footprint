package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RepoHooksPath(repoRoot string) (string, error) {
	cmd := exec.Command("git", "-C", repoRoot, "rev-parse", "--git-path", "hooks")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	hooksPath := filepath.Clean(strings.TrimSpace(string(out)))

	// If the path is relative, join it with repoRoot
	if !filepath.IsAbs(hooksPath) {
		hooksPath = filepath.Join(repoRoot, hooksPath)
	}

	// Ensure the path is absolute
	return filepath.Abs(hooksPath)
}

func GlobalHooksPath() (string, error) {
	cmd := exec.Command("git", "config", "--global", "core.hooksPath")
	out, err := cmd.Output()
	if err == nil {
		path := strings.TrimSpace(string(out))
		if path != "" {
			return path, nil
		}
	}

	homeCmd := exec.Command("git", "config", "--global", "--path", "core.hooksPath")
	homeOut, _ := homeCmd.Output()

	if strings.TrimSpace(string(homeOut)) != "" {
		return strings.TrimSpace(string(homeOut)), nil
	}

	return filepath.Join(defaultHome(), ".git", "hooks"), nil
}

func defaultHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}
