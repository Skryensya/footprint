//go:build !windows

package tracking

import (
	"fmt"
	"syscall"

	"github.com/footprint-tools/cli/internal/log"
)

// checkDiskSpace verifies there's enough space to write the estimated bytes.
func checkDiskSpace(dir string, requiredBytes int64) error {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(dir, &stat); err != nil {
		// If we can't check, proceed anyway (might be a virtual filesystem)
		log.Debug("export: could not check disk space for %s: %v", dir, err)
		return nil
	}

	// Available space = available blocks * block size
	available := int64(stat.Bavail) * int64(stat.Bsize)

	// Require at least 2x the estimated size for safety margin
	if available < requiredBytes*2 {
		return fmt.Errorf("need %d bytes, only %d available", requiredBytes*2, available)
	}

	return nil
}
