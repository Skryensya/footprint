package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/repo"
	"github.com/Skryensya/footprint/internal/usage"
)

func RepoTrack(args []string, flags []string) error {
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

	added, err := repo.Track(id)
	if err != nil {
		return err
	}

	if !added {
		fmt.Printf("already tracking %s\n", id)
		return nil
	}

	fmt.Printf("tracking %s\n", id)
	return nil
}
