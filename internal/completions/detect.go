package completions

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Shell represents a supported shell type
type Shell string

const (
	ShellBash Shell = "bash"
	ShellZsh  Shell = "zsh"
	ShellFish Shell = "fish"
)

// ShellInfo contains information about an installed shell
type ShellInfo struct {
	Type Shell
	Path string
}

// DetectShells finds all installed shells
func DetectShells() []ShellInfo {
	var shells []ShellInfo

	shellPaths := map[Shell][]string{
		ShellBash: {"bash"},
		ShellZsh:  {"zsh"},
		ShellFish: {"fish"},
	}

	for shell, names := range shellPaths {
		for _, name := range names {
			if path, err := exec.LookPath(name); err == nil {
				shells = append(shells, ShellInfo{
					Type: shell,
					Path: path,
				})
				break
			}
		}
	}

	return shells
}

// CurrentShell returns the user's current shell based on $SHELL env var
func CurrentShell() Shell {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return ""
	}

	base := filepath.Base(shellEnv)
	switch {
	case strings.Contains(base, "bash"):
		return ShellBash
	case strings.Contains(base, "zsh"):
		return ShellZsh
	case strings.Contains(base, "fish"):
		return ShellFish
	default:
		return ""
	}
}

// RunningShell detects which shell is actually executing the current process
// by checking shell-specific environment variables
func RunningShell() Shell {
	// These env vars are only set when running inside that specific shell
	if os.Getenv("FISH_VERSION") != "" {
		return ShellFish
	}
	if os.Getenv("ZSH_VERSION") != "" {
		return ShellZsh
	}
	if os.Getenv("BASH_VERSION") != "" {
		return ShellBash
	}
	// Fallback to $SHELL
	return CurrentShell()
}

// IsShellAvailable checks if a specific shell is available
func IsShellAvailable(shell Shell) bool {
	shells := DetectShells()
	for _, s := range shells {
		if s.Type == shell {
			return true
		}
	}
	return false
}

// bashCompletionPaths are known locations for bash-completion installation
var bashCompletionPaths = []string{
	// Linux
	"/usr/share/bash-completion/bash_completion",
	"/etc/bash_completion",
	// macOS Homebrew (Apple Silicon)
	"/opt/homebrew/etc/profile.d/bash_completion.sh",
	"/opt/homebrew/etc/bash_completion",
	// macOS Homebrew (Intel)
	"/usr/local/etc/profile.d/bash_completion.sh",
	"/usr/local/etc/bash_completion",
}

// IsBashCompletionInstalled checks if bash-completion package is installed
func IsBashCompletionInstalled() bool {
	for _, path := range bashCompletionPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}
