package output

import (
	"encoding/json"
	"fmt"
)

// PrintlnFunc is the function signature for printing with a newline
type PrintlnFunc func(a ...any) (int, error)

// JSON marshals data with 2-space indentation and prints it
func JSON(println PrintlnFunc, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, _ = println(string(bytes))
	return nil
}

// JSONLine marshals data without indentation (for newline-delimited JSON)
func JSONLine(println PrintlnFunc, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, _ = println(string(bytes))
	return nil
}

// JSONEmpty prints an empty JSON array
func JSONEmpty(println PrintlnFunc) {
	_, _ = println("[]")
}

// JSONError prints a standardized error response
func JSONError(println PrintlnFunc, code string, message string) error {
	return JSON(println, map[string]string{
		"error":   code,
		"message": message,
	})
}
