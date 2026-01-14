# Footprint

footprint is a local tool that records activity from Git repositories and stores it on your machine.

It helps you keep a structured history of work without using external services or sending data anywhere.

## What it does

footprint records events that happen in a Git repository and saves them locally.

The data is stored in a way that can be inspected, processed, or exported later by the user.

Everything stays under your control.

## What it does not do

* footprint does not upload data
* footprint does not track time
* footprint does not monitor behavior
* footprint does not depend on online services

## Current state

footprint records Git activity and stores it in a local database.

That is its full responsibility today.

## Usage

footprint is a command line tool called `fp`.

It is meant to be used inside Git repositories you want to track.

The basic workflow is simple.

You install the Git hooks.
You choose which repositories to track.
footprint records activity automatically.
You can later inspect what was recorded.

### Install Git hooks

The first step is to install the Git hooks.

Run this inside a Git repository:

```
fp hooks install
```

You can check which hooks are installed:

```
fp hooks status
```

You can remove them later with:

```
fp hooks uninstall
```

### Track a repository

After installing the hooks, choose which repositories to track.

Run this inside a Git repository:

```
fp repo track
```

You can check the status at any time:

```
fp repo status
```

Stop tracking a repository:

```
fp repo untrack
```

### Inspect recorded activity

To see what has been recorded:

```
fp activity list
```

This shows recorded activity, newest first.

### Manage repositories

List all tracked repositories:

```
fp repo list
```

If a repository remote changes and you want to keep the same history:

```
fp repo adopt-remote
```

### Configuration

You can manage configuration values using:

```
fp config list
fp config get <key>
fp config set <key> <value>
fp config unset <key>
```

Configuration is stored locally and applies only to your machine.

### Help and version

To see all commands:

```
fp --help
```

To see help for a specific command:

```
fp help <command>
```

To check the installed version:

```
fp version
```

## Roadmap

Planned features include:

* Exporting data to CSV
* Optional sync to a user owned Git repository
* Adding a full testing suite

## Privacy

All data stays on your machine.
Nothing is shared unless you choose to export it.
There is no telemetry.

## License

MIT License. See the LICENSE file for details.
