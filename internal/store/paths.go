package store

import (
	"path/filepath"

	"github.com/footprint-tools/footprint-cli/internal/paths"
)

func DBPath() string {
	return filepath.Join(paths.AppDataDir(), "store.db")
}
