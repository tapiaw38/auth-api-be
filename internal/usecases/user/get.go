package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	GetUsecase interface {
		Execute(context.Context, GetFilterOptions) (*GetOutput, error)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data UserOutputData `json:"data"`
	}

	GetFilterOptions user.GetFilterOptions
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{
		contextFactory: contextFactory,
	}
}

func (u *getUsecase) Execute(ctx context.Context, filters GetFilterOptions) (*GetOutput, error) {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(ctx, user.GetFilterOptions(filters))
	if err != nil {
		return nil, err
	}

	return &GetOutput{
		Data: toUserOutputData(user),
	}, nil
}
