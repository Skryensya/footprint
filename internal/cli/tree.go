package cli

import (
	"github.com/Skryensya/footprint/internal/actions"
	configactions "github.com/Skryensya/footprint/internal/actions/config"
	logsactions "github.com/Skryensya/footprint/internal/actions/logs"
	setupactions "github.com/Skryensya/footprint/internal/actions/setup"
	themeactions "github.com/Skryensya/footprint/internal/actions/theme"
	trackingactions "github.com/Skryensya/footprint/internal/actions/tracking"
	"github.com/Skryensya/footprint/internal/dispatchers"
)

func BuildTree() *dispatchers.DispatchNode {
	root := dispatchers.Root(dispatchers.RootSpec{
		Name:    "fp",
		Summary: "Track and inspect Git activity across repositories",
		Usage:   "fp [--help] [--version] [--no-color] [--no-pager] <command> [<args>]",
		Flags:   RootFlags,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "version",
		Parent:  root,
		Summary: "Show current fp version",
		Description: `Prints the version number of the fp binary.

The version includes the git tag, commit count since tag, and commit hash
when built from a non-release commit.`,
		Usage:    "fp version",
		Action:   actions.ShowVersion,
		Category: dispatchers.CategoryInspectActivity,
	})

	// -- config

	config := dispatchers.Group(dispatchers.GroupSpec{
		Name:    "config",
		Parent:  root,
		Summary: "Manage configuration",
		Description: `Read and write fp configuration values.

Configuration is stored in ~/.fprc as simple key=value pairs.
Use 'fp config list' to see all current settings.`,
		Usage: "fp config <command>",
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "get",
		Parent:  config,
		Summary: "Get a config value",
		Description: `Prints the value of a configuration key.

If the key does not exist, nothing is printed and the command exits
with a non-zero status.`,
		Usage:    "fp config get <key>",
		Args:     ConfigKeyArg,
		Action:   configactions.Get,
		Category: dispatchers.CategoryConfig,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "set",
		Parent:  config,
		Summary: "Set a config value",
		Description: `Sets a configuration key to the specified value.

If the key already exists, its value is overwritten.
The configuration file is created if it does not exist.`,
		Usage:    "fp config set <key> <value>",
		Args:     ConfigKeyValueArgs,
		Action:   configactions.Set,
		Category: dispatchers.CategoryConfig,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "unset",
		Parent:  config,
		Summary: "Remove a config value",
		Description: `Removes a configuration key from the config file.

Use --all to remove all configuration values and reset to defaults.`,
		Usage: "fp config unset <key>",
		Flags: ConfigUnsetFlags,
		Args: []dispatchers.ArgSpec{
			{
				Name:        "key",
				Description: "Configuration key to delete",
				Required:    false,
			},
		},
		Action:   configactions.Unset,
		Category: dispatchers.CategoryConfig,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "list",
		Parent:  config,
		Summary: "List all the configuration as key=value pairs",
		Description: `Prints all configuration values in key=value format.

This shows the current state of ~/.fprc. If no configuration has been
set, the output will be empty.`,
		Usage:    "fp config list",
		Action:   configactions.List,
		Category: dispatchers.CategoryConfig,
	})

	// -- theme

	theme := dispatchers.Group(dispatchers.GroupSpec{
		Name:    "theme",
		Parent:  root,
		Summary: "Manage color themes",
		Description: `View and change color themes.

fp includes several built-in themes optimized for different terminal
backgrounds. Each theme has a dark and light variant.

Available themes:
  default-dark, default-light   Classic terminal colors
  ocean-dark, ocean-light       Blue/teal palette
  forest-dark, forest-light     Green/earth palette

Use 'fp theme' to see all themes and the current selection.
Use 'fp theme set <name>' to change the theme.`,
		Usage: "fp theme [command]",
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "list",
		Parent:  theme,
		Summary: "List available themes",
		Description: `Shows all available color themes.

The current theme is marked with an asterisk (*).`,
		Usage:    "fp theme list",
		Action:   themeactions.List,
		Category: dispatchers.CategoryConfig,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "set",
		Parent:  theme,
		Summary: "Set the color theme",
		Description: `Changes the current color theme.

The theme is saved to ~/.fprc and takes effect immediately for new
commands. Choose a -dark theme for dark terminal backgrounds or a
-light theme for light terminal backgrounds.`,
		Usage:    "fp theme set <name>",
		Args:     ThemeNameArg,
		Action:   themeactions.Set,
		Category: dispatchers.CategoryConfig,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "show",
		Parent:  theme,
		Summary: "Show current theme color values",
		Description: `Displays the color values for the current theme.

Shows each semantic color (success, warning, error, info, muted, header)
and source colors (1-6) with their corresponding ANSI color codes.

This is useful to see exactly what colors are being used and to
customize individual colors with 'fp config set color_<name> <value>'.`,
		Usage:    "fp theme show",
		Action:   themeactions.Show,
		Category: dispatchers.CategoryConfig,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "pick",
		Parent:  theme,
		Summary: "Interactively select a theme",
		Description: `Opens an interactive picker to browse and select themes.

Use arrow keys or j/k to navigate through available themes.
Press enter or space to select and apply the highlighted theme.
Press q or esc to cancel without making changes.

The picker shows a live preview of each theme's colors as you navigate.`,
		Usage:    "fp theme pick",
		Action:   themeactions.Pick,
		Category: dispatchers.CategoryConfig,
	})

	// -- tracking

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "track",
		Parent:  root,
		Summary: "Start tracking a repository",
		Description: `Marks a git repository for activity tracking.

When a repository is tracked, fp records git events (commits, merges,
checkouts, rebases, pushes) that occur in it. Events are stored locally
in a SQLite database.

The repository is identified by its remote URL (usually 'origin'). If no
'origin' remote exists but exactly one other remote is available, that
remote is used. If multiple remotes exist without 'origin', use --remote
to specify which one to use.

If no remote exists at all, a local path identifier is used instead.

If no path is provided, the current directory is used.`,
		Usage:    "fp track [--remote=<name>] [path]",
		Args:     OptionalRepoPathArg,
		Flags:    TrackFlags,
		Action:   trackingactions.Track,
		Category: dispatchers.CategoryGetStarted,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "untrack",
		Parent:  root,
		Summary: "Stop tracking a repository",
		Description: `Removes a repository from activity tracking.

Future git events in this repository will no longer be recorded.
Existing recorded events are not deleted.

If no path is provided, the current directory is used.

Use --id to untrack by repository ID directly. This is useful when the
repository directory no longer exists but the tracking entry remains.
Repository IDs can be found via 'fp repos'.`,
		Usage:    "fp untrack [path] [--id=<repo-id>]",
		Args:     OptionalRepoPathArg,
		Flags:    UntrackFlags,
		Action:   trackingactions.Untrack,
		Category: dispatchers.CategoryManageRepos,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "repos",
		Parent:  root,
		Summary: "Show all tracked repositories",
		Description: `Lists all repositories currently being tracked by fp.

Each entry shows the repository identifier and, when available, the
local path where it was last seen.`,
		Usage:    "fp repos",
		Action:   trackingactions.Repos,
		Category: dispatchers.CategoryInspectActivity,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "list",
		Parent:  root,
		Summary: "Show all tracked repositories (alias for repos)",
		Description: `Lists all repositories currently being tracked by fp.

This is an alias for 'fp repos'.`,
		Usage:    "fp list",
		Action:   trackingactions.Repos,
		Category: dispatchers.CategoryInspectActivity,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "status",
		Parent:  root,
		Summary: "Show repository tracking status",
		Description: `Shows whether a repository is currently being tracked.

Displays the repository identifier, tracking state, and hook installation
status for the given path.

If no path is provided, the current directory is used.`,
		Usage:    "fp status [path]",
		Args:     OptionalRepoPathArg,
		Action:   trackingactions.Status,
		Category: dispatchers.CategoryInspectActivity,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "sync-remote",
		Parent:  root,
		Summary: "Update repository id to reflect new remote",
		Description: `Updates the repository identifier after a remote URL change.

Repository IDs are derived from the remote URL. If you change the remote
(e.g., after transferring a repo to a new host), run this command to
update the tracking ID. This ensures future events are associated with
the new identifier.

If no path is provided, the current directory is used.`,
		Usage:    "fp sync-remote [path]",
		Args:     OptionalRepoPathArg,
		Action:   trackingactions.Adopt,
		Category: dispatchers.CategoryManageRepos,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "record",
		Parent:  root,
		Summary: "Record a git event (invoked automatically by git hooks)",
		Description: `Records a git event to the local database.

This command is normally invoked automatically by git hooks installed
via 'fp setup'. You typically don't need to run it manually.

When called, it checks if the current repository is tracked. If so, it
records the event. If not, it exits silently without error.`,
		Usage:    "fp record",
		Flags:    RecordFlags,
		Action:   trackingactions.Record,
		Category: dispatchers.CategoryPlumbing,
	})

	// -- activity

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "activity",
		Parent:  root,
		Summary: "List recorded activity (newest first)",
		Description: `Shows recorded git events across all tracked repositories.

Events are displayed in reverse chronological order (newest first).
Each entry shows the timestamp, event type, repository, and relevant
details like commit SHA or branch name.

Use -n to limit the number of entries shown.`,
		Usage:    "fp activity",
		Action:   trackingactions.Activity,
		Flags:    ActivityFlags,
		Category: dispatchers.CategoryInspectActivity,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "watch",
		Parent:  root,
		Summary: "Stream new events in real time",
		Description: `Watches for new git events and prints them as they occur.

This command runs continuously like 'tail -f', polling the database
for new events. Only events that occur after the command starts are
shown; historical events are not displayed.

Press Ctrl+C to stop.`,
		Usage:    "fp watch",
		Action:   trackingactions.Log,
		Flags:    WatchFlags,
		Category: dispatchers.CategoryInspectActivity,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "export",
		Parent:  root,
		Summary: "Export pending events to CSV",
		Description: `Exports all pending events to a CSV file in the export repository.

By default, exports only run if the configured interval has passed since
the last export. Use --force to export immediately regardless of interval.

The export repository is a git repository where CSV files are committed.
Default location: ~/.config/Footprint/exports

Configuration:
  export_interval  Seconds between exports (default: 3600 = 1 hour)
  export_repo      Path to export repository
  export_last      Unix timestamp of last export (managed internally)

Use 'fp config set export_interval 1800' to change the interval.`,
		Usage:    "fp export [--force] [--dry-run]",
		Action:   trackingactions.Export,
		Flags:    ExportFlags,
		Category: dispatchers.CategoryManageRepos,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "backfill",
		Parent:  root,
		Summary: "Import historical commits from a repository",
		Description: `Imports existing commits from a git repository into the database.

This allows you to retroactively track commit history that occurred before
fp was installed. The command scans the git log and inserts each commit
as a pending event.

The backfill runs in the background while 'fp watch --oneline' shows
progress in real-time. Press Ctrl+C to stop watching (the import continues).

Use --dry-run to preview what would be imported without modifying the database.

Events are inserted with source "BACKFILL" and status "pending". Run
'fp export --force' afterward to export the backfilled events.

Duplicate commits (same repo + commit hash + source) are automatically skipped.`,
		Usage:    "fp backfill [path] [--since=<date>] [--until=<date>] [--limit=<n>]",
		Args:     OptionalRepoPathArg,
		Flags:    BackfillFlags,
		Action:   trackingactions.Backfill,
		Category: dispatchers.CategoryManageRepos,
	})

	// -- setup

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "setup",
		Parent:  root,
		Summary: "Install fp git hooks",
		Description: `Installs git hooks that automatically record activity.

By default, hooks are installed in the current repository's .git/hooks/
directory. Use --global to install hooks in git's global hooks directory,
which applies to all repositories.

Installed hooks:
  post-commit     Records commits
  post-merge      Records merges (including pulls)
  post-checkout   Records branch switches
  post-rewrite    Records rebases and amends
  pre-push        Records push attempts

If existing hooks are found, they are backed up before installation.`,
		Usage:    "fp setup [--global]",
		Flags:    SetupFlags,
		Action:   setupactions.Setup,
		Category: dispatchers.CategoryGetStarted,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "teardown",
		Parent:  root,
		Summary: "Remove fp git hooks",
		Description: `Removes git hooks installed by fp.

By default, removes hooks from the current repository. Use --global to
remove hooks from git's global hooks directory.

If hooks were backed up during installation, the original hooks are
restored.`,
		Usage:    "fp teardown [--global]",
		Flags:    TeardownFlags,
		Action:   setupactions.Teardown,
		Category: dispatchers.CategoryManageRepos,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "check",
		Parent:  root,
		Summary: "Show installed fp hooks status",
		Description: `Shows the installation status of fp git hooks.

Displays which hooks are installed, whether they are fp hooks or
third-party hooks, and if any backups exist.

Use --global to check the global hooks directory instead of the
current repository.`,
		Usage:    "fp check [--global]",
		Flags:    CheckFlags,
		Action:   setupactions.Check,
		Category: dispatchers.CategoryInspectActivity,
	})

	// -- logs

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "logs",
		Parent:  root,
		Summary: "View application logs",
		Description: `Shows the fp application log file.

By default, shows the last 50 lines. Use -n to change the number of lines.

Use --tail or -f to follow logs in real time (like tail -f).
Use --clear to empty the log file.

Log settings can be configured with:
  fp config set log_enabled false    Disable logging
  fp config set log_level debug      Set log level (debug, info, warn, error)`,
		Usage:    "fp logs [-n <lines>] [--tail] [--clear]",
		Flags:    LogsFlags,
		Action:   logsAction,
		Category: dispatchers.CategoryConfig,
	})

	// -- help

	dispatchers.NewNode(
		"help",
		root,
		"Show help for a command",
		"", // description not needed for help itself
		"fp help [command]",
		nil,
		nil,
		nil,
	)

	return root
}

// logsAction handles the logs command with its various flags
func logsAction(args []string, flags *dispatchers.ParsedFlags) error {
	if flags.Has("--clear") {
		return logsactions.Clear(args, flags)
	}
	if flags.Has("--tail") || flags.Has("-f") {
		return logsactions.Tail(args, flags)
	}
	return logsactions.View(args, flags)
}
