package actions

import (
	"fmt"

	"github.com/Skryensya/footprint/internal/config"
)

func ConfigList(_ []string, _ []string) error {
	lines, err := config.ReadLines()
	if err != nil {
		return err
	}

	configMap, err := config.Parse(lines)
	if err != nil {
		return err
	}

	for key, value := range configMap {
		fmt.Printf("%s=%s\n", key, value)
	}

	return nil
}
