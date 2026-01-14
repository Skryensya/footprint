package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/hooks"
	"github.com/Skryensya/footprint/internal/usage"
)

func HooksInstall(_ []string, flags []string) error {
	force := hasFlag(flags, "--force")
	global := hasFlag(flags, "--global")
	repo := hasFlag(flags, "--repo")

	if global == repo {
		return usage.InvalidFlag("--repo | --global")
	}

	var (
		hooksPath string
		err       error
	)

	if repo {
		root, err := git.RepoRoot(".")
		if err != nil {
			return usage.NotInGitRepo()
		}
		hooksPath, err = git.RepoHooksPath(root)
	} else {
		hooksPath, err = git.GlobalHooksPath()
	}

	if err != nil {
		return err
	}

	needsConfirm := false
	status := hooks.Status(hooksPath)

	for _, installed := range status {
		if installed {
			needsConfirm = true
			break
		}
	}

	if needsConfirm && !force {
		fmt.Println("fp detected existing git hooks")
		fmt.Println("they will be backed up and replaced")
		fmt.Print("continue? [y/N]: ")

		var resp string
		fmt.Scanln(&resp)
		if resp != "y" && resp != "yes" {
			return nil
		}
	}

	if err := hooks.Install(hooksPath); err != nil {
		return err
	}

	fmt.Println("hooks installed")
	return nil
}
