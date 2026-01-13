package actions

import (
	"os"
	"path/filepath"
)

func hasFlag(flags []string, name string) bool {
	for _, f := range flags {
		if f == name {
			return true
		}
	}
	return false
}

func resolvePath(args []string) (string, error) {
	p := "."
	if len(args) > 0 {
		p = args[0]
	}

	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", os.ErrInvalid
	}

	return abs, nil
}
