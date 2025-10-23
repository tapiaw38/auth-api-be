package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(context.Context, GetFilterOptions) (*GetOutput, apperrors.ApplicationError)
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

func (u *getUsecase) Execute(ctx context.Context, filters GetFilterOptions) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	user, appErr := app.Repositories.User.Get(ctx, user.GetFilterOptions(filters))
	if appErr != nil {
		return nil, appErr
	}

	if user == nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetNotFoundError, nil)
	}

	return &GetOutput{
		Data: toUserOutputData(user),
	}, nil
}
