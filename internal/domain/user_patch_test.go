package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func TestUserApplyPatch(t *testing.T) {
	firstName := " John "
	lastName := " Doe "
	email := " john@example.com "
	picture := " https://example.com/avatar.jpg "
	phoneNumber := " +54 9 383 1234567 "
	address := " 123 Main St "
	emptyString := ""
	invalidEmail := "not-an-email"
	active := false
	verified := true

	tests := map[string]struct {
		user        domain.User
		patch       domain.UserPatch
		expected    domain.User
		expectedErr string
	}{
		"apply valid patch trims and updates fields": {
			user: domain.User{
				FirstName:     "Jane",
				LastName:      "Smith",
				Email:         "jane@example.com",
				IsActive:      true,
				VerifiedEmail: false,
			},
			patch: domain.UserPatch{
				FirstName: domain.Optional[string]{Set: true, Value: &firstName},
				LastName:  domain.Optional[string]{Set: true, Value: &lastName},
				Email:     domain.Optional[string]{Set: true, Value: &email},
				Picture:   domain.Optional[string]{Set: true, Value: &picture},
				PhoneNumber: domain.Optional[string]{
					Set:   true,
					Value: &phoneNumber,
				},
				Address:       domain.Optional[string]{Set: true, Value: &address},
				IsActive:      domain.Optional[bool]{Set: true, Value: &active},
				VerifiedEmail: domain.Optional[bool]{Set: true, Value: &verified},
			},
			expected: domain.User{
				FirstName:     "John",
				LastName:      "Doe",
				Email:         "john@example.com",
				Picture:       strPtr("https://example.com/avatar.jpg"),
				PhoneNumber:   strPtr("+54 9 383 1234567"),
				Address:       strPtr("123 Main St"),
				IsActive:      false,
				VerifiedEmail: true,
			},
		},
		"clear nullable fields with null": {
			user: domain.User{
				Picture:     strPtr("https://example.com/avatar.jpg"),
				PhoneNumber: strPtr("+54 9 383 1234567"),
				Address:     strPtr("123 Main St"),
			},
			patch: domain.UserPatch{
				Picture:     domain.Optional[string]{Set: true, Value: nil},
				PhoneNumber: domain.Optional[string]{Set: true, Value: nil},
				Address:     domain.Optional[string]{Set: true, Value: nil},
			},
			expected: domain.User{
				Picture:     nil,
				PhoneNumber: nil,
				Address:     nil,
			},
		},
		"reject empty patch": {
			user:        domain.User{},
			patch:       domain.UserPatch{},
			expectedErr: "at least one field is required",
		},
		"reject null first name": {
			user: domain.User{},
			patch: domain.UserPatch{
				FirstName: domain.Optional[string]{Set: true, Value: nil},
			},
			expectedErr: "first_name cannot be null",
		},
		"reject empty nullable field": {
			user: domain.User{},
			patch: domain.UserPatch{
				Picture: domain.Optional[string]{Set: true, Value: &emptyString},
			},
			expectedErr: "picture cannot be empty; use null to clear it",
		},
		"reject invalid email": {
			user: domain.User{},
			patch: domain.UserPatch{
				Email: domain.Optional[string]{Set: true, Value: &invalidEmail},
			},
			expectedErr: "invalid email format",
		},
		"reject null bool": {
			user: domain.User{},
			patch: domain.UserPatch{
				IsActive: domain.Optional[bool]{Set: true, Value: nil},
			},
			expectedErr: "is_active cannot be null",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			user := tc.user

			err := user.ApplyPatch(tc.patch)

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, user)
		})
	}
}

func TestUserPatchHasChanges(t *testing.T) {
	value := "John"
	flag := true

	tests := map[string]struct {
		patch    domain.UserPatch
		expected bool
	}{
		"empty patch has no changes": {
			patch:    domain.UserPatch{},
			expected: false,
		},
		"string field counts as change": {
			patch: domain.UserPatch{
				FirstName: domain.Optional[string]{Set: true, Value: &value},
			},
			expected: true,
		},
		"null string field counts as change": {
			patch: domain.UserPatch{
				Address: domain.Optional[string]{Set: true, Value: nil},
			},
			expected: true,
		},
		"bool field counts as change": {
			patch: domain.UserPatch{
				IsActive: domain.Optional[bool]{Set: true, Value: &flag},
			},
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.patch.HasChanges())
		})
	}
}

func strPtr(value string) *string {
	return &value
}
