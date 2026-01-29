package output

import (
	"strings"
	"testing"
)

// mockPrintln captures output for testing
func mockPrintln() (PrintlnFunc, *strings.Builder) {
	var output strings.Builder
	fn := func(a ...any) (int, error) {
		for i, v := range a {
			if i > 0 {
				output.WriteString(" ")
			}
			output.WriteString(v.(string))
		}
		output.WriteString("\n")
		return output.Len(), nil
	}
	return fn, &output
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		expected string
	}{
		{
			name:     "simple map",
			data:     map[string]string{"key": "value"},
			expected: "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name:     "simple slice",
			data:     []string{"a", "b"},
			expected: "[\n  \"a\",\n  \"b\"\n]\n",
		},
		{
			name:     "nested struct",
			data:     struct{ Name string }{"test"},
			expected: "{\n  \"Name\": \"test\"\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println, output := mockPrintln()
			err := JSON(println, tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output.String() != tt.expected {
				t.Errorf("got %q, want %q", output.String(), tt.expected)
			}
		})
	}
}

func TestJSONLine(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		expected string
	}{
		{
			name:     "simple map",
			data:     map[string]string{"key": "value"},
			expected: "{\"key\":\"value\"}\n",
		},
		{
			name:     "simple slice",
			data:     []string{"a", "b"},
			expected: "[\"a\",\"b\"]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println, output := mockPrintln()
			err := JSONLine(println, tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output.String() != tt.expected {
				t.Errorf("got %q, want %q", output.String(), tt.expected)
			}
		})
	}
}

func TestJSONEmpty(t *testing.T) {
	println, output := mockPrintln()
	JSONEmpty(println)
	expected := "[]\n"
	if output.String() != expected {
		t.Errorf("got %q, want %q", output.String(), expected)
	}
}

func TestJSONError(t *testing.T) {
	println, output := mockPrintln()
	err := JSONError(println, "NOT_FOUND", "item not found")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Check that output contains expected fields
	out := output.String()
	if !strings.Contains(out, `"error": "NOT_FOUND"`) {
		t.Errorf("output should contain error code, got: %s", out)
	}
	if !strings.Contains(out, `"message": "item not found"`) {
		t.Errorf("output should contain message, got: %s", out)
	}
}

func TestJSONMarshalError(t *testing.T) {
	println, _ := mockPrintln()
	// channels cannot be marshaled to JSON
	err := JSON(println, make(chan int))
	if err == nil {
		t.Error("expected error for unmarshallable type")
	}
	if !strings.Contains(err.Error(), "failed to marshal JSON") {
		t.Errorf("error should mention marshal failure, got: %v", err)
	}
}

func TestJSONLineMarshalError(t *testing.T) {
	println, _ := mockPrintln()
	// channels cannot be marshaled to JSON
	err := JSONLine(println, make(chan int))
	if err == nil {
		t.Error("expected error for unmarshallable type")
	}
	if !strings.Contains(err.Error(), "failed to marshal JSON") {
		t.Errorf("error should mention marshal failure, got: %v", err)
	}
}
