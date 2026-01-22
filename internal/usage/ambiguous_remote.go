package usage

import (
	"fmt"
	"strings"
)

func AmbiguousRemote(remotes []string) *Error {
	return &Error{
		Kind: ErrAmbiguousRemote,
		Message: fmt.Sprintf(
			"fp: repository has multiple remotes but no 'origin'\nAvailable remotes: %s\nUse: fp track --remote=<%s>",
			strings.Join(remotes, ", "),
			remotes[0],
		),
	}
}
