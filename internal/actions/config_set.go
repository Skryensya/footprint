package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/config"
	"github.com/Skryensya/footprint/internal/usage"
)

func ConfigSet(args []string, _ []string) error {
	if len(args) < 2 {
		return usage.MissingArgument("key value")
	}

	key := args[0]
	value := args[1]

	lines, err := config.ReadLines()
	if err != nil {
		return err
	}

	lines, updated := config.Set(lines, key, value)

	if err := config.WriteLines(lines); err != nil {
		return err
	}

	action := "added"
	if updated {
		action = "updated"
	}

	fmt.Printf("%s %s=%s\n", action, key, value)

	return nil
}
