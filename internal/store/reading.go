package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type EventFilter struct {
	Status *Status
	Source *Source
	Since  *time.Time
	Until  *time.Time
	RepoID *string
	Limit  int
}

func ListEvents(db *sql.DB, filter EventFilter) ([]RepoEvent, error) {

	base := `
		SELECT
			id,
			repo_id,
			repo_path,
			commit_hash,
			commit_message,
			branch,
			timestamp,
			status_id,
			source_id
		FROM repo_events
	`

	var (
		clauses []string
		args    []any
	)

	if filter.Status != nil {
		clauses = append(clauses, "status_id = ?")
		args = append(args, int(*filter.Status))
	}

	if filter.Source != nil {
		clauses = append(clauses, "source_id = ?")
		args = append(args, int(*filter.Source))
	}

	if filter.Since != nil {
		clauses = append(clauses, "timestamp >= ?")
		args = append(args, filter.Since.Format(time.RFC3339))
	}

	if filter.Until != nil {
		clauses = append(clauses, "timestamp <= ?")
		args = append(args, filter.Until.Format(time.RFC3339))
	}

	if filter.RepoID != nil {
		clauses = append(clauses, "repo_id = ?")
		args = append(args, *filter.RepoID)
	}

	query := base

	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	query += " ORDER BY timestamp DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RepoEvent

	for rows.Next() {
		var (
			e        RepoEvent
			ts       string
			statusID int
			sourceID int
		)

		if err := rows.Scan(
			&e.ID,
			&e.RepoID,
			&e.RepoPath,
			&e.Commit,
			&e.CommitMessage,
			&e.Branch,
			&ts,
			&statusID,
			&sourceID,
		); err != nil {
			return nil, err
		}

		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}

		e.Timestamp = t
		e.Status = Status(statusID)
		e.Source = Source(sourceID)

		out = append(out, e)
	}

	return out, rows.Err()
}

// GetMaxEventID returns the highest event ID in the database.
// Returns 0 if no events exist.
func GetMaxEventID(db *sql.DB) (int64, error) {
	var maxID sql.NullInt64
	err := db.QueryRow("SELECT MAX(id) FROM repo_events").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	if !maxID.Valid {
		return 0, nil
	}
	return maxID.Int64, nil
}

// ListEventsSince returns events with ID greater than afterID, ordered by ID ascending.
// Used for polling new events in real-time.
func ListEventsSince(db *sql.DB, afterID int64) ([]RepoEvent, error) {
	query := `
		SELECT
			id,
			repo_id,
			repo_path,
			commit_hash,
			commit_message,
			branch,
			timestamp,
			status_id,
			source_id
		FROM repo_events
		WHERE id > ?
		ORDER BY id ASC
	`

	rows, err := db.Query(query, afterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RepoEvent

	for rows.Next() {
		var (
			e        RepoEvent
			ts       string
			statusID int
			sourceID int
		)

		if err := rows.Scan(
			&e.ID,
			&e.RepoID,
			&e.RepoPath,
			&e.Commit,
			&e.CommitMessage,
			&e.Branch,
			&ts,
			&statusID,
			&sourceID,
		); err != nil {
			return nil, err
		}

		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}

		e.Timestamp = t
		e.Status = Status(statusID)
		e.Source = Source(sourceID)

		out = append(out, e)
	}

	return out, rows.Err()
}
