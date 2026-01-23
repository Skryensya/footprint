package cli

import (
	"github.com/footprint-tools/footprint-cli/internal/actions"
	configactions "github.com/footprint-tools/footprint-cli/internal/actions/config"
	logsactions "github.com/footprint-tools/footprint-cli/internal/actions/logs"
	setupactions "github.com/footprint-tools/footprint-cli/internal/actions/setup"
	themeactions "github.com/footprint-tools/footprint-cli/internal/actions/theme"
	trackingactions "github.com/footprint-tools/footprint-cli/internal/actions/tracking"
	updateactions "github.com/footprint-tools/footprint-cli/internal/actions/update"
	"github.com/footprint-tools/footprint-cli/internal/dispatchers"
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

	addConfigCommands(root)
	addThemeCommands(root)
	addTrackingCommands(root)
	addActivityCommands(root)
	addSetupCommands(root)
	addLogsCommand(root)
	addUpdateCommand(root)
	addHelpCommand(root)

	return root
}

func addConfigCommands(root *dispatchers.DispatchNode) {
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
The configuration file is created if it does not exist.

Common configuration keys:
  export_remote       Remote URL for syncing exports (configures git remote)
  export_interval_sec Seconds between auto-exports (default: 3600)
  theme               Color theme name
  enable_log          Enable/disable logging (true/false)
  display_date        Date format: Jan 02, dd/mm/yyyy, mm/dd/yyyy, or Go format
  display_time        Time format: 12h, 24h

Example:
  fp config set export_remote git@github.com:user/my-exports.git`,
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
}

func addThemeCommands(root *dispatchers.DispatchNode) {
	theme := dispatchers.Group(dispatchers.GroupSpec{
		Name:    "theme",
		Parent:  root,
		Summary: "Manage color themes",
		Description: `View and change color themes.

fp includes several built-in themes optimized for different terminal
backgrounds. Each theme has a dark and light variant.

Available themes:
  default    Traditional terminal colors
  neon       Vivid cyberpunk colors
  aurora     Dreamy purples and teals
  mono       Minimalist grayscale with accent
  ocean      Cool blues and teals
  sunset     Warm oranges to purple gradient
  candy      Sweet pastel colors
  contrast   Maximum readability primaries

Each theme has -dark and -light variants.

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
		Category: dispatchers.CategoryTheme,
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
		Category: dispatchers.CategoryTheme,
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
		Category: dispatchers.CategoryTheme,
	})
}

func addTrackingCommands(root *dispatchers.DispatchNode) {
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
		Action:   trackingactions.SyncRemote,
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
}

func addActivityCommands(root *dispatchers.DispatchNode) {
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
		Summary: "Export pending events to CSV (invoked automatically)",
		Description: `Exports all pending events to CSV files in the export repository.

This command is normally invoked automatically after recording events.
You typically don't need to run it manually.

Events are exported to a flat CSV structure with year-based rotation:
  commits.csv       Current year's events
  commits-2024.csv  Events from 2024
  commits-2023.csv  Events from 2023

The export repository is a git repository where CSV files are committed.
Default location: ~/.config/Footprint/exports

Configuration (via 'fp config set'):
  export_remote    Remote URL for syncing exports
  export_interval_sec  Seconds between auto-exports (default: 3600)
  export_path      Path to export repository

Use --force to export immediately regardless of interval.
Use --open to open the export directory in file manager.`,
		Usage:    "fp export [--force] [--dry-run] [--open]",
		Action:   trackingactions.Export,
		Flags:    ExportFlags,
		Category: dispatchers.CategoryPlumbing,
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
}

