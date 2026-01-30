package update

import (
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/footprint-tools/cli/internal/app"
)

type Deps struct {
	Stdout         io.Writer
	Stderr         io.Writer
	HTTPClient     HTTPClient
	CurrentVersion string
	ExecutablePath func() (string, error)
	RunCommand     func(name string, args ...string) error
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func DefaultDeps() Deps {
	return Deps{
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		HTTPClient:     http.DefaultClient,
		CurrentVersion: app.Version,
		ExecutablePath: os.Executable,
		RunCommand: func(name string, args ...string) error {
			cmd := exec.Command(name, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
	}
}
