package bencode

import (
	"testing"
)

func TestExtractInfoBytes(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    []byte
		expectError string
	}{
		{
			name:     "extracts info value when it is a string",
			input:    []byte("d4:info5:helloe"),
			expected: []byte("5:hello"),
		},
		{
			name:     "extracts info value when it is an integer",
			input:    []byte("d4:infoi42ee"),
			expected: []byte("i42e"),
		},
		{
			name:     "extracts info value when it is a dictionary",
			input:    []byte("d4:infod3:foo3:baree"),
			expected: []byte("d3:foo3:bare"),
		},
		{
			name:     "extracts info value when it is a list",
			input:    []byte("d4:infoli1ei2eee"),
			expected: []byte("li1ei2ee"),
		},
		{
			name:     "extracts info when it is not the first key",
			input:    []byte("d4:name5:hello4:infoi99ee"),
			expected: []byte("i99e"),
		},
		{
			name:     "extracts info from a dict with multiple keys after it",
			input:    []byte("d4:infoi7e5:other3:fooe"),
			expected: []byte("i7e"),
		},
		{
			name:        "returns error for empty input",
			input:       []byte(""),
			expectError: "invalid bencoded data",
		},
		{
			name:        "returns error when input is not a dictionary",
			input:       []byte("5:hello"),
			expectError: "invalid bencoded data",
		},
		{
			name:        "returns error when info key is absent",
			input:       []byte("d3:foo3:bare"),
			expectError: "info section not found",
		},
		{
			name:        "returns error for non-string key",
			input:       []byte("di42e3:fooe"),
			expectError: "expected string key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractInfoBytes(tt.input)

			if tt.expectError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectError)
				}
				if !containsSubstring(err.Error(), tt.expectError) {
					t.Fatalf("expected error containing %q, got %q", tt.expectError, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != string(tt.expected) {
				t.Fatalf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
