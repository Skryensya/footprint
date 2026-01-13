package cli

import (
	"github.com/Skryensya/footprint/internal/actions"
	"github.com/Skryensya/footprint/internal/dispatchers"
)

func BuildTree() *dispatchers.DispatchNode {
	root := dispatchers.NewNode(
		"fp",
		nil,
		"Track your work across repositories",
		"fp <command> [flags]",
		[]dispatchers.FlagDescriptor{
			{
				Names:       []string{"--help", "-h"},
				Description: "Show help",
				Scope:       dispatchers.FlagScopeGlobal,
			},
			{
				Names:       []string{"--version", "-v"},
				Description: "Show version",
				Scope:       dispatchers.FlagScopeGlobal,
			},
		},
		nil,
		nil,
	)

	dispatchers.NewNode(
		"version",
		root,
		"Show current fp version",
		"fp version",
		nil,
		nil,
		actions.ShowVersion,
	).Category = dispatchers.CategoryInfo

	dispatchers.NewNode(
		"hello-world",
		root,
		"Run a basic sanity check",
		"fp hello-world",
		nil,
		nil,
		actions.HelloWorld,
	).Category = dispatchers.CategoryInfo

	config := dispatchers.NewNode(
		"config",
		root,
		"Manage configuration",
		"fp config <command>",
		nil,
		nil,
		nil,
	)

	dispatchers.NewNode(
		"get",
		config,
		"Get a config value",
		"fp config get <key>",
		nil,
		[]dispatchers.ArgSpec{
			{
				Name:        "key",
				Description: "Configuration key to read",
				Required:    true,
			},
		},
		actions.ConfigGet,
	).Category = dispatchers.CategoryConfig

	dispatchers.NewNode(
		"set",
		config,
		"Set a config value",
		"fp config set <key> <value>",
		nil,
		[]dispatchers.ArgSpec{
			{
				Name:        "key",
				Description: "Configuration key to write",
				Required:    true,
			},
			{
				Name:        "value",
				Description: "Value to assign",
				Required:    true,
			},
		},
		actions.ConfigSet,
	).Category = dispatchers.CategoryConfig

	dispatchers.NewNode(
		"unset",
		config,
		"Remove a config value",
		"fp config unset <key> <value>",
		[]dispatchers.FlagDescriptor{
			{
				Names:       []string{"--all"},
				Description: "Delete all the config key=value pairs",
				Scope:       dispatchers.FlagScopeLocal,
			},
		},
		[]dispatchers.ArgSpec{
			{
				Name:        "key",
				Description: "Configuration key to delete",
				Required:    false,
			},
		},
		actions.ConfigUnset,
	).Category = dispatchers.CategoryConfig

	dispatchers.NewNode(
		"list",
		config,
		"List all the configuration as key=value pairs",
		"fp config list",
		nil,
		nil,
		actions.ConfigList,
	).Category = dispatchers.CategoryConfig

	repo := dispatchers.NewNode(
		"repo",
		root,
		"Manage repository tracking",
		"fp repo <command>",
		nil,
		nil,
		nil,
	)

	dispatchers.NewNode(
		"track",
		repo,
		"Start tracking a repository",
		"fp repo track <path>",
		nil,
		[]dispatchers.ArgSpec{
			{
				Name:        "path",
				Description: "Path to a git repository (defaults to current directory)",
				Required:    false,
			},
		},
		actions.RepoTrack,
	).Category = dispatchers.CategoryRepo

	dispatchers.NewNode(
		"untrack",
		repo,
		"Stop tracking a repository",
		"fp repo untrack <path>",
		nil,
		[]dispatchers.ArgSpec{
			{
				Name:        "path",
				Description: "Path to a git repository (defaults to current directory)",
				Required:    false,
			},
		},
		actions.RepoUntrack,
	).Category = dispatchers.CategoryRepo

	dispatchers.NewNode(
		"list",
		repo,
		"Show all the repositories being tracked",
		"fp repo list",
		nil,
		nil,
		actions.RepoList,
	).Category = dispatchers.CategoryRepo

	dispatchers.NewNode(
		"status",
		repo,
		"Show repository tracking status",
		"fp repo status <path>",
		nil,
		[]dispatchers.ArgSpec{
			{
				Name:        "path",
				Description: "Path to a git repository (defaults to current directory)",
				Required:    false,
			},
		},
		actions.RepoStatus,
	).Category = dispatchers.CategoryRepo

	dispatchers.NewNode(
		"adopt-remote",
		repo,
		"Update repository id to reflect new remote",
		"fp repo adopt-remote <path>",
		nil,
		[]dispatchers.ArgSpec{
			{
				Name:        "path",
				Description: "Path to a git repository (defaults to current directory)",
				Required:    false,
			},
		},
		actions.RepoAdoptRemote,
	).Category = dispatchers.CategoryRepo

	dispatchers.NewNode(
		"help",
		root,
		"Show help for a command",
		"fp help [command]",
		nil,
		nil,
		nil,
	)

	return root
}
