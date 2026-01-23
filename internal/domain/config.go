package domain

// ConfigKey defines a configuration key with its metadata.
type ConfigKey struct {
	Name        string
	Default     string
	Description string
	Hidden      bool // Hidden keys are not shown in help or config list
	HideIfEmpty bool // Only show in config list if explicitly set
}

// ConfigKeys defines all available configuration keys.
// This is the single source of truth for configuration.
// Order determines display order in `fp config list`.
var ConfigKeys = []ConfigKey{
	// Display
	{
		Name:        "pager",
		Default:     "less -FRSX",
		Description: "Pager command for long output",
	},
	{
		Name:        "theme",
		Default:     "default",
		Description: "Color theme: default, neon, aurora, mono, ocean, sunset, candy, contrast",
	},
	{
		Name:        "display_date",
		Default:     "Jan 02",
		Description: "Display date format: dd/mm/yyyy, mm/dd/yyyy, yyyy-mm-dd, or Go format (e.g., Jan 02, Jan 02 2006)",
	},
	{
		Name:        "display_time",
		Default:     "24h",
		Description: "Display time format: 12h, 24h",
	},
	{
		Name:        "color_success",
		Description: "ANSI color code override for success messages",
		HideIfEmpty: true,
	},
	{
		Name:        "color_warning",
		Description: "ANSI color code override for warning messages",
		HideIfEmpty: true,
	},
	{
		Name:        "color_error",
		Description: "ANSI color code override for error messages",
		HideIfEmpty: true,
	},
	{
		Name:        "color_info",
		Description: "ANSI color code override for info messages",
		HideIfEmpty: true,
	},
	{
		Name:        "color_muted",
		Description: "ANSI color code override for muted text",
		HideIfEmpty: true,
	},
	{
		Name:        "color_header",
		Description: "Style override for headers (e.g., 'bold')",
		HideIfEmpty: true,
	},
	// Logging
	{
		Name:        "enable_log",
		Default:     "true",
		Description: "Enable logging to file",
	},
	// Export
	{
		Name:        "export_interval_sec",
		Default:     "3600",
		Description: "Seconds between automatic exports",
	},
	{
		Name:        "export_path",
		Default:     "", // Set dynamically to paths.ExportRepoDir()
		Description: "Path to the export repository",
	},
	{
		Name:        "export_remote",
		Default:     "",
		Description: "Remote URL for syncing exports",
	},
	// Hidden (internal)
	{
		Name:        "export_last",
		Default:     "0",
		Description: "Unix timestamp of last export",
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
