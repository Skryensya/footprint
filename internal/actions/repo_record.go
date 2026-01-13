package actions

import (
	"time"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/repo"
	"github.com/Skryensya/footprint/internal/telemetry"
)

func RepoRecord(args []string, flags []string) error {
	defer func() { _ = recover() }()

	if !git.IsAvailable() {
		return nil
	}

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		return nil
	}

	remoteURL, _ := git.OriginURL(repoRoot)

	repoID, err := repo.DeriveID(remoteURL, repoRoot)
	if err != nil {
		return nil
	}

	tracked, err := repo.IsTracked(repoID)
	if err != nil || !tracked {
		return nil
	}

	commit, err := git.HeadCommit()
	if err != nil {
		return nil
	}

	branch, _ := git.CurrentBranch()

	db, err := telemetry.Open(telemetry.DBPath())
	if err != nil {
		return nil
	}

	_ = telemetry.Init(db)

	_ = telemetry.InsertCommit(db, telemetry.CommitEvent{
		RepoID:    string(repoID),
		RepoPath:  repoRoot,
		Commit:    commit,
		Branch:    branch,
		Timestamp: time.Now().UTC(),
		Status:    telemetry.StatusPending,
	})

	return nil
}
