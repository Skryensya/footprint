package setup

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/footprint-tools/footprint-cli/internal/dispatchers"
)

// =========== SETUP TESTS ===========

func TestSetup_GlobalSuccess(t *testing.T) {
	var installedPath string
	var printedLines []string
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit":   false,
				"post-merge":    false,
				"post-checkout": false,
			}
		},
		HooksInstall: func(path string) error {
			installedPath = path
			return nil
		},
		Printf: func(format string, a ...any) (int, error) {
			printedLines = append(printedLines, fmt.Sprintf(format, a...))
			return 0, nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := setup([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/home/user/.config/git/hooks", installedPath)
	require.True(t, len(printedLines) > 0)
	require.Contains(t, printedLines[0], "global")
}

func TestSetup_RepoSuccess(t *testing.T) {
	var installedPath string
	var printedLines []string
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "/path/to/repo", nil
		},
		RepoHooksPath: func(root string) (string, error) {
			return "/path/to/repo/.git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{}
		},
		HooksInstall: func(path string) error {
			installedPath = path
			return nil
		},
		Printf: func(format string, a ...any) (int, error) {
			printedLines = append(printedLines, fmt.Sprintf(format, a...))
			return 0, nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := setup([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/path/to/repo/.git/hooks", installedPath)
	require.True(t, len(printedLines) > 0)
	require.Contains(t, printedLines[0], "repository")
}

func TestSetup_NotInGitRepo(t *testing.T) {
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "", errors.New("not a git repo")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := setup([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "git repo")
}

func TestSetup_ExistingHooksWithConfirmation(t *testing.T) {
	var installedPath string
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit": true, // Existing hook
			}
		},
		HooksInstall: func(path string) error {
			installedPath = path
			return nil
		},
		Printf: func(format string, a ...any) (int, error) {
			return 0, nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			if ptr, ok := a[0].(*string); ok {
				*ptr = "y"
			}
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := setup([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/home/user/.config/git/hooks", installedPath)
}

func TestSetup_ExistingHooksDeclined(t *testing.T) {
	installCalled := false
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit": true,
			}
		},
		HooksInstall: func(path string) error {
			installCalled = true
			return nil
		},
		Printf: func(format string, a ...any) (int, error) {
			return 0, nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			if ptr, ok := a[0].(*string); ok {
				*ptr = "n"
			}
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := setup([]string{}, flags, deps)

	require.NoError(t, err)
	require.False(t, installCalled, "should not install when user declines")
}

func TestSetup_ExistingHooksWithForce(t *testing.T) {
	var installedPath string
	scanlnCalled := false
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit": true,
			}
		},
		HooksInstall: func(path string) error {
			installedPath = path
			return nil
		},
		Printf: func(format string, a ...any) (int, error) {
			return 0, nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			scanlnCalled = true
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--force"})
	err := setup([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/home/user/.config/git/hooks", installedPath)
	require.False(t, scanlnCalled, "should not prompt with --force")
}

func TestSetup_GlobalHooksPathError(t *testing.T) {
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "", errors.New("cannot get global hooks path")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := setup([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot get global hooks path")
}

func TestSetup_RepoHooksPathError(t *testing.T) {
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "/path/to/repo", nil
		},
		RepoHooksPath: func(root string) (string, error) {
			return "", errors.New("cannot get repo hooks path")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := setup([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot get repo hooks path")
}

func TestSetup_InstallError(t *testing.T) {
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{}
		},
		HooksInstall: func(path string) error {
			return errors.New("cannot install hooks")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := setup([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot install hooks")
}

// =========== TEARDOWN TESTS ===========

func TestTeardown_GlobalSuccess(t *testing.T) {
	var uninstalledPath string
	var printedLine string
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksUninstall: func(path string) error {
			uninstalledPath = path
			return nil
		},
		Println: func(a ...any) (int, error) {
			if len(a) > 0 {
				printedLine = fmt.Sprint(a...)
			}
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			if ptr, ok := a[0].(*string); ok {
				*ptr = "y"
			}
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := teardown([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/home/user/.config/git/hooks", uninstalledPath)
	require.Contains(t, printedLine, "teardown complete")
}

func TestTeardown_RepoSuccess(t *testing.T) {
	var uninstalledPath string
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "/path/to/repo", nil
		},
		RepoHooksPath: func(root string) (string, error) {
			return "/path/to/repo/.git/hooks", nil
		},
		HooksUninstall: func(path string) error {
			uninstalledPath = path
			return nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			if ptr, ok := a[0].(*string); ok {
				*ptr = "yes"
			}
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := teardown([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/path/to/repo/.git/hooks", uninstalledPath)
}

func TestTeardown_NotInGitRepo(t *testing.T) {
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "", errors.New("not a git repo")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := teardown([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "git repo")
}

func TestTeardown_Declined(t *testing.T) {
	uninstallCalled := false
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksUninstall: func(path string) error {
			uninstallCalled = true
			return nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			if ptr, ok := a[0].(*string); ok {
				*ptr = "n"
			}
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := teardown([]string{}, flags, deps)

	require.NoError(t, err)
	require.False(t, uninstallCalled, "should not uninstall when user declines")
}

func TestTeardown_WithForce(t *testing.T) {
	var uninstalledPath string
	scanlnCalled := false
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksUninstall: func(path string) error {
			uninstalledPath = path
			return nil
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			scanlnCalled = true
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--force"})
	err := teardown([]string{}, flags, deps)

	require.NoError(t, err)
	require.Equal(t, "/home/user/.config/git/hooks", uninstalledPath)
	require.False(t, scanlnCalled, "should not prompt with --force")
}

func TestTeardown_GlobalHooksPathError(t *testing.T) {
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "", errors.New("cannot get global hooks path")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := teardown([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot get global hooks path")
}

func TestTeardown_UninstallError(t *testing.T) {
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksUninstall: func(path string) error {
			return errors.New("cannot uninstall hooks")
		},
		Println: func(a ...any) (int, error) {
			return 0, nil
		},
		Print: func(a ...any) (int, error) {
			return 0, nil
		},
		Scanln: func(a ...any) (int, error) {
			if ptr, ok := a[0].(*string); ok {
				*ptr = "y"
			}
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := teardown([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot uninstall hooks")
}

// =========== CHECK TESTS ===========

func TestCheck_GlobalScope(t *testing.T) {
	var printedLines []string
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit":   true,
				"post-merge":    false,
				"post-checkout": true,
			}
		},
		Printf: func(format string, a ...any) (int, error) {
			printedLines = append(printedLines, fmt.Sprintf(format, a...))
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := check([]string{}, flags, deps)

	require.NoError(t, err)
	require.True(t, len(printedLines) > 0)
	require.Contains(t, printedLines[0], "global")
}

func TestCheck_RepoScope(t *testing.T) {
	var printedLines []string
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "/path/to/repo", nil
		},
		RepoHooksPath: func(root string) (string, error) {
			return "/path/to/repo/.git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit": true,
			}
		},
		Printf: func(format string, a ...any) (int, error) {
			printedLines = append(printedLines, fmt.Sprintf(format, a...))
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := check([]string{}, flags, deps)

	require.NoError(t, err)
	require.True(t, len(printedLines) > 0)
	require.Contains(t, printedLines[0], "repository")
}

func TestCheck_NotInGitRepo(t *testing.T) {
	deps := Deps{
		RepoRoot: func(path string) (string, error) {
			return "", errors.New("not a git repo")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{"--repo"})
	err := check([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "git repo")
}

func TestCheck_GlobalHooksPathError(t *testing.T) {
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "", errors.New("cannot get global hooks path")
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := check([]string{}, flags, deps)

	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot get global hooks path")
}

func TestCheck_DisplaysInstalledAndMissing(t *testing.T) {
	var printedLines []string
	deps := Deps{
		GlobalHooksPath: func() (string, error) {
			return "/home/user/.config/git/hooks", nil
		},
		HooksStatus: func(path string) map[string]bool {
			return map[string]bool{
				"post-commit": true,
				"post-merge":  false,
			}
		},
		Printf: func(format string, a ...any) (int, error) {
			printedLines = append(printedLines, fmt.Sprintf(format, a...))
			return 0, nil
		},
	}

	flags := dispatchers.NewParsedFlags([]string{})
	err := check([]string{}, flags, deps)

	require.NoError(t, err)

	// Should have scope line + 2 hook status lines
	require.Len(t, printedLines, 3)

	hasInstalled := false
	hasMissing := false
	for _, line := range printedLines {
		if containsStr(line, "installed") {
			hasInstalled = true
		}
		if containsStr(line, "missing") {
			hasMissing = true
		}
	}
	require.True(t, hasInstalled, "should show installed hooks")
	require.True(t, hasMissing, "should show missing hooks")
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstr(s, substr)))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
