package tracking

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/store"
)

// FormatEvent formats a single event for display.
// If oneline is true, uses compact single-line format.
func FormatEvent(e store.RepoEvent, oneline bool) string {
	if oneline {
		message := truncateMessage(e.CommitMessage, 15)
		return fmt.Sprintf(
			"%s %-9s %-13s %-20s %-8s %.7s %s",
			e.Timestamp.Format("2006-01-02 15:04"),
			e.Status.String(),
			e.Source.String(),
			e.RepoID,
			e.Branch,
			e.Commit,
			message,
		)
	}

	return fmt.Sprintf(
		"%s %-9s %-13s %-20s %-8s %.7s\n    %s",
		e.Timestamp.Format("2006-01-02 15:04"),
		e.Status.String(),
		e.Source.String(),
		e.RepoID,
		e.Branch,
		e.Commit,
		e.CommitMessage,
	)
}
