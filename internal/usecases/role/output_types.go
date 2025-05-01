package role

import "github.com/tapiaw38/auth-api-be/internal/domain"

type (
	RoleOutputData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func toRoleOutputData(role domain.Role) RoleOutputData {
	return RoleOutputData{
		ID:   role.ID,
		Name: string(role.Name),
	}
}
