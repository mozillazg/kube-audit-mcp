package alibaba

import (
	"testing"
)

func TestGetSLSFilterExp(t *testing.T) {
	tests := []struct {
		name     string
		keyword  string
		expected string
	}{
		{
			name:     "wildcard suffix",
			keyword:  "test*",
			expected: "test*",
		},
		{
			name:     "exact match",
			keyword:  "test",
			expected: `"test"`,
		},
		{
			name:     "empty string",
			keyword:  "",
			expected: `""`,
		},
		{
			name:     "special characters",
			keyword:  "test-123",
			expected: `"test-123"`,
		},
		{
			name:     "multiple wildcards but not at suffix",
			keyword:  "*test",
			expected: `"*test"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSLSFilterExp(tt.keyword)
			if got != tt.expected {
				t.Errorf("getSLSFilterExp(%q) = %q, want %q", tt.keyword, got, tt.expected)
			}
		})
	}
}
