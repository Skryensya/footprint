package store

import (
	"path/filepath"

	"github.com/Skryensya/footprint/internal/paths"
)

func DBPath() string {
	return filepath.Join(paths.AppDataDir(), "store.db")
}
