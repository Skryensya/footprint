package setup

import (
	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/usage"
)

func Teardown(args []string, flags *dispatchers.ParsedFlags) error {
	return teardown(args, flags, DefaultDeps())
}

func teardown(_ []string, flags *dispatchers.ParsedFlags, deps Deps) error {
	force := flags.Has("--force")
	repo := flags.Has("--repo")

	var hooksPath string
	var err error

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

	if !force {
		deps.Println("fp will remove its git hooks")
		deps.Println("previous hooks will be restored if available")
		deps.Print("continue? [y/N]: ")

		var resp string
		deps.Scanln(&resp)
		if resp != "y" && resp != "yes" {
			return nil
		}
	}

	if err := deps.HooksUninstall(hooksPath); err != nil {
		return err
	}

	deps.Println("fp teardown complete")
	return nil
}
