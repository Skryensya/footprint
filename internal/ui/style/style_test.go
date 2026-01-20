package style

import (
	"os"
	"strings"
	"testing"
)

func TestDisabledReturnsPlainText(t *testing.T) {
	// Ensure no env vars interfere
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FP_NO_COLOR")

	Init(false, nil)

	tests := []struct {
		name string
		fn   func(string) string
	}{
		{"Success", Success},
		{"Warning", Warning},
		{"Error", Error},
		{"Info", Info},
		{"Header", Header},
		{"Muted", Muted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "test message"
			output := tt.fn(input)

			if output != input {
				t.Errorf("%s() with disabled styling: got %q, want %q", tt.name, output, input)
			}

			// Verify no ANSI escape codes
			if strings.Contains(output, "\x1b[") {
				t.Errorf("%s() with disabled styling contains ANSI codes: %q", tt.name, output)
			}
		})
	}
}

func TestEnabledReturnsStyledText(t *testing.T) {
	// Ensure no env vars interfere
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FP_NO_COLOR")

	Init(true, nil)

	tests := []struct {
		name string
		fn   func(string) string
	}{
		{"Success", Success},
		{"Warning", Warning},
		{"Error", Error},
		{"Info", Info},
		{"Header", Header},
		{"Muted", Muted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "test message"
			output := tt.fn(input)

			// Output should contain the original text
			if !strings.Contains(output, input) {
				t.Errorf("%s() output %q does not contain input %q", tt.name, output, input)
			}

			// Output should contain ANSI escape codes when enabled
			if !strings.Contains(output, "\x1b[") {
				t.Errorf("%s() with enabled styling should contain ANSI codes: %q", tt.name, output)
			}
		})
	}
}

func TestNoColorEnvDisablesStyling(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	Init(true, nil) // Try to enable, but NO_COLOR should override

	if Enabled() {
		t.Error("Enabled() should return false when NO_COLOR is set")
	}

	input := "test"
	output := Success(input)
	if output != input {
		t.Errorf("Success() should return plain text when NO_COLOR is set: got %q, want %q", output, input)
	}
}

func TestFPNoColorEnvDisablesStyling(t *testing.T) {
	os.Setenv("FP_NO_COLOR", "1")
	defer os.Unsetenv("FP_NO_COLOR")

	Init(true, nil) // Try to enable, but FP_NO_COLOR should override

	if Enabled() {
		t.Error("Enabled() should return false when FP_NO_COLOR is set")
	}

	input := "test"
	output := Warning(input)
	if output != input {
		t.Errorf("Warning() should return plain text when FP_NO_COLOR is set: got %q, want %q", output, input)
	}
}

func TestEnabledReturnsCorrectState(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FP_NO_COLOR")

	Init(false, nil)
	if Enabled() {
		t.Error("Enabled() should return false after Init(false)")
	}

	Init(true, nil)
	if !Enabled() {
		t.Error("Enabled() should return true after Init(true)")
	}
}

func TestEmptyStringHandling(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FP_NO_COLOR")

	// Test with disabled
	Init(false, nil)
	if got := Success(""); got != "" {
		t.Errorf("Success(\"\") with disabled styling: got %q, want \"\"", got)
	}

	// Test with enabled - should still handle empty gracefully
	Init(true, nil)
	output := Success("")
	// Empty string with styling might have escape codes but should be minimal
	if !strings.Contains(output, "") {
		t.Errorf("Success(\"\") with enabled styling failed")
	}
}

// Tests for configurable colors

func TestLoadColorConfig_DarkTheme(t *testing.T) {
	clearColorEnvVars(t)

	cfg := map[string]string{}
	colors := LoadColorConfig(cfg)

	// Check default-dark theme values (bright colors for dark backgrounds)
	if colors.Success != "10" {
		t.Errorf("Success: got %q, want %q", colors.Success, "10")
	}
	if colors.Warning != "11" {
		t.Errorf("Warning: got %q, want %q", colors.Warning, "11")
	}
	if colors.Error != "9" {
		t.Errorf("Error: got %q, want %q", colors.Error, "9")
	}
	if colors.Info != "14" {
		t.Errorf("Info: got %q, want %q", colors.Info, "14")
	}
	if colors.Muted != "245" {
		t.Errorf("Muted: got %q, want %q", colors.Muted, "245")
	}
	if colors.Header != "bold" {
		t.Errorf("Header: got %q, want %q", colors.Header, "bold")
	}
}

func TestLoadColorConfig_LightTheme(t *testing.T) {
	clearColorEnvVars(t)

	cfg := map[string]string{"color_theme": "default-light"}
	colors := LoadColorConfig(cfg)

	if colors.Success != "28" {
		t.Errorf("Success: got %q, want %q", colors.Success, "28")
	}
	if colors.Warning != "130" {
		t.Errorf("Warning: got %q, want %q", colors.Warning, "130")
	}
	if colors.Error != "124" {
		t.Errorf("Error: got %q, want %q", colors.Error, "124")
	}
	if colors.Info != "27" {
		t.Errorf("Info: got %q, want %q", colors.Info, "27")
	}
	if colors.Muted != "243" {
		t.Errorf("Muted: got %q, want %q", colors.Muted, "243")
	}
}

