package tracking

import (
	"github.com/Skryensya/footprint/internal/usage"
)

func Track(args []string, flags []string) error {
	return track(args, flags, DefaultDeps())
}

func track(args []string, flags []string, deps Deps) error {
	if !deps.GitIsAvailable() {
		return usage.GitNotInstalled()
	}

	path, err := resolvePath(args)
	if err != nil {
		return usage.InvalidPath()
	}

	repoRoot, err := deps.RepoRoot(path)
	if err != nil {
		return usage.NotInGitRepo()
	}

	// Determine the remote URL to use
	remoteURL, err := resolveRemoteURL(repoRoot, flags, deps)
	if err != nil {
		return err
	}

	id, err := deps.DeriveID(remoteURL, repoRoot)
	if err != nil {
		return usage.InvalidRepo()
	}

	added, err := deps.Track(id)
	if err != nil {
		return err
	}

	if !added {
		deps.Printf("already tracking %s\n", id)
		return nil
	}

	deps.Printf("tracking %s\n", id)
	return nil
}

// resolveRemoteURL determines which remote URL to use for tracking.
// Priority:
//  1. --remote flag if specified
//  2. 'origin' remote if it exists
//  3. Single remote if only one exists
//  4. Error if multiple remotes exist without 'origin'
//  5. Empty string if no remotes (will use local path)
func resolveRemoteURL(repoRoot string, flags []string, deps Deps) (string, error) {
	// Check for --remote flag
	specifiedRemote := getFlagValue(flags, "--remote")
	if specifiedRemote != "" {
		url, err := deps.GetRemoteURL(repoRoot, specifiedRemote)
		if err != nil {
			return "", usage.MissingRemote()
		}
		return url, nil
	}

	// Try origin first
	if url, err := deps.OriginURL(repoRoot); err == nil && url != "" {
		return url, nil
	}

	// No origin, check available remotes
	remotes, err := deps.ListRemotes(repoRoot)
	if err != nil || len(remotes) == 0 {
		// No remotes at all, will use local path
		return "", nil
	}

	if len(remotes) == 1 {
		// Exactly one remote, use it
		url, err := deps.GetRemoteURL(repoRoot, remotes[0])
		if err != nil {
			return "", nil
		}
		return url, nil
	}

	// Multiple remotes without origin - ambiguous
	return "", usage.AmbiguousRemote(remotes)
}

