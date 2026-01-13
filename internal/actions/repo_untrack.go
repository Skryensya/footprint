package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/repo"
	"github.com/Skryensya/footprint/internal/usage"
)

func RepoUntrack(args []string, flags []string) error {
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

	remoteURL, _ := git.OriginURL(repoRoot)

	id, err := repo.DeriveID(remoteURL, repoRoot)
	if err != nil {
		return usage.InvalidRepo()
	}

	removed, err := repo.Untrack(id)
	if err != nil {
		return err
	}

	if !removed {
		fmt.Printf("repository not tracked: %s\n", id)
		return nil
	}

	fmt.Printf("untracked %s\n", id)
	return nil
}
