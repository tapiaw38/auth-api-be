package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	GetTokenVersionUsecase interface {
		Execute(context.Context, string) (uint, error)
	}

	getTokenVersionUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewGetTokenVersionUsecase(contextFactory appcontext.Factory) GetTokenVersionUsecase {
	return &getTokenVersionUsecase{
		contextFactory: contextFactory,
	}
}

func (u *getTokenVersionUsecase) Execute(ctx context.Context, username string) (uint, error) {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(ctx, user.GetFilterOptions{
		Username: username,
	})
	if err != nil {
		return 0, err
	}

	if user == nil {
		return 0, nil
	}

	return user.TokenVersion, nil
}
