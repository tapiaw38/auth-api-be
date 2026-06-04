package web_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	platformweb "github.com/tapiaw38/auth-api-be/internal/platform/web"
)

func TestRequestOptionalStringUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		input       string
		expectedSet bool
		expected    *string
		expectErr   bool
	}{
		"when valid string is passed": {
			input:       `"John"`,
			expectedSet: true,
			expected:    strPtr("John"),
			expectErr:   false,
		},
		"when null is passed": {
			input:       `null`,
			expectedSet: true,
			expected:    nil,
			expectErr:   false,
		},
		"when invalid json is passed": {
			input:       `123`,
			expectedSet: true,
			expected:    nil,
			expectErr:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var result platformweb.RequestOptional[string]
			err := json.Unmarshal([]byte(tc.input), &result)

			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.expectedSet, result.Set)
			assert.Equal(t, tc.expected, result.Value)
		})
	}
}

func TestRequestOptionalBoolUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		input       string
		expectedSet bool
		expected    *bool
		expectErr   bool
	}{
		"when valid bool is passed": {
			input:       `true`,
			expectedSet: true,
			expected:    boolPtr(true),
			expectErr:   false,
		},
		"when null is passed": {
			input:       `null`,
			expectedSet: true,
			expected:    nil,
			expectErr:   false,
		},
		"when invalid json is passed": {
			input:       `"true"`,
			expectedSet: true,
			expected:    nil,
			expectErr:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var result platformweb.RequestOptional[bool]
			err := json.Unmarshal([]byte(tc.input), &result)

			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.expectedSet, result.Set)
			assert.Equal(t, tc.expected, result.Value)
		})
	}
}

func strPtr(value string) *string {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}
