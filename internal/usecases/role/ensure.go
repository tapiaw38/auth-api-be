package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	EnsureUseCase interface {
		Execute(context.Context) error
	}

	ensureUseCase struct {
		contextFactory appcontext.Factory
	}
)

func NewEnsureUseCase(contextFactory appcontext.Factory) EnsureUseCase {
	return &ensureUseCase{
		contextFactory: contextFactory,
	}
}

func (e *ensureUseCase) Execute(ctx context.Context) error {
	app := e.contextFactory()

	roleNames := []domain.RoleName{
		domain.RoleSuperAdmin,
		domain.RoleAdmin,
		domain.RoleUser,
	}

	for _, roleName := range roleNames {
		existingRole, err := app.Repositories.Role.Get(
			ctx, role.GetFilterOptions{Name: string(roleName)},
		)
		if err != nil {
			return err
		}

		if existingRole != nil {
			continue
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		newRole := domain.Role{
			ID:   id.String(),
			Name: roleName,
		}
		if _, err := app.Repositories.Role.Create(ctx, newRole); err != nil {
			return err
		}
	}

	return nil
}
