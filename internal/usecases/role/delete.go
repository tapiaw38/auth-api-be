package role

import (
	"context"

	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	DeleteUsecase interface {
		Execute(context.Context, string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{
		contextFactory: contextFactory,
	}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	if id == "" {
		return apperrors.NewApplicationError(mappings.RoleDeleteIDRequiredError, nil)
	}

	role, appErr := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if appErr != nil {
		return appErr
	}
	if role == nil {
		return apperrors.NewApplicationError(mappings.RoleDeleteNotFoundError, nil)
	}

	if appErr := app.Repositories.Role.Delete(ctx, id); appErr != nil {
		return appErr
	}

	return nil
}
