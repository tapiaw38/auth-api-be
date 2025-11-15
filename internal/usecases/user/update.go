package user

import (
	"context"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(context.Context, string, *domain.User) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateOutput struct {
		Data UserOutputData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{
		contextFactory: contextFactory,
	}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, user *domain.User) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	updatedID, err := app.Repositories.User.Update(ctx, id, user)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateQueryError, err)
	}

	updatedUser, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: updatedID,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, err)
	}

	return &UpdateOutput{
		Data: toUserOutputData(updatedUser),
	}, nil
}
