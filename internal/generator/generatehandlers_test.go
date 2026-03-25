package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMediaType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain application/json",
			input:    "application/json",
			expected: "application/json",
		},
		{
			name:     "application/json with charset",
			input:    "application/json; charset=utf-8",
			expected: "application/json",
		},
		{
			name:     "comma-separated content types",
			input:    "application/json; charset=utf-8, application/json",
			expected: "application/json",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "text/plain with charset",
			input:    "text/plain; charset=utf-8",
			expected: "text/plain",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseMediaType(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
