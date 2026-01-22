package setup

import (
	"strings"

	"github.com/footprint-tools/footprint-cli/internal/dispatchers"
	"github.com/footprint-tools/footprint-cli/internal/hooks"
	"github.com/footprint-tools/footprint-cli/internal/usage"
)

func Setup(args []string, flags *dispatchers.ParsedFlags) error {
	return setup(args, flags, DefaultDeps())
}

func setup(_ []string, flags *dispatchers.ParsedFlags, deps Deps) error {
	force := flags.Has("--force")
	repo := flags.Has("--repo")

	var (
		hooksPath string
		err       error
	)

	// Default to global behavior unless --repo is explicitly passed
	if repo {
		root, err := deps.RepoRoot(".")
		if err != nil {
			return usage.NotInGitRepo()
		}
		hooksPath, err = deps.RepoHooksPath(root)
		if err != nil {
			return err
		}
	} else {
		hooksPath, err = deps.GlobalHooksPath()
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	// Check existing hooks before install
	statusBefore := deps.HooksStatus(hooksPath)
	backedUp := 0
	for _, installed := range statusBefore {
		if installed {
			backedUp++
		}
	}

	if backedUp > 0 && !force {
		deps.Println("fp detected existing git hooks")
		deps.Println("they will be backed up and replaced")
		deps.Print("continue? [y/N]: ")

		var resp string
		deps.Scanln(&resp)
		if resp != "y" && resp != "yes" {
			return nil
		}
	}

	if err := deps.HooksInstall(hooksPath); err != nil {
		return err
	}

	// Output summary
	location := "global"
	if repo {
		location = "repository"
	}

	if backedUp > 0 {
		deps.Printf("Installed %d %s hooks (%d backed up)\n", len(hooks.ManagedHooks), location, backedUp)
	} else {
		deps.Printf("Installed %d %s hooks\n", len(hooks.ManagedHooks), location)
	}
	deps.Printf("  %s\n", strings.Join(hooks.ManagedHooks, ", "))

	deps.Println("")
	deps.Println("Run 'fp track' in a repo to start recording activity")

	return nil
}
