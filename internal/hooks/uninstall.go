package hooks

import (
	"os"
	"path/filepath"

	"github.com/Skryensya/footprint/internal/log"
)

func Uninstall(hooksPath string) error {
	log.Debug("hooks: uninstalling from %s", hooksPath)
	backupDir := BackupDir(hooksPath)

	for _, hook := range ManagedHooks {
		target := filepath.Join(hooksPath, hook)
		backup := filepath.Join(backupDir, hook)

		if Exists(backup) {
			_ = os.Remove(target)
			if err := os.Rename(backup, target); err != nil {
				log.Error("hooks: failed to restore backup for %s: %v", hook, err)
				return err
			}
			log.Debug("hooks: restored %s from backup", hook)
			continue
		}

		if Exists(target) {
			if err := os.Remove(target); err != nil {
				log.Error("hooks: failed to remove %s: %v", hook, err)
				return err
			}
			log.Debug("hooks: removed %s", hook)
		}
	}

	log.Info("hooks: uninstalled from %s", hooksPath)
	return nil
}
