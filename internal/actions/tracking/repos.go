package tracking

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/footprint-tools/cli/internal/dispatchers"
	"github.com/footprint-tools/cli/internal/output"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ReposList lists repositories with recorded activity.
func ReposList(args []string, flags *dispatchers.ParsedFlags) error {
	return reposList(args, flags, DefaultDeps())
}

func reposList(_ []string, flags *dispatchers.ParsedFlags, deps Deps) error {
	jsonOutput := flags.Has("--json")

	s, err := deps.OpenStore(deps.DBPath())
	if err != nil {
		return err
	}
	defer func() { _ = s.Close() }()

	repos, err := s.ListRepos()
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		if jsonOutput {
			output.JSONEmpty(deps.Println)
		} else {
			_, _ = deps.Println("no tracked repositories")
			_, _ = deps.Println("run 'fp setup' in a repo to install hooks")
		}
		return nil
	}

	if jsonOutput {
		type repoJSON struct {
			Path     string `json:"path"`
			AddedAt  string `json:"added_at,omitempty"`
			LastSeen string `json:"last_seen,omitempty"`
		}
		out := make([]repoJSON, 0, len(repos))
		for _, r := range repos {
			out = append(out, repoJSON{Path: r.Path, AddedAt: r.AddedAt, LastSeen: r.LastSeen})
		}
		return output.JSON(deps.Println, out)
	}

	for _, r := range repos {
		_, _ = deps.Println(r.Path)
	}

	return nil
}

// ReposScan scans directories for git repositories and shows their hook status.
func ReposScan(args []string, flags *dispatchers.ParsedFlags) error {
	return reposScan(args, flags, DefaultDeps())
}

func reposScan(_ []string, flags *dispatchers.ParsedFlags, deps Deps) error {
	jsonOutput := flags.Has("--json")
	root := flags.String("--root", ".")

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("invalid path %s: %w", root, err)
	}
	root = absRoot

	maxDepth := flags.Int("--depth", 25)

	if !jsonOutput {
		_, _ = deps.Printf("Scanning for git repositories in %s...\n", root)
	}
	repos, err := scanForRepos(root, maxDepth)
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		if jsonOutput {
			output.JSONEmpty(deps.Println)
		} else {
			_, _ = deps.Println("No git repositories found")
		}
		return nil
	}

	if jsonOutput {
		return reposScanJSON(repos, deps)
	}

	_, _ = deps.Printf("Found %d repositories\n\n", len(repos))

	// Get home for path shortening
	home, _ := os.UserHomeDir()

	// Print repos with status
	for _, repo := range repos {
		displayPath := repo.Path
		if home != "" {
			if rel, err := filepath.Rel(home, repo.Path); err == nil && !strings.HasPrefix(rel, "..") {
				displayPath = "~/" + rel
			}
		}

		var status string
		switch {
		case repo.HasHooks:
			status = style.Success("[✓]")
		case repo.Inspection.Status.CanInstall():
			status = style.Muted("[ ]")
		default:
			status = style.Error("[×]") + " " + style.Warning(repo.Inspection.Status.String())
		}

		_, _ = deps.Printf("%s %s\n", status, displayPath)
	}

	// Summary
	_, _ = deps.Println()
	installed := 0
	canInstall := 0
	blocked := 0
	for _, r := range repos {
		switch {
		case r.HasHooks:
			installed++
		case r.Inspection.Status.CanInstall():
			canInstall++
		default:
			blocked++
		}
	}

	_, _ = deps.Printf("Installed: %d, Available: %d", installed, canInstall)
	if blocked > 0 {
		_, _ = deps.Printf(", Blocked: %d", blocked)
	}
	_, _ = deps.Println()

	if canInstall > 0 {
		_, _ = deps.Printf("\nUse 'fp repos -i' to install hooks interactively, or 'fp setup <path>' for individual repos.\n")
	}

	return nil
}

func reposScanJSON(repos []RepoEntry, deps Deps) error {
	type repoJSON struct {
		Path       string `json:"path"`
		Name       string `json:"name"`
		HasHooks   bool   `json:"has_hooks"`
		CanInstall bool   `json:"can_install"`
		Status     string `json:"status,omitempty"`
	}

	out := make([]repoJSON, 0, len(repos))
	for _, r := range repos {
		entry := repoJSON{
			Path:       r.Path,
			Name:       r.Name,
			HasHooks:   r.HasHooks,
			CanInstall: r.Inspection.Status.CanInstall(),
		}
		if !r.HasHooks && !r.Inspection.Status.CanInstall() {
			entry.Status = r.Inspection.Status.String()
		}
		out = append(out, entry)
	}

	return output.JSON(deps.Println, out)
}
