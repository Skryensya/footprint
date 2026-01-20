package theme

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/ui/style"
	"github.com/Skryensya/footprint/internal/usage"
)

func Set(args []string, flags *dispatchers.ParsedFlags) error {
	return setTheme(args, flags, DefaultDeps())
}

func setTheme(args []string, _ *dispatchers.ParsedFlags, deps Deps) error {
	if len(args) < 1 {
		return usage.MissingArgument("theme")
	}

	themeName := args[0]

	// Validate theme exists
	if _, ok := deps.Themes[themeName]; !ok {
		deps.Printf("%s Unknown theme: %s\n", style.Error("error:"), themeName)
		deps.Println("")
		deps.Println("Available themes:")
		for _, name := range deps.ThemeNames {
			deps.Printf("  %s\n", name)
		}
		return fmt.Errorf("unknown theme: %s", themeName)
	}

	lines, err := deps.ReadLines()
	if err != nil {
		return err
	}

	lines, _ = deps.Set(lines, "color_theme", themeName)

	if err := deps.WriteLines(lines); err != nil {
		return err
	}

	deps.Printf("Theme set to %s\n", style.Success(themeName))

	return nil
}
