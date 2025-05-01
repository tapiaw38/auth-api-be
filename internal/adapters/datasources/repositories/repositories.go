package repositories

import (
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	user_role "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type Repositories struct {
	User     user.Repository
	Role     role.Repository
	UserRole user_role.Repository
}

type Factory func() *Repositories

func NewFactory(
	datasources *datasources.Datasources,
	configService *config.ConfigurationService,
) func() *Repositories {
	return func() *Repositories {
		return &Repositories{
			User:     user.NewRepository(datasources.DB),
			Role:     role.NewRepository(datasources.DB),
			UserRole: user_role.NewRepository(datasources.DB),
		}
	}
}
