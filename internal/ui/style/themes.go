package style

import "os"

// ColorConfig holds all configurable colors for the UI.
// Values can be ANSI color numbers (0-255) or "bold" for bold styling.
type ColorConfig struct {
	Success string
	Warning string
	Error   string
	Info    string
	Muted   string
	Header  string
	Color1  string
	Color2  string
	Color3  string
	Color4  string
	Color5  string
	Color6  string
}

// ThemeNames lists available themes in display order.
var ThemeNames = []string{
	"default-dark", "default-light",
	"ocean-dark", "ocean-light",
	"forest-dark", "forest-light",
}

// Themes contains the built-in color themes.
// Dark themes use BRIGHT colors (high contrast on dark backgrounds).
// Light themes use DARK colors (high contrast on light/white backgrounds).
var Themes = map[string]ColorConfig{
	// Default dark - bright classic colors for dark terminal backgrounds.
	"default-dark": {
		Success: "10",  // bright green
		Warning: "11",  // bright yellow
		Error:   "9",   // bright red
		Info:    "14",  // bright cyan
		Muted:   "246", // light gray
		Header:  "bold",
		Color1:  "10", // POST-COMMIT  (green)
		Color2:  "13", // POST-REWRITE (magenta)
		Color3:  "12", // POST-CHECKOUT (blue)
		Color4:  "14", // POST-MERGE (cyan)
		Color5:  "11", // PRE-PUSH (yellow)
		Color6:  "15", // MANUAL (white)
	},
	// Default light - dark colors for light terminal backgrounds.
	"default-light": {
		Success: "28",  // strong green
		Warning: "136", // amber
		Error:   "124", // strong red
		Info:    "25",  // deep blue
		Muted:   "245", // medium gray (not too faint)
		Header:  "bold",
		Color1:  "28",  // POST-COMMIT
		Color2:  "90",  // POST-REWRITE
		Color3:  "25",  // POST-CHECKOUT
		Color4:  "30",  // POST-MERGE
		Color5:  "136", // PRE-PUSH
		Color6:  "250", // MANUAL
	},
	// Ocean dark - bright blue/cyan palette for dark backgrounds.
	"ocean-dark": {
		Success: "51",  // neon cyan
		Warning: "220", // gold (contrast anchor)
		Error:   "203", // coral red
		Info:    "117", // sky blue
		Muted:   "245",
		Header:  "bold",

		// Event sources (cold spectrum)
		Color1: "48",  // commit (teal-green)
		Color2: "135", // rewrite (violet)
		Color3: "75",  // checkout (blue)
		Color4: "51",  // merge (cyan)
		Color5: "220", // push (gold)
		Color6: "159", // manual (ice white)
	},
	// Ocean light - deep blue/teal for light backgrounds.
	"ocean-light": {
		Success: "29",  // deep teal
		Warning: "136", // amber
		Error:   "124", // strong red
		Info:    "17",  // navy (very dark blue)
		Muted:   "244",
		Header:  "bold",

		// Event sources (ink spectrum)
		Color1: "30",  // commit (deep teal)
		Color2: "97",  // rewrite (purple)
		Color3: "27",  // checkout (navy)
		Color4: "31",  // merge (cyan-teal)
		Color5: "136", // push (amber)
		Color6: "250", // manual
	},
	// Forest dark - bright green/earth palette for dark backgrounds.
	"forest-dark": {
		Success: "120",
		Warning: "222",
		Error:   "203",
		Info:    "151",
		Muted:   "245",
		Header:  "bold",
		Color1:  "120", // POST-COMMIT   → leaf green
		Color2:  "135", // POST-REWRITE  → purple (danger / anomaly)
		Color3:  "110", // POST-CHECKOUT → shadow blue
		Color4:  "151", // POST-MERGE    → pale mint (integration)
		Color5:  "220", // PRE-PUSH      → sunlight yellow
		Color6:  "180", // MANUAL        → bone / neutral

	},
	// Forest light - deep green/earth for light backgrounds.
	"forest-light": {
		Success: "28", // forest green
		Warning: "94",
		Error:   "124",
		Info:    "29", // sea green
		Muted:   "244",
		Header:  "bold",
		Color1:  "28",  // POST-COMMIT   → dark green
		Color2:  "95",  // POST-REWRITE  → muted purple
		Color3:  "24",  // POST-CHECKOUT → ink blue
		Color4:  "30",  // POST-MERGE    → teal
		Color5:  "136", // PRE-PUSH      → amber
		Color6:  "250", // MANUAL        → near-white

	},
}

// colorConfigKeys maps config/env key names to ColorConfig field names.
var colorConfigKeys = map[string]string{
	"color_success": "Success",
	"color_warning": "Warning",
	"color_error":   "Error",
	"color_info":    "Info",
	"color_muted":   "Muted",
	"color_header":  "Header",
	"color_1":       "Color1",
	"color_2":       "Color2",
	"color_3":       "Color3",
	"color_4":       "Color4",
	"color_5":       "Color5",
	"color_6":       "Color6",
}

// LoadColorConfig builds a ColorConfig from the given configuration map.
// Resolution priority:
// 1. Environment variable (FP_COLOR_*)
// 2. Config file value
// 3. Theme value (from color_theme config)
// 4. Default theme
func LoadColorConfig(cfg map[string]string) ColorConfig {
	// Start with default-dark theme
	themeName := "default-dark"

	// Check env for theme override
	if envTheme := os.Getenv("FP_COLOR_THEME"); envTheme != "" {
		themeName = envTheme
	} else if cfgTheme, ok := cfg["color_theme"]; ok && cfgTheme != "" {
		themeName = cfgTheme
	}

	// Get base theme (fall back to default-dark if unknown)
	theme, ok := Themes[themeName]
	if !ok {
		theme = Themes["default-dark"]
	}

	// Apply overrides from config and env
	result := theme

	for configKey, fieldName := range colorConfigKeys {
		// Check env first (highest priority)
		envKey := "FP_" + toUpperSnake(configKey)
		if envVal := os.Getenv(envKey); envVal != "" {
			setColorField(&result, fieldName, envVal)
			continue
		}

		// Check config file
		if cfgVal, ok := cfg[configKey]; ok && cfgVal != "" {
			setColorField(&result, fieldName, cfgVal)
		}
	}

	return result
}

// setColorField sets a field on ColorConfig by name.
func setColorField(c *ColorConfig, field, value string) {
	switch field {
	case "Success":
		c.Success = value
	case "Warning":
		c.Warning = value
	case "Error":
		c.Error = value
	case "Info":
		c.Info = value
	case "Muted":
		c.Muted = value
	case "Header":
		c.Header = value
	case "Color1":
		c.Color1 = value
	case "Color2":
		c.Color2 = value
	case "Color3":
		c.Color3 = value
	case "Color4":
		c.Color4 = value
	case "Color5":
		c.Color5 = value
	case "Color6":
		c.Color6 = value
	}
}

// toUpperSnake converts "color_success" to "COLOR_SUCCESS".
func toUpperSnake(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 'a' + 'A'
		} else {
			result[i] = c
		}
	}
	return string(result)
}
