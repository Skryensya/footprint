package domain

// ConfigKey defines a configuration key with its metadata.
type ConfigKey struct {
	Name        string
	Default     string
	Description string
	Hidden      bool // Hidden keys are not shown in help
}

// ConfigKeys defines all available configuration keys.
// This is the single source of truth for configuration.
var ConfigKeys = []ConfigKey{
	{
		Name:        "export_interval",
		Default:     "3600",
		Description: "Seconds between automatic exports (default: 3600 = 1 hour)",
	},
	{
		Name:        "export_repo",
		Default:     "", // Will be set dynamically to paths.ExportRepoDir()
		Description: "Path to the export repository",
	},
	{
		Name:        "export_last",
		Default:     "0",
		Description: "Unix timestamp of last export (managed internally)",
		Hidden:      true,
	},
	{
		Name:        "export_remote",
		Default:     "",
		Description: "Remote URL for syncing exports",
	},
	{
		Name:        "log_enabled",
		Default:     "true",
		Description: "Enable debug logging to file",
	},
	{
		Name:        "log_level",
		Default:     "debug",
		Description: "Minimum log level: debug, info, warn, error",
	},
	{
		Name:        "color_theme",
		Default:     "default",
		Description: "Color theme: default, default-dark, default-light",
	},
	{
		Name:        "color_success",
		Default:     "",
		Description: "ANSI color code for success messages",
	},
	{
		Name:        "color_warning",
		Default:     "",
		Description: "ANSI color code for warning messages",
	},
	{
		Name:        "color_error",
		Default:     "",
		Description: "ANSI color code for error messages",
	},
	{
		Name:        "color_info",
		Default:     "",
		Description: "ANSI color code for info messages",
	},
	{
		Name:        "color_muted",
		Default:     "",
		Description: "ANSI color code for muted text",
	},
	{
		Name:        "color_header",
		Default:     "",
		Description: "Style for headers (e.g., 'bold')",
	},
	{
		Name:        "pager",
		Default:     "",
		Description: "Pager command (default: less -FRSX)",
	},
	{
		Name:        "update_last_check",
		Default:     "",
		Description: "ISO8601 timestamp of last update check (managed internally)",
		Hidden:      true,
	},
	{
		Name:        "update_latest_version",
		Default:     "",
		Description: "Latest known version from GitHub (managed internally)",
		Hidden:      true,
	},
}

// configKeyMap is a lookup map for configuration keys.
var configKeyMap map[string]ConfigKey

func init() {
	configKeyMap = make(map[string]ConfigKey, len(ConfigKeys))
	for _, key := range ConfigKeys {
		configKeyMap[key.Name] = key
	}
}

// GetConfigKey returns the ConfigKey for a given name.
func GetConfigKey(name string) (ConfigKey, bool) {
	key, ok := configKeyMap[name]
	return key, ok
}

// IsValidConfigKey checks if a key name is valid.
func IsValidConfigKey(name string) bool {
	_, ok := configKeyMap[name]
	return ok
}

// GetDefaultValue returns the default value for a config key.
func GetDefaultValue(name string) (string, bool) {
	if key, ok := configKeyMap[name]; ok {
		return key.Default, true
	}
	return "", false
}

// VisibleConfigKeys returns all non-hidden configuration keys.
func VisibleConfigKeys() []ConfigKey {
	var visible []ConfigKey
	for _, key := range ConfigKeys {
		if !key.Hidden {
			visible = append(visible, key)
		}
	}
	return visible
}
