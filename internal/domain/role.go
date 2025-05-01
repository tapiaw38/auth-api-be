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
)
