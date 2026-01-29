package hooks

import (
	"os"
)

func Status(hooksPath string) map[string]bool {
	out := make(map[string]bool, len(ManagedHooks))

	// Initialize all hooks as false
	for _, hook := range ManagedHooks {
		out[hook] = false
	}

	// Read directory once instead of calling os.Stat for each hook
	entries, err := os.ReadDir(hooksPath)
	if err != nil {
		return out
	}

	// Build set of existing files
	existing := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			existing[entry.Name()] = struct{}{}
		}
	}

	// Check which managed hooks exist
	for _, hook := range ManagedHooks {
		_, out[hook] = existing[hook]
	}

	return out
}
