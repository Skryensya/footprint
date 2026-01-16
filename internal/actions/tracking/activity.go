package tracking

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/Skryensya/footprint/internal/store"
)

func Activity(args []string, flags []string) error {
	return activity(args, flags, DefaultDeps())
}

func activity(_ []string, flags []string, deps Deps) error {
	db, err := deps.OpenDB(deps.DBPath())
	if err != nil {
		return nil
	}

	var (
		filter  store.EventFilter
		oneline bool
	)

	for _, flag := range flags {
		if flag == "--oneline" {
			oneline = true
			continue
		}

		if strings.HasPrefix(flag, "--status=") {
			if status, ok := parseStatus(strings.TrimPrefix(flag, "--status=")); ok {
				filter.Status = &status
			}
			continue
		}

		if strings.HasPrefix(flag, "--source=") {
			if source, ok := parseSource(strings.TrimPrefix(flag, "--source=")); ok {
				filter.Source = &source
			}
			continue
		}

		if strings.HasPrefix(flag, "--since=") {
			if t, ok := parseDate(strings.TrimPrefix(flag, "--since=")); ok {
				filter.Since = &t
			}
			continue
		}

		if strings.HasPrefix(flag, "--until=") {
			if t, ok := parseDate(strings.TrimPrefix(flag, "--until=")); ok {
				filter.Until = &t
			}
			continue
		}

		if strings.HasPrefix(flag, "--repo=") {
			repoID := strings.TrimPrefix(flag, "--repo=")
			filter.RepoID = &repoID
			continue
		}

		if strings.HasPrefix(flag, "--limit=") {
			if n, err := strconv.Atoi(strings.TrimPrefix(flag, "--limit=")); err == nil && n > 0 {
				filter.Limit = n
			}
			continue
		}
	}

	events, err := deps.ListEvents(db, filter)
	if err != nil || len(events) == 0 {
		return nil
	}

	var output bytes.Buffer

	for _, event := range events {
		output.WriteString(FormatEvent(event, oneline))
		output.WriteString("\n")
	}

	deps.Pager(output.String())
	return nil
}
