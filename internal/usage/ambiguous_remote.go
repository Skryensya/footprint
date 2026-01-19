package usage

import (
	"fmt"
	"strings"
)

func AmbiguousRemote(remotes []string) *Error {
	return &Error{
		Message: fmt.Sprintf(
			"Repository has multiple remotes but no 'origin'.\nAvailable remotes: %s\nUse: fp track --remote=<%s>",
			strings.Join(remotes, ", "),
			remotes[0],
		),
		ExitCode: 2,
	}
}
