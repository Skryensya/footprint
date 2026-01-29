package hooks

import "os"

// File permission constants for hooks and directories.
const (
	// dirPermPrivate is for directories containing sensitive data (owner only).
	dirPermPrivate os.FileMode = 0700
	// filePemExecutable is for hook scripts that need to be executable.
	filePermExecutable os.FileMode = 0755
)

var ManagedHooks = []string{
	"post-commit",
	"post-merge",
	"post-checkout",
	"post-rewrite",
	"pre-push",
}
