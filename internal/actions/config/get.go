package config

import (
	"github.com/Skryensya/footprint/internal/usage"
)

func Get(args []string, flags []string) error {
	return get(args, flags, DefaultDeps())
}

func get(args []string, _ []string, deps Deps) error {
	if len(args) < 1 {
		return usage.MissingArgument("key")
	}

	key := args[0]

	value, found := deps.Get(key)
	if !found {
		return usage.InvalidConfigKey(key)
	}

	deps.Println(value)
	return nil
}
