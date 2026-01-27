package completions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// InstallResult describes what happened during installation
type InstallResult struct {
	Shell         Shell
	Installed     bool   // true if file was written
	Path          string // path where file was written (if installed)
	NeedsManual   bool   // true if user needs to add eval to rc file
	Instructions  string // instructions for manual setup (if needed)
}

// InstallSilently attempts to install completions for the running shell.
// For Fish: always writes to auto-load directory
// For Bash: writes to auto-load directory if bash-completion is installed
// For Zsh: returns instructions (no auto-install)
// Returns nil if shell is not detected or not supported.
func InstallSilently() *InstallResult {
	shell := RunningShell()
	if shell == "" {
		return nil
	}

	return InstallForShell(shell)
}

// InstallForShell installs completions for a specific shell
func InstallForShell(shell Shell) *InstallResult {
	root := GetCommandTree()
	if root == nil {
		return nil
	}

	commands := ExtractCommands(root)
	autoPath := AutoInstallPath(shell)

	// Can we auto-install?
	if autoPath != "" {
		script := generateScript(shell, commands)
		if err := writeCompletionFile(autoPath, script); err == nil {
			return &InstallResult{
				Shell:     shell,
				Installed: true,
				Path:      autoPath,
			}
		}
		// If write failed, fall through to manual instructions
	}

	// Need manual installation
	return &InstallResult{
		Shell:        shell,
		NeedsManual:  true,
		Instructions: manualInstructions(shell),
	}
}

// InstallAll attempts to install completions for all detected shells
func InstallAll() []InstallResult {
	var results []InstallResult

	root := GetCommandTree()
	if root == nil {
		return results
	}

	commands := ExtractCommands(root)
	shells := DetectShells()

	for _, shellInfo := range shells {
		shell := shellInfo.Type
		autoPath := AutoInstallPath(shell)

		if autoPath != "" {
			script := generateScript(shell, commands)
			if err := writeCompletionFile(autoPath, script); err == nil {
				results = append(results, InstallResult{
					Shell:     shell,
					Installed: true,
					Path:      autoPath,
				})
				continue
			}
		}

		// Need manual installation
		results = append(results, InstallResult{
			Shell:        shell,
			NeedsManual:  true,
			Instructions: manualInstructions(shell),
		})
	}

	return results
}

func generateScript(shell Shell, commands []CommandInfo) string {
	switch shell {
	case ShellBash:
		return GenerateBash(commands)
	case ShellZsh:
		return GenerateZsh(commands)
	case ShellFish:
		return GenerateFish(commands)
	default:
		return ""
	}
}

func manualInstructions(shell Shell) string {
	rcFile := RcFile(shell)
	sourceLine := SourceInstructions(shell)

	return fmt.Sprintf("To enable completions, add to %s:\n  %s", rcFile, sourceLine)
}

func writeCompletionFile(path, content string) error {
	// Create parent directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write the file
	return os.WriteFile(path, []byte(content), 0644)
}

// PrintCompletions writes the completion script for the given shell to w
func PrintCompletions(w io.Writer, shell Shell) error {
	root := GetCommandTree()
	if root == nil {
		return fmt.Errorf("command tree not registered")
	}

	commands := ExtractCommands(root)
	script := generateScript(shell, commands)
	if script == "" {
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	_, err := fmt.Fprint(w, script)
	return err
}

// AppendToRcFile adds a line to the shell's rc file if not already present
func AppendToRcFile(shell Shell, line string) error {
	rcPath := RcFile(shell)
	if rcPath == "" {
		return fmt.Errorf("unknown rc file for shell: %s", shell)
	}

	// Expand ~ to home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if len(rcPath) > 0 && rcPath[0] == '~' {
		rcPath = filepath.Join(home, rcPath[1:])
	}

	// Create parent directory if needed (e.g., ~/.config/fish/)
	dir := filepath.Dir(rcPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if line already exists
	content, err := os.ReadFile(rcPath)
	if err == nil && strings.Contains(string(content), line) {
		return nil // Already present
	}

	// Append to file
	f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Add newline before if file doesn't end with one
	if len(content) > 0 && content[len(content)-1] != '\n' {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}

	_, err = f.WriteString(line + "\n")
	return err
}
