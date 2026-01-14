package hooks

import (
	"os"
	"path/filepath"
)

func Install(hooksPath string) error {
	fpPath, err := os.Executable()
	if err != nil {
		return err
	}

	for _, hook := range ManagedHooks {
		target := filepath.Join(hooksPath, hook)

		if Exists(target) {
			if err := BackupHook(hooksPath, hook); err != nil {
				return err
			}
		}

		script := Script(fpPath, hook)

		if err := os.WriteFile(target, []byte(script), 0755); err != nil {
			return err
		}
	}

	return nil
}
