package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := map[string]struct {
		password    string
		expectedErr string
	}{
		"valid password with all requirements": {
			password:    "Password123!",
			expectedErr: "",
		},
		"valid password with different special char": {
			password:    "MyP@ssw0rd",
			expectedErr: "",
		},
		"valid password complex": {
			password:    "C0mpl3x&Pass",
			expectedErr: "",
		},
		"password too short": {
			password:    "Pass1!",
			expectedErr: "password must be at least 8 characters long",
		},
		"password without uppercase": {
			password:    "password123!",
			expectedErr: "password must contain at least one uppercase letter",
		},
		"password without lowercase": {
			password:    "PASSWORD123!",
			expectedErr: "password must contain at least one lowercase letter",
		},
		"password without number": {
			password:    "Password!",
			expectedErr: "password must contain at least one number",
		},
		"password without special character": {
			password:    "Password123",
			expectedErr: "password must contain at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)",
		},
		"empty password": {
			password:    "",
			expectedErr: "password must be at least 8 characters long",
		},
		"password with only letters": {
			password:    "PasswordOnly",
			expectedErr: "password must contain at least one number",
		},
		"password with only numbers": {
			password:    "12345678",
			expectedErr: "password must contain at least one uppercase letter",
		},
		"password with valid length but missing multiple requirements": {
			password:    "password",
			expectedErr: "password must contain at least one uppercase letter",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := auth.ValidatePasswordStrength(tc.password)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
			}
		})
	}
}
