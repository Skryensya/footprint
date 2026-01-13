package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/repo"
	"github.com/Skryensya/footprint/internal/usage"
)

func RepoAdoptRemote(args []string, flags []string) error {
	if !git.IsAvailable() {
		return usage.GitNotInstalled()
	}

	path, err := resolvePath(args)
	if err != nil {
		return usage.InvalidPath()
	}

	repoRoot, err := git.RepoRoot(path)
	if err != nil {
		return usage.NotInGitRepo()
	}

	remoteURL, err := git.OriginURL(repoRoot)
	if err != nil || remoteURL == "" {
		return usage.MissingRemote()
	}

	localID, err := repo.DeriveID("", repoRoot)
	if err != nil {
		return usage.InvalidRepo()
	}

	remoteID, err := repo.DeriveID(remoteURL, repoRoot)
	if err != nil {
		return usage.InvalidRepo()
	}

	isLocalTracked, err := repo.IsTracked(localID)
	if err != nil {
		return err
	}

	if !isLocalTracked {
		return usage.InvalidRepo()
	}

	if _, err := repo.Untrack(localID); err != nil {
		return err
	}

	if _, err := repo.Track(remoteID); err != nil {
		return err
	}

	fmt.Printf("adopted identity:\n  %s\n→ %s\n", localID, remoteID)
	return nil
}
