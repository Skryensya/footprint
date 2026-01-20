package usage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// =========== ERROR TYPE TESTS ===========

func TestError_Error(t *testing.T) {
	err := &Error{
		Message:  "test error message",
		ExitCode: 1,
	}

	require.Equal(t, "test error message", err.Error())
}

func TestError_ExitCode(t *testing.T) {
	err := &Error{
		Message:  "test",
		ExitCode: 42,
	}

	require.Equal(t, 42, err.ExitCode)
}

// =========== UNKNOWN COMMAND TESTS ===========

func TestUnknownCommand(t *testing.T) {
	err := UnknownCommand("foobar")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "foobar")
	require.Contains(t, err.Message, "not a fp command")
	require.Equal(t, 1, err.ExitCode)
}

// =========== MISSING ARGUMENT TESTS ===========

func TestMissingArgument(t *testing.T) {
	err := MissingArgument("path")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "path")
	require.Contains(t, err.Message, "missing required argument")
	require.Equal(t, 2, err.ExitCode)
}

// =========== INVALID FLAG TESTS ===========

func TestInvalidFlag(t *testing.T) {
	err := InvalidFlag("--unknown")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "--unknown")
	require.Contains(t, err.Message, "invalid flag")
	require.Equal(t, 2, err.ExitCode)
}

// =========== INVALID CONFIG KEY TESTS ===========

func TestInvalidConfigKey(t *testing.T) {
	err := InvalidConfigKey("nonexistent_key")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "nonexistent_key")
	require.Equal(t, 1, err.ExitCode)
}

// =========== NOT IN GIT REPO TESTS ===========

func TestNotInGitRepo(t *testing.T) {
	err := NotInGitRepo()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "git repository")
	require.Equal(t, 1, err.ExitCode)
}

// =========== INVALID PATH TESTS ===========

func TestInvalidPath(t *testing.T) {
	err := InvalidPath()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "Invalid path")
	require.Equal(t, 1, err.ExitCode)
}

// =========== INVALID REPO TESTS ===========

func TestInvalidRepo(t *testing.T) {
	err := InvalidRepo()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "Invalid repository")
	require.Equal(t, 1, err.ExitCode)
}

// =========== GIT NOT INSTALLED TESTS ===========

func TestGitNotInstalled(t *testing.T) {
	err := GitNotInstalled()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "Git")
	require.Contains(t, err.Message, "not installed")
	require.Equal(t, 1, err.ExitCode)
}

// =========== FAILED CONFIG PATH TESTS ===========

func TestFailedConfigPath(t *testing.T) {
	err := FailedConfigPath()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "Config")
	require.Contains(t, err.Message, "failed")
	require.Equal(t, 1, err.ExitCode)
}

// =========== MISSING REMOTE TESTS ===========

func TestMissingRemote(t *testing.T) {
	err := MissingRemote()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "remote")
	require.Equal(t, 2, err.ExitCode)
}

// =========== AMBIGUOUS REMOTE TESTS ===========

func TestAmbiguousRemote(t *testing.T) {
	remotes := []string{"upstream", "fork", "backup"}
	err := AmbiguousRemote(remotes)

	require.NotNil(t, err)
	require.Contains(t, err.Message, "upstream")
	require.Contains(t, err.Message, "fork")
	require.Contains(t, err.Message, "backup")
	require.Contains(t, err.Message, "multiple remotes")
	require.Equal(t, 2, err.ExitCode)
}

func TestAmbiguousRemote_SingleRemote(t *testing.T) {
	remotes := []string{"custom"}
	err := AmbiguousRemote(remotes)

	require.NotNil(t, err)
	require.Contains(t, err.Message, "custom")
	require.Contains(t, err.Message, "--remote=<custom>")
}

// =========== ERROR INTERFACE COMPLIANCE ===========

func TestError_ImplementsError(t *testing.T) {
	var err error = UnknownCommand("test")
	require.NotNil(t, err)
	require.NotEmpty(t, err.Error())
}
