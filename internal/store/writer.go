package store

import (
	"database/sql"
	"time"
)

func InsertEvent(db *sql.DB, e RepoEvent) error {
	_, err := db.Exec(
		`INSERT INTO repo_events
		 (repo_id, repo_path, commit_hash, commit_message, branch, author, timestamp, status_id, source_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(repo_id, commit_hash, source_id)
		 DO UPDATE SET
		   timestamp = excluded.timestamp,
		   commit_message = excluded.commit_message,
		   author = excluded.author`,
		e.RepoID,
		e.RepoPath,
		e.Commit,
		e.CommitMessage,
		e.Branch,
		e.Author,
		e.Timestamp.Format(time.RFC3339),
		int(e.Status),
		int(e.Source),
	)
	return err
}
