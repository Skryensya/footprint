package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/repo"
	"github.com/Skryensya/footprint/internal/usage"
)

func RepoStatus(args []string, flags []string) error {
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

	localID, err := repo.DeriveID("", repoRoot)
	if err != nil {
		return usage.InvalidRepo()
	}

	remoteID := repo.RepoID("")
	if remoteURL != "" {
		remoteID, _ = repo.DeriveID(remoteURL, repoRoot)
	}

	isLocalTracked, err := repo.IsTracked(localID)
	if err != nil {
		return err
	}

	isRemoteTracked := false
	if remoteID != "" {
		isRemoteTracked, err = repo.IsTracked(remoteID)
		if err != nil {
			return err
		}
	}

	if isRemoteTracked {
		fmt.Printf("tracked %s\n", remoteID)
		return nil
	}

	if isLocalTracked {
		fmt.Printf("tracked %s\n", localID)

		if remoteID != "" && localID != remoteID {
			fmt.Printf("remote detected %s\n", remoteID)
			fmt.Println("run 'fp repo adopt-remote' to update identity")
		}

		return nil
	}

	if remoteID != "" {
		fmt.Printf("not tracked %s\n", remoteID)
	} else {
		fmt.Printf("not tracked %s\n", localID)
	}

	return nil
}
