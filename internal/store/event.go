package store

import "time"

type RepoEvent struct {
	ID            int64
	RepoID        string
	RepoPath      string
	Commit        string
	CommitMessage string
	Branch        string
	Timestamp     time.Time
	Status        Status
	Source        Source
}
