package telemetry

import "time"

type CommitEvent struct {
	RepoID    string
	RepoPath  string
	Commit    string
	Branch    string
	Timestamp time.Time
	Status    Status
}
