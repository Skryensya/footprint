package usage

import (
	"fmt"
	"strings"
)

func AmbiguousRemote(remotes []string) *Error {
	var sb strings.Builder
	sb.WriteString("fp: repository has multiple remotes but no 'origin'\n\n")
	sb.WriteString("Available remotes:\n")
	for _, remote := range remotes {
		sb.WriteString(fmt.Sprintf("  - %s\n", remote))
	}
	sb.WriteString("\nSpecify which remote to use:\n")
	for _, remote := range remotes {
		sb.WriteString(fmt.Sprintf("  fp track --remote=%s\n", remote))
	}

	return &Error{
		Kind:    ErrAmbiguousRemote,
		Message: sb.String(),
	}
}
