package actions

import (
	"runtime"

	"github.com/footprint-tools/cli/internal/dispatchers"
	"github.com/footprint-tools/cli/internal/output"
)

func ShowVersion(args []string, flags *dispatchers.ParsedFlags) error {
	return showVersion(args, flags, DefaultDeps())
}

func showVersion(_ []string, flags *dispatchers.ParsedFlags, deps Deps) error {
	if flags.Has("--json") {
		return showVersionJSON(deps)
	}
	_, _ = deps.Printf("fp version %v\n", deps.Version())
	return nil
}

func showVersionJSON(deps Deps) error {
	type versionInfo struct {
		Version   string `json:"version"`
		GoVersion string `json:"go_version"`
		OS        string `json:"os"`
		Arch      string `json:"arch"`
	}

	info := versionInfo{
		Version:   deps.Version(),
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	return output.JSON(deps.Println, info)
}
