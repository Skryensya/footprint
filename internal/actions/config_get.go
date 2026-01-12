package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/config"
	"github.com/Skryensya/footprint/internal/usage"
)

func ConfigGet(args []string, _ []string) error {
	if len(args) < 1 {
		return usage.MissingArgument("key")
	}

	key := args[0]

	lines, err := config.ReadLines()
	if err != nil {
		return err
	}

	configMap, err := config.Parse(lines)
	if err != nil {
		return err
	}

	value, found := configMap[key]

	if !found {
		return usage.InvalidConfigKey(key)
	}

	fmt.Println(value)
	return nil
}
