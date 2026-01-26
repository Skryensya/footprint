package completions

import (
	"os"
	"path/filepath"
)

// SourceInstructions returns shell-specific instructions for loading completions
func SourceInstructions(shell Shell) string {
	switch shell {
	case ShellBash:
		return `eval "$(fp completions --script)"`
	case ShellZsh:
		return `eval "$(fp completions --script)"`
	case ShellFish:
		return `fp completions --script | source`
	default:
		return ""
	}
}

// RcFile returns the rc file path for the given shell
func RcFile(shell Shell) string {
	switch shell {
	case ShellBash:
		return "~/.bashrc"
	case ShellZsh:
		return "~/.zshrc"
	case ShellFish:
		return "~/.config/fish/config.fish"
	default:
		return ""
	}
}

// AutoInstallPath returns the path where completions can be auto-loaded from.
// Returns empty string if auto-install is not supported for this shell.
func AutoInstallPath(shell Shell) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch shell {
	case ShellFish:
		// Fish always auto-loads from this directory
		return filepath.Join(home, ".config", "fish", "completions", "fp.fish")
	case ShellBash:
		// Only if bash-completion is installed
		if IsBashCompletionInstalled() {
			return filepath.Join(home, ".local", "share", "bash-completion", "completions", "fp")
		}
		return ""
	default:
		return ""
	}
}