func addSetupCommands(root *dispatchers.DispatchNode) {
	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "setup",
		Parent:  root,
		Summary: "Install fp git hooks (global)",
		Description: `Installs git hooks globally that automatically record activity.

Hooks are installed in git's global hooks directory, which applies to
all repositories. Use --repo to install only in the current repository.

Installed hooks:
  post-commit     Records commits
  post-merge      Records merges (including pulls)
  post-checkout   Records branch switches
  post-rewrite    Records rebases and amends
  pre-push        Records push attempts

If existing hooks are found, they are backed up before installation.

After setup, use 'fp track' in each repo you want to track.`,
		Usage:    "fp setup [--repo]",
		Flags:    SetupFlags,
		Action:   setupactions.Setup,
		Category: dispatchers.CategoryGetStarted,
	})

	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "teardown",
		Parent:  root,
		Summary: "Remove fp git hooks",
		Description: `Removes git hooks installed by fp.

By default, removes global hooks. Use --repo to remove hooks from
the current repository only.

If hooks were backed up during installation, the original hooks are
restored.`,
		Usage:    "fp teardown [--repo]",
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

By default checks global hooks. Use --repo to check the current
repository instead.`,
		Usage:    "fp check [--repo]",
		Flags:    CheckFlags,
		Action:   setupactions.Check,
		Category: dispatchers.CategoryInspectActivity,
	})
}

func addLogsCommand(root *dispatchers.DispatchNode) {
	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "logs",
		Parent:  root,
		Summary: "View application logs",
		Description: `Shows the fp application log file.

By default, shows the last 50 lines. Use -n to change the number of lines.

Use -i for an interactive viewer with filtering, navigation, and auto-tail.
Use --tail or -f to follow logs in real time (like tail -f).
Use --clear to empty the log file.

Log settings can be configured with:
  fp config set enable_log false    Disable logging
  fp config set log_level debug      Set log level (debug, info, warn, error)`,
		Usage:    "fp logs [-i] [-n <lines>] [--tail] [--clear]",
		Flags:    LogsFlags,
		Action:   logsAction,
		Category: dispatchers.CategoryInspectActivity,
	})
}

func addUpdateCommand(root *dispatchers.DispatchNode) {
	dispatchers.Command(dispatchers.CommandSpec{
		Name:    "update",
		Parent:  root,
		Summary: "Update fp to the latest version",
		Description: `Downloads and installs a newer version of fp.

By default, installs the latest GitHub release. Specify a version to install
a specific release (e.g., 'fp update v0.1.0').

The command downloads a pre-built binary for your OS and architecture.
If no binary is available, it falls back to building from source using
'go install' (requires Go to be installed).

Use --tag to install from a git tag that hasn't been released yet.
This always uses 'go install' and requires Go.

Examples:
  fp update              Install latest release
  fp update v0.1.0       Install specific release
  fp update --tag v0.2.0-beta  Build from tag using go install`,
		Usage:    "fp update [version] [--tag]",
		Args:     OptionalVersionArg,
		Flags:    UpdateFlags,
		Action:   updateactions.Update,
		Category: dispatchers.CategoryManageRepos,
	})
}

func addHelpCommand(root *dispatchers.DispatchNode) {
	dispatchers.NewNode(
		"help",
		root,
		"Show help for a command",
		`Shows help for fp commands and topics.

Use -i or --interactive to open an interactive TUI browser where you can
navigate through all commands and topics with keyboard shortcuts.

In interactive mode:
  - Use arrow keys or j/k to navigate the sidebar
  - Press / to search and filter commands
  - Press Tab to switch focus between sidebar and content
  - Use PgUp/PgDn or u/d to scroll content
  - Press q or Esc to exit`,
		"fp help [-i] [command]",
		[]dispatchers.FlagDescriptor{
			{Names: []string{"-i", "--interactive"}, Description: "Open interactive help browser"},
		},
		nil,
		nil,
	)
}

// logsAction handles the logs command with its various flags
func logsAction(args []string, flags *dispatchers.ParsedFlags) error {
	if flags.Has("--clear") {
		return logsactions.Clear(args, flags)
	}
	if flags.Has("-i") || flags.Has("--interactive") {
		return logsactions.Interactive(args, flags)
	}
	if flags.Has("--tail") || flags.Has("-f") {
		return logsactions.Tail(args, flags)
	}
	return logsactions.View(args, flags)
}
