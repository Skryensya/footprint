package domain

import (
	"fmt"
	"strings"
)

// ErrorKind represents the type of usage error.
type ErrorKind int

const (
	ErrUnknown ErrorKind = iota
	ErrInvalidFlag
	ErrMissingArgument
	ErrUnknownCommand
	ErrAmbiguousCommand
	ErrNotInGitRepo
	ErrInvalidRepo
	ErrInvalidPath
	ErrMissingRemote
	ErrAmbiguousRemote
	ErrGitNotInstalled
	ErrInvalidConfigKey
	ErrFailedConfigPath
)

// exitCodes maps error kinds to their exit codes.
var exitCodes = map[ErrorKind]int{
	ErrUnknown:          1,
	ErrInvalidFlag:      2,
	ErrMissingArgument:  2,
	ErrUnknownCommand:   1,
	ErrAmbiguousCommand: 1,
	ErrNotInGitRepo:     1,
	ErrInvalidRepo:      1,
	ErrInvalidPath:      1,
	ErrMissingRemote:    1,
	ErrAmbiguousRemote:  1,
	ErrGitNotInstalled:  1,
	ErrInvalidConfigKey: 1,
	ErrFailedConfigPath: 1,
}

// UsageError represents a user-facing error with semantic type information.
type UsageError struct {
	Kind    ErrorKind
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *UsageError) Error() string {
	return e.Message
}

// Unwrap returns the underlying cause for error wrapping.
func (e *UsageError) Unwrap() error {
	return e.Cause
}

// ExitCode returns the appropriate exit code for this error kind.
func (e *UsageError) ExitCode() int {
	if code, ok := exitCodes[e.Kind]; ok {
		return code
	}
	return 1
}

// IsKind checks if the error is of a specific kind.
func (e *UsageError) IsKind(kind ErrorKind) bool {
	return e.Kind == kind
}

// NewUsageError creates a new usage error.
func NewUsageError(kind ErrorKind, message string) *UsageError {
	return &UsageError{
		Kind:    kind,
		Message: message,
	}
}

// NewUsageErrorf creates a new usage error with formatted message.
func NewUsageErrorf(kind ErrorKind, format string, args ...any) *UsageError {
	return &UsageError{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}
}

// WrapUsageError wraps an existing error with usage context.
func WrapUsageError(kind ErrorKind, message string, cause error) *UsageError {
	return &UsageError{
		Kind:    kind,
		Message: message,
		Cause:   cause,
	}
}

// Common error constructors for convenience.

func ErrInvalidFlagError(flag string) *UsageError {
	return NewUsageErrorf(ErrInvalidFlag, "invalid flag: %s\nUse 'fp help' for usage information.", flag)
}

func ErrMissingArgumentError(arg string) *UsageError {
	return NewUsageErrorf(ErrMissingArgument, "missing required argument: %s", arg)
}

func ErrUnknownCommandError(cmd string) *UsageError {
	return NewUsageErrorf(ErrUnknownCommand, "unknown command: %s\nRun 'fp help' for a list of commands.", cmd)
}

func ErrNotInGitRepoError() *UsageError {
	return NewUsageError(ErrNotInGitRepo, "not in a git repository (or any parent)")
}

func ErrInvalidRepoError() *UsageError {
	return NewUsageError(ErrInvalidRepo, "could not determine repository identity")
}

func ErrInvalidPathError() *UsageError {
	return NewUsageError(ErrInvalidPath, "path does not exist or is not accessible")
}

func ErrMissingRemoteError() *UsageError {
	return NewUsageError(ErrMissingRemote, "repository has no remote configured\nAdd a remote with: git remote add origin <url>")
}

func ErrAmbiguousRemoteError(remotes []string) *UsageError {
	var b strings.Builder
	b.WriteString("multiple remotes found, please specify one:\n")
	for _, r := range remotes {
		b.WriteString("  ")
		b.WriteString(r)
		b.WriteString("\n")
	}
	return NewUsageError(ErrAmbiguousRemote, b.String())
}

func ErrGitNotInstalledError() *UsageError {
	return NewUsageError(ErrGitNotInstalled, "git is not installed or not in PATH")
}

func ErrInvalidConfigKeyError(key string) *UsageError {
	return NewUsageErrorf(ErrInvalidConfigKey, "unknown configuration key: %s", key)
}

func ErrFailedConfigPathError() *UsageError {
	return NewUsageError(ErrFailedConfigPath, "could not determine config file path")
}
