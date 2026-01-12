package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/app"
)

func ShowVersion(args []string, flags []string) error {
	fmt.Printf("fp version %v\n", app.Version)
	return nil

}
