package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	EnsureUseCase interface {
		Execute(context.Context) apperrors.ApplicationError
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

func (e *ensureUseCase) Execute(ctx context.Context) apperrors.ApplicationError {
	app := e.contextFactory()

	roleNames := []domain.RoleName{
		domain.RoleSuperAdmin,
		domain.RoleAdmin,
		domain.RoleUser,
	}

	for _, roleName := range roleNames {
		existingRole, appErr := app.Repositories.Role.Get(
			ctx, role.GetFilterOptions{Name: string(roleName)},
		)
		if appErr != nil {
			return appErr
		}

		if existingRole != nil {
			continue
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return apperrors.NewApplicationError(mappings.InternalServerError, err)
		}
		newRole := domain.Role{
			ID:   id.String(),
			Name: roleName,
		}
		if _, appErr := app.Repositories.Role.Create(ctx, newRole); appErr != nil {
			return appErr
		}
	}

	return nil
}
