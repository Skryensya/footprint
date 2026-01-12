package config

import (
	"fmt"
	"strings"
)

func Parse(lines []string) (map[string]string, error) {
	cfg := make(map[string]string)

	for i, line := range lines {
		if i == 0 {
			line = strings.TrimPrefix(line, "\uFEFF") // BOM safety
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config format at line %d", i+1)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("invalid empty key at line %d", i+1)
		}

		cfg[key] = value
	}

	return cfg, nil
}
