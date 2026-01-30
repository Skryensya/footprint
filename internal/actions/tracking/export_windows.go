//go:build windows

package tracking

import (
	"fmt"

	"golang.org/x/sys/windows"

	"github.com/footprint-tools/cli/internal/log"
)

// checkDiskSpace verifies there's enough space to write the estimated bytes.
func checkDiskSpace(dir string, requiredBytes int64) error {
	var freeBytesAvailable, totalBytes, totalFreeBytes uint64

	dirPtr, err := windows.UTF16PtrFromString(dir)
	if err != nil {
		log.Debug("export: could not convert path %s: %v", dir, err)
		return nil
	}

	err = windows.GetDiskFreeSpaceEx(dirPtr, &freeBytesAvailable, &totalBytes, &totalFreeBytes)
	if err != nil {
		log.Debug("export: could not check disk space for %s: %v", dir, err)
		return nil
	}

	available := int64(freeBytesAvailable)

	// Require at least 2x the estimated size for safety margin
	if available < requiredBytes*2 {
		return fmt.Errorf("need %d bytes, only %d available", requiredBytes*2, available)
	}

	return nil
}
