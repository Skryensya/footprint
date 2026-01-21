package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/app"
)

type actionDependencies struct {
	Printf  func(format string, a ...any) (n int, err error)
	Version func() string
}

func defaultDeps() actionDependencies {
	return actionDependencies{
		Printf:  fmt.Printf,
		Version: func() string { return app.Version },
	}
}
