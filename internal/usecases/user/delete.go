package user

import (
	"context"

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
		return apperrors.NewApplicationError(mappings.UserDeleteNotFoundError, nil)
	}

	appErr := app.Repositories.User.Delete(ctx, id)
	if appErr != nil {
		return appErr
	}

	return nil
}
