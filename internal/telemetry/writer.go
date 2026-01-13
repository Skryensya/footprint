package telemetry

import (
	"database/sql"
	"time"
)

func InsertCommit(db *sql.DB, e CommitEvent) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO commit_events
		 (repo_id, repo_path, commit, branch, timestamp, status)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		e.RepoID,
		e.RepoPath,
		e.Commit,
		e.Branch,
		e.Timestamp.Format(time.RFC3339),
		int(e.Status),
	)
	return err
}
