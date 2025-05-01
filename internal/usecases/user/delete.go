package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	DeleteUsecase interface {
		Execute(context.Context, string) (string, error)
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

func (u *deleteUsecase) Execute(ctx context.Context, id string) (string, error) {
	app := u.contextFactory()

	err := app.Repositories.User.Delete(ctx, id)
	if err != nil {
		return "", err
	}

	return id, nil
}