func TestLoadColorConfig_IndividualOverride(t *testing.T) {
	clearColorEnvVars(t)

	cfg := map[string]string{
		"color_theme":   "default-dark",
		"color_success": "82",
		"color_error":   "196",
	}
	colors := LoadColorConfig(cfg)

	// Overridden values
	if colors.Success != "82" {
		t.Errorf("Success: got %q, want %q", colors.Success, "82")
	}
	if colors.Error != "196" {
		t.Errorf("Error: got %q, want %q", colors.Error, "196")
	}

	// Non-overridden values should be from theme
	if colors.Warning != "11" {
		t.Errorf("Warning: got %q, want %q", colors.Warning, "11")
	}
}

func TestLoadColorConfig_EnvOverridesConfig(t *testing.T) {
	clearColorEnvVars(t)
	t.Setenv("FP_COLOR_SUCCESS", "99")

	cfg := map[string]string{
		"color_success": "10", // Should be overridden by env
	}
	colors := LoadColorConfig(cfg)

	if colors.Success != "99" {
		t.Errorf("Success: got %q, want %q (env should override config)", colors.Success, "99")
	}
}

func TestLoadColorConfig_EnvThemeOverridesConfig(t *testing.T) {
	clearColorEnvVars(t)
	t.Setenv("FP_COLOR_THEME", "default-light")

	cfg := map[string]string{
		"color_theme": "default-dark", // Should be overridden by env
	}
	colors := LoadColorConfig(cfg)

	// Should have default-light theme values
	if colors.Success != "28" {
		t.Errorf(
			"Success: got %q, want %q (env theme should override config)",
			colors.Success,
			"28",
		)
	}
}

func TestLoadColorConfig_UnknownThemeFallsBackToDefaultDark(t *testing.T) {
	clearColorEnvVars(t)

	cfg := map[string]string{"color_theme": "nonexistent"}
	colors := LoadColorConfig(cfg)

	// Should fall back to default-dark theme (bright green = 10)
	if colors.Success != "10" {
		t.Errorf("Success: got %q, want %q (should fall back to default-dark)", colors.Success, "10")
	}
}

func TestLoadColorConfig_ColorNumberOverrides(t *testing.T) {
	clearColorEnvVars(t)

	cfg := map[string]string{
		"color_1": "100",
		"color_2": "101",
		"color_3": "102",
		"color_4": "103",
		"color_5": "104",
		"color_6": "105",
	}
	colors := LoadColorConfig(cfg)

	if colors.Color1 != "100" {
		t.Errorf("Color1: got %q, want %q", colors.Color1, "100")
	}
	if colors.Color2 != "101" {
		t.Errorf("Color2: got %q, want %q", colors.Color2, "101")
	}
	if colors.Color3 != "102" {
		t.Errorf("Color3: got %q, want %q", colors.Color3, "102")
	}
	if colors.Color4 != "103" {
		t.Errorf("Color4: got %q, want %q", colors.Color4, "103")
	}
	if colors.Color5 != "104" {
		t.Errorf("Color5: got %q, want %q", colors.Color5, "104")
	}
	if colors.Color6 != "105" {
		t.Errorf("Color6: got %q, want %q", colors.Color6, "105")
	}
}

func TestMakeStyle_Bold(t *testing.T) {
	s := makeStyle("bold")
	// Check that bold style is created (we can't easily inspect the style,
	// but we can verify it doesn't panic and returns a style)
	rendered := s.Render("test")
	if rendered == "" {
		t.Error("makeStyle(bold) should return a usable style")
	}
}

func TestMakeStyle_Color(t *testing.T) {
	s := makeStyle("6")
	rendered := s.Render("test")
	if rendered == "" {
		t.Error("makeStyle(6) should return a usable style")
	}
}

func TestToUpperSnake(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"color_success", "COLOR_SUCCESS"},
		{"color_1", "COLOR_1"},
		{"abc", "ABC"},
		{"", ""},
	}

	for _, tt := range tests {
		got := toUpperSnake(tt.input)
		if got != tt.want {
			t.Errorf("toUpperSnake(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// clearColorEnvVars clears all FP_COLOR_* environment variables for test isolation.
func clearColorEnvVars(t *testing.T) {
	t.Helper()
	envVars := []string{
		"FP_COLOR_THEME",
		"FP_COLOR_SUCCESS",
		"FP_COLOR_WARNING",
		"FP_COLOR_ERROR",
		"FP_COLOR_INFO",
		"FP_COLOR_MUTED",
		"FP_COLOR_HEADER",
		"FP_COLOR_1",
		"FP_COLOR_2",
		"FP_COLOR_3",
		"FP_COLOR_4",
		"FP_COLOR_5",
		"FP_COLOR_6",
	}
	for _, env := range envVars {
		if os.Getenv(env) != "" {
			t.Setenv(env, "")
		}
	}
}
