package config

import (
	"github.com/Skryensya/footprint/internal/paths"
)

// Default configuration values (in code, not persisted)
var Defaults = map[string]func() string{
	"export_interval": func() string { return "3600" },
	"export_repo":     func() string { return paths.ExportRepoDir() },
	"export_last":     func() string { return "0" },
	"color_theme":     func() string { return "default-dark" },
	"log_enabled":     func() string { return "true" },
	"log_level":       func() string { return "debug" }, // debug, info, warn, error
}

// Get returns the value for a config key.
// It checks the config file first, then falls back to the default.
// Returns the value and whether it was found (in file or defaults).
func Get(key string) (string, bool) {
	lines, err := ReadLines()
	if err != nil {
		// On error, try defaults
		if defaultFn, ok := Defaults[key]; ok {
			return defaultFn(), true
		}
		return "", false
	}

	cfg, err := Parse(lines)
	if err != nil {
		if defaultFn, ok := Defaults[key]; ok {
			return defaultFn(), true
		}
		return "", false
	}

	// Check config file first
	if value, exists := cfg[key]; exists {
		return value, true
	}

	// Fall back to default
	if defaultFn, ok := Defaults[key]; ok {
		return defaultFn(), true
	}

	return "", false
}

// GetAll returns all config values (user overrides merged with defaults).
func GetAll() (map[string]string, error) {
	result := make(map[string]string)

	// Start with defaults
	for key, valueFn := range Defaults {
		result[key] = valueFn()
	}

	// Override with user config
	lines, err := ReadLines()
	if err != nil {
		return result, nil // Return defaults on error
	}

	cfg, err := Parse(lines)
	if err != nil {
		return result, nil
	}

	for key, value := range cfg {
		result[key] = value
	}

	return result, nil
}
