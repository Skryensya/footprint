package components

import "github.com/charmbracelet/bubbles/key"

// CommonKeyMap defines common keybindings used across TUI views.
type CommonKeyMap struct {
	Quit        key.Binding
	Tab         key.Binding
	Enter       key.Binding
	Escape      key.Binding
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Home        key.Binding
	End         key.Binding
	Help        key.Binding
	Filter      key.Binding
	ClearFilter key.Binding
}

// NewCommonKeyMap returns the default common keybindings.
func NewCommonKeyMap() CommonKeyMap {
	return CommonKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "switch"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "back"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("PgUp", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("PgDn", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g", "top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G", "bottom"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "clear"),
		),
	}
}

// NavigationBindings returns bindings for list navigation.
func NavigationBindings() []key.Binding {
	km := NewCommonKeyMap()
	return []key.Binding{
		km.Tab,
		km.Up,
		km.Down,
		km.Enter,
		km.Quit,
	}
}

// FilterBindings returns bindings when a filter is active.
func FilterBindings() []key.Binding {
	km := NewCommonKeyMap()
	return []key.Binding{
		km.Escape,
		km.ClearFilter,
	}
}

// EditBindings returns bindings for edit mode.
func EditBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "save")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc", "cancel")),
	}
}

// ConfigKeyMap defines keybindings for the config editor.
type ConfigKeyMap struct {
	Common  CommonKeyMap
	Edit    key.Binding
	Default key.Binding
	Unset   key.Binding
}

// NewConfigKeyMap returns keybindings for the config editor.
func NewConfigKeyMap() ConfigKeyMap {
	return ConfigKeyMap{
		Common: NewCommonKeyMap(),
		Edit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "edit"),
		),
		Default: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "default"),
		),
		Unset: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "unset"),
		),
	}
}

// LogsKeyMap defines keybindings for the logs viewer.
type LogsKeyMap struct {
	Common     CommonKeyMap
	Pause      key.Binding
	AutoScroll key.Binding
	LevelError key.Binding
	LevelWarn  key.Binding
	LevelInfo  key.Binding
	LevelDebug key.Binding
}

// NewLogsKeyMap returns keybindings for the logs viewer.
func NewLogsKeyMap() LogsKeyMap {
	return LogsKeyMap{
		Common: NewCommonKeyMap(),
		Pause: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause"),
		),
		AutoScroll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "auto"),
		),
		LevelError: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "ERROR"),
		),
		LevelWarn: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "WARN"),
		),
		LevelInfo: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "INFO"),
		),
		LevelDebug: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "DEBUG"),
		),
	}
}

// WatchKeyMap defines keybindings for the watch view.
type WatchKeyMap struct {
	Common  CommonKeyMap
	Pause   key.Binding
	Source1 key.Binding
	Source2 key.Binding
	Source3 key.Binding
	Source4 key.Binding
	Source5 key.Binding
	Source6 key.Binding
	Source7 key.Binding
}

// NewWatchKeyMap returns keybindings for the watch view.
func NewWatchKeyMap() WatchKeyMap {
	return WatchKeyMap{
		Common: NewCommonKeyMap(),
		Pause: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause"),
		),
		Source1: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "commit"),
		),
		Source2: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "rewrite"),
		),
		Source3: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "checkout"),
		),
		Source4: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "merge"),
		),
		Source5: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "push"),
		),
		Source6: key.NewBinding(
			key.WithKeys("6"),
			key.WithHelp("6", "manual"),
		),
		Source7: key.NewBinding(
			key.WithKeys("7"),
			key.WithHelp("7", "backfill"),
		),
	}
}
