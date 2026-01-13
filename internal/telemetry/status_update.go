package telemetry

import "database/sql"

func UpdateStatus(
	db *sql.DB,
	repoID string,
	commit string,
	status Status,
) error {
	_, err := db.Exec(
		`UPDATE commit_events
		 SET status = ?
		 WHERE repo_id = ? AND commit = ?`,
		int(status),
		repoID,
		commit,
	)
	return err
}
