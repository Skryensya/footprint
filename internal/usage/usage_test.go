package usage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// =========== ERROR TYPE TESTS ===========

func TestError_Error(t *testing.T) {
	err := &Error{
		Kind:    ErrUnknown,
		Message: "test error message",
	}

	require.Equal(t, "test error message", err.Error())
}

func TestError_GetExitCode(t *testing.T) {
	err := &Error{
		Kind:    ErrInvalidFlag,
		Message: "test",
	}

	require.Equal(t, 2, err.GetExitCode())
}

func TestError_GetExitCode_ExplicitOverride(t *testing.T) {
	err := &Error{
		Kind:     ErrUnknown,
		Message:  "test",
		ExitCode: 42,
	}

	require.Equal(t, 42, err.GetExitCode())
}

func TestError_GetExitCode_UnknownKind(t *testing.T) {
	err := &Error{
		Kind:    ErrorKind(999), // Unknown kind
		Message: "test",
	}

	require.Equal(t, 1, err.GetExitCode())
}

// =========== UNKNOWN COMMAND TESTS ===========

func TestUnknownCommand(t *testing.T) {
	err := UnknownCommand("foobar")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "foobar")
	require.Contains(t, err.Message, "not a fp command")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrUnknownCommand, err.Kind)
}

// =========== MISSING ARGUMENT TESTS ===========

func TestMissingArgument(t *testing.T) {
	err := MissingArgument("path")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "path")
	require.Contains(t, err.Message, "missing required argument")
	require.Equal(t, 2, err.GetExitCode())
	require.Equal(t, ErrMissingArgument, err.Kind)
}

// =========== INVALID FLAG TESTS ===========

func TestInvalidFlag(t *testing.T) {
	err := InvalidFlag("--unknown")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "--unknown")
	require.Contains(t, err.Message, "invalid flag")
	require.Equal(t, 2, err.GetExitCode())
	require.Equal(t, ErrInvalidFlag, err.Kind)
}

// =========== INVALID CONFIG KEY TESTS ===========

func TestInvalidConfigKey(t *testing.T) {
	err := InvalidConfigKey("nonexistent_key")

	require.NotNil(t, err)
	require.Contains(t, err.Message, "nonexistent_key")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrInvalidConfigKey, err.Kind)
}

// =========== NOT IN GIT REPO TESTS ===========

func TestNotInGitRepo(t *testing.T) {
	err := NotInGitRepo()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "git repository")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrNotInGitRepo, err.Kind)
}

// =========== INVALID PATH TESTS ===========

func TestInvalidPath(t *testing.T) {
	err := InvalidPath()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "path")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrInvalidPath, err.Kind)
}

// =========== INVALID REPO TESTS ===========

func TestInvalidRepo(t *testing.T) {
	err := InvalidRepo()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "repository")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrInvalidRepo, err.Kind)
}

// =========== GIT NOT INSTALLED TESTS ===========

func TestGitNotInstalled(t *testing.T) {
	err := GitNotInstalled()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "git")
	require.Contains(t, err.Message, "not found")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrGitNotInstalled, err.Kind)
}

// =========== FAILED CONFIG PATH TESTS ===========

func TestFailedConfigPath(t *testing.T) {
	err := FailedConfigPath()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "config")
	require.Contains(t, err.Message, "could not")
	require.Equal(t, 1, err.GetExitCode())
	require.Equal(t, ErrFailedConfigPath, err.Kind)
}

// =========== MISSING REMOTE TESTS ===========

func TestMissingRemote(t *testing.T) {
	err := MissingRemote()

	require.NotNil(t, err)
	require.Contains(t, err.Message, "remote")
	require.Equal(t, 2, err.GetExitCode())
	require.Equal(t, ErrMissingRemote, err.Kind)
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
	require.Equal(t, 2, err.GetExitCode())
	require.Equal(t, ErrAmbiguousRemote, err.Kind)
}

func TestAmbiguousRemote_SingleRemote(t *testing.T) {
	remotes := []string{"custom"}
	err := AmbiguousRemote(remotes)

	require.NotNil(t, err)
	require.Contains(t, err.Message, "custom")
	require.Contains(t, err.Message, "--remote=custom")
	require.Equal(t, ErrAmbiguousRemote, err.Kind)
}

// =========== ERROR INTERFACE COMPLIANCE ===========

func TestError_ImplementsError(t *testing.T) {
	var err error = UnknownCommand("test")
	require.NotNil(t, err)
	require.NotEmpty(t, err.Error())
}
