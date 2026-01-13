package telemetry

import (
	"database/sql"
	"time"
)

func PendingCommits(db *sql.DB) ([]CommitEvent, error) {
	rows, err := db.Query(
		`SELECT repo_id, repo_path, commit, branch, timestamp, status
		 FROM commit_events
		 WHERE status = ?`,
		int(StatusPending),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CommitEvent

	for rows.Next() {
		var (
			e      CommitEvent
			ts     string
			status int
		)

		if err := rows.Scan(
			&e.RepoID,
			&e.RepoPath,
			&e.Commit,
			&e.Branch,
			&ts,
			&status,
		); err != nil {
			return nil, err
		}

		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}

		e.Timestamp = t
		e.Status = Status(status)

		out = append(out, e)
	}

	return out, rows.Err()
}
