package repo

import (
	"errors"
	"net/url"
	"strings"

	"github.com/footprint-tools/footprint-cli/internal/config"
)

type RepoID string

const trackedReposKey = "trackedRepos"

// decodeRepoID decodes a URL-encoded repo ID.
// Used for backward compatibility with old comma-separated format.
func decodeRepoID(encoded string) RepoID {
	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		return RepoID(encoded)
	}
	return RepoID(decoded)
}

// containsPathTraversal checks if a string contains path traversal sequences
func containsPathTraversal(s string) bool {
	// Check for common path traversal patterns
	if strings.Contains(s, "..") {
		return true
	}
	// Check for null bytes which could be used to bypass checks
	if strings.Contains(s, "\x00") {
		return true
	}
	return false
}

func DeriveID(remoteURL, repoRoot string) (RepoID, error) {
	remoteURL = strings.TrimSpace(remoteURL)
	repoRoot = strings.TrimSpace(repoRoot)

	if remoteURL != "" {
		remoteURL = strings.TrimSuffix(remoteURL, ".git")

		if strings.HasPrefix(remoteURL, "git@") {
			parts := strings.SplitN(remoteURL, ":", 2)
			if len(parts) != 2 {
				return "", errors.New("invalid ssh remote url")
			}
			host := strings.TrimPrefix(parts[0], "git@")
			path := parts[1]

			// Validate against path traversal
			if containsPathTraversal(host) || containsPathTraversal(path) {
				return "", errors.New("invalid remote url: contains path traversal sequence")
			}

			// Normalize remote URLs to lowercase to prevent duplicates
			return RepoID(strings.ToLower(host + "/" + path)), nil
		}

		if strings.HasPrefix(remoteURL, "https://") || strings.HasPrefix(remoteURL, "http://") {
			remoteURL = strings.TrimPrefix(remoteURL, "https://")
			remoteURL = strings.TrimPrefix(remoteURL, "http://")

			// Validate against path traversal
			if containsPathTraversal(remoteURL) {
				return "", errors.New("invalid remote url: contains path traversal sequence")
			}

			// Normalize remote URLs to lowercase to prevent duplicates
			return RepoID(strings.ToLower(remoteURL)), nil
		}

		// Support git:// protocol (read-only git protocol)
		if strings.HasPrefix(remoteURL, "git://") {
			remoteURL = strings.TrimPrefix(remoteURL, "git://")

			// Validate against path traversal
			if containsPathTraversal(remoteURL) {
				return "", errors.New("invalid remote url: contains path traversal sequence")
			}

			// Normalize remote URLs to lowercase to prevent duplicates
			return RepoID(strings.ToLower(remoteURL)), nil
		}

		// Support file:// protocol (local repositories)
		if strings.HasPrefix(remoteURL, "file://") {
			path := strings.TrimPrefix(remoteURL, "file://")
			return RepoID("local:" + path), nil
		}

		return "", errors.New("unsupported remote url format: only git@, https://, http://, git://, and file:// are supported")
	}

	if repoRoot != "" {
		clean := strings.TrimRight(repoRoot, "/")
		return RepoID("local:" + clean), nil
	}

	return "", errors.New("cannot derive repo id")
}

func ListTracked() ([]RepoID, error) {
	lines, err := config.ReadLines()
	if err != nil {
		return nil, err
	}

	// Try new array format first (trackedRepos[]=value)
	values := config.ParseArray(lines, trackedReposKey)
	if len(values) > 0 {
		out := make([]RepoID, 0, len(values))
		for _, v := range values {
			out = append(out, RepoID(v))
		}
		return out, nil
	}

	// Fall back to old comma-separated format (trackedRepos=a,b,c)
	cfg, err := config.Parse(lines)
	if err != nil {
		return nil, err
	}

	value, ok := cfg[trackedReposKey]
	if !ok || strings.TrimSpace(value) == "" {
		return []RepoID{}, nil
	}

	parts := strings.Split(value, ",")
	out := make([]RepoID, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, decodeRepoID(p))
		}
	}

	return out, nil
}

func Track(id RepoID) (bool, error) {
	var added bool
	err := config.WithLock(func() error {
		lines, err := config.ReadLines()
		if err != nil {
			return err
		}

		// Migrate from old format if needed
		lines, err = migrateToArrayFormat(lines)
		if err != nil {
			return err
		}

		// Add to array (AppendArray checks for duplicates)
		var wasAdded bool
		lines, wasAdded = config.AppendArray(lines, trackedReposKey, string(id))
		added = wasAdded

		if !wasAdded {
			return nil // Already tracked
		}

		return config.WriteLines(lines)
	})
	return added, err
}

func Untrack(id RepoID) (bool, error) {
	var removed bool
	err := config.WithLock(func() error {
		lines, err := config.ReadLines()
		if err != nil {
			return err
		}

		// Migrate from old format if needed
		lines, err = migrateToArrayFormat(lines)
		if err != nil {
			return err
		}

		// Remove from array
		var wasRemoved bool
		lines, wasRemoved = config.RemoveFromArray(lines, trackedReposKey, string(id))
		removed = wasRemoved

		if !wasRemoved {
			return nil // Wasn't tracked
		}

		return config.WriteLines(lines)
	})
	return removed, err
}

// migrateToArrayFormat converts old comma-separated format to new array format.
// If already in array format or no tracked repos, returns lines unchanged.
func migrateToArrayFormat(lines []string) ([]string, error) {
	// Check if already using array format
	values := config.ParseArray(lines, trackedReposKey)
	if len(values) > 0 {
		return lines, nil // Already migrated
	}

	// Check for old format
	cfg, err := config.Parse(lines)
	if err != nil {
		return nil, err
	}

	value, ok := cfg[trackedReposKey]
	if !ok || strings.TrimSpace(value) == "" {
		return lines, nil // Nothing to migrate
	}

	// Parse old format
	parts := strings.Split(value, ",")
	var repos []RepoID
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			repos = append(repos, decodeRepoID(p))
		}
	}

	// Remove old format
	lines, _ = config.Unset(lines, trackedReposKey)

	// Add in new format
	for _, repo := range repos {
		lines, _ = config.AppendArray(lines, trackedReposKey, string(repo))
	}

	return lines, nil
}

func IsTracked(id RepoID) (bool, error) {
	current, err := ListTracked()
	if err != nil {
		return false, err
	}

	for _, existing := range current {
		if existing == id {
			return true, nil
		}
	}

	return false, nil
}

// ToFilesystemSafe converts a RepoID to a filesystem-safe directory name.
// Transforms:
//   - "github.com/user/repo" -> "github.com__user__repo"
//   - "local:/path/to/repo" -> "local__path__to__repo"
// The transformation is deterministic and reversible (for display).
func (id RepoID) ToFilesystemSafe() string {
	idString := string(id)

	// Replace colon (from local: prefix) with double underscore
	idString = strings.ReplaceAll(idString, ":", "__")

	// Replace path separators with double underscores
	idString = strings.ReplaceAll(idString, "/", "__")

	// Remove leading underscores
	for len(idString) > 0 && idString[0] == '_' {
		idString = idString[1:]
	}

	return idString
}
