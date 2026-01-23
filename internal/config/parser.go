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

		// Skip array keys (key[]=value) - handled by ParseArray
		if strings.HasSuffix(key, "[]") {
			continue
		}

		// Strip surrounding quotes if present
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		cfg[key] = value
	}

	return cfg, nil
}

// ParseArray extracts all values for an array key (key[]=value format).
// Returns values in order of appearance.
func ParseArray(lines []string, arrayKey string) []string {
	var values []string
	prefix := arrayKey + "[]="

	for i, line := range lines {
		if i == 0 {
			line = strings.TrimPrefix(line, "\uFEFF") // BOM safety
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, prefix) {
			value := strings.TrimPrefix(trimmed, prefix)
			value = strings.TrimSpace(value)
			if value != "" {
				values = append(values, value)
			}
		}
	}

	return values
}
