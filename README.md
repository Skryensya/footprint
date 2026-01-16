# Footprint

Footprint is a local tool that records activity from Git repositories and stores it on your machine.

It helps you keep a structured history of work without using external services or sending data anywhere.

## What it does

Footprint records events that happen in Git repositories and saves them locally in a SQLite database.

Recorded events include commits, merges, checkouts, rebases, and pushes.

The data can be inspected, filtered, or exported later. Everything stays under your control.

## What it does not do

* Footprint does not upload data
* Footprint does not track time
* Footprint does not monitor behavior
* Footprint does not depend on online services

## Installation

Build from source:

```
make build
```

This creates the `fp` binary in the current directory.

## Quick start

```
fp setup              # Install git hooks (in current repo)
fp track              # Start tracking the current repository
```

From now on, footprint records activity automatically. You can inspect it with:

```
fp activity           # Show recorded events
fp log                # Watch for new events in real time
```

## Commands

### Getting started

Install git hooks in the current repository:

```
fp setup
```

Install git hooks globally (applies to all repositories):

```
fp setup --global
```

Start tracking a repository:

```
fp track [path]
```

### Inspecting activity

Show recorded activity (newest first):

```
fp activity
```

Filter activity with flags:

```
fp activity --oneline              # Compact one-line format
fp activity --since=2024-01-01     # Events after date
fp activity --until=2024-12-31     # Events before date
fp activity --status=pending       # Filter by status
fp activity --source=post-commit   # Filter by source
fp activity --repo=<id>            # Filter by repository
fp activity --limit=50             # Limit results
```

Watch for new events in real time:

```
fp log
fp log --oneline
```

This runs continuously like `tail -f`. Press Ctrl+C to stop.

### Repository status

Check tracking status of current repository:

```
fp status [path]
```

List all tracked repositories:

```
fp repos
fp list              # Alias for repos
```

### Managing repositories

Stop tracking a repository:

```
fp untrack [path]
```

Update repository ID after remote URL changes:

```
fp sync-remote [path]
```

### Managing hooks

Check installed hooks:

```
fp check
fp check --global
```

Remove hooks:

```
fp teardown
fp teardown --global
```

### Configuration

```
fp config list                  # Show all config values
fp config get <key>             # Get a value
fp config set <key> <value>     # Set a value
fp config unset <key>           # Remove a value
fp config unset --all           # Remove all values
```

Configuration is stored in `~/.fprc`.

### Help

```
fp --help                # Show all commands
fp help <command>        # Help for a specific command
fp help <topic>          # Conceptual documentation
fp version               # Show version
```

Available help topics: `overview`, `workflow`, `hooks`, `data`.

## Data storage

Events are stored in a SQLite database at:

- macOS: `~/Library/Application Support/Footprint/store.db`
- Linux: `~/.config/Footprint/store.db`

Tracked repositories are stored in `~/.fprc`.

## Privacy

All data stays on your machine.
Nothing is shared unless you choose to export it.
There is no telemetry.

## License

MIT License. See the LICENSE file for details.
