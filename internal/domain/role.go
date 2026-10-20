package domain

const (
	RoleSuperAdmin RoleName = "superadmin"
	RoleAdmin      RoleName = "admin"
	RoleUser       RoleName = "user"
)

type (
	RoleName string

	Role struct {
		ID   string
		Name RoleName
	}

	roleClaimDescriptor interface {
		GetID() string
		GetName() string
	}
)

func CanManageUsers(roles []Role) bool {
	for _, role := range roles {
		if role.Name == RoleSuperAdmin {
			return true
		}
	}

	return false
}

func RolesFromClaims[T roleClaimDescriptor](roleClaims []T) []Role {
	roles := make([]Role, 0, len(roleClaims))
	for _, roleClaim := range roleClaims {
		roles = append(roles, Role{
			ID:   roleClaim.GetID(),
			Name: RoleName(roleClaim.GetName()),
		})
	}

	return roles
}
