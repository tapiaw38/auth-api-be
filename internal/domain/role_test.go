package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

func TestCanManageUsers(t *testing.T) {
	tests := map[string]struct {
		roles    []domain.Role
		expected bool
	}{
		"superadmin can manage users": {
			roles: []domain.Role{
				{ID: "role-1", Name: domain.RoleSuperAdmin},
			},
			expected: true,
		},
		"admin can manage users": {
			roles: []domain.Role{
				{ID: "role-1", Name: domain.RoleAdmin},
			},
			expected: true,
		},
		"user cannot manage users": {
			roles: []domain.Role{
				{ID: "role-1", Name: domain.RoleUser},
			},
			expected: false,
		},
		"mixed roles including admin can manage users": {
			roles: []domain.Role{
				{ID: "role-1", Name: domain.RoleUser},
				{ID: "role-2", Name: domain.RoleAdmin},
			},
			expected: true,
		},
		"empty roles cannot manage users": {
			roles:    []domain.Role{},
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, domain.CanManageUsers(tc.roles))
		})
	}
}

func TestRolesFromClaims(t *testing.T) {
	tests := map[string]struct {
		roleClaims []auth.RoleClaim
		expected   []domain.Role
	}{
		"maps claims to roles": {
			roleClaims: []auth.RoleClaim{
				{ID: "role-1", Name: "user"},
				{ID: "role-2", Name: "admin"},
			},
			expected: []domain.Role{
				{ID: "role-1", Name: domain.RoleUser},
				{ID: "role-2", Name: domain.RoleAdmin},
			},
		},
		"empty claims returns empty roles": {
			roleClaims: []auth.RoleClaim{},
			expected:   []domain.Role{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, domain.RolesFromClaims(tc.roleClaims))
		})
	}
}
