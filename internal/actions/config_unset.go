package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/config"
	"github.com/Skryensya/footprint/internal/usage"
)

func ConfigUnset(args []string, flags []string) error {
	if hasFlag(flags, "--all") {
		if len(args) > 0 {
			return usage.InvalidFlag("--all does not take arguments")
		}

		if err := config.WriteLines([]string{}); err != nil {
			return err
		}

		fmt.Println("all config entries removed")
		return nil
	}

	if len(args) < 1 {
		return usage.MissingArgument("key")
	}

	key := args[0]

	lines, err := config.ReadLines()
	if err != nil {
		return err
	}

	lines, removed := config.Unset(lines, key)
	if !removed {
		return usage.InvalidConfigKey(key)
	}

	if err := config.WriteLines(lines); err != nil {
		return err
	}

	fmt.Printf("unset %s\n", key)
	return nil
}
