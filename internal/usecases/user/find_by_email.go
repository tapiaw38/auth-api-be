package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	FindByEmailUsecase interface {
		Execute(context.Context, string) (*FindByEmailOutput, apperrors.ApplicationError)
	}

	findByEmailUsecase struct {
		contextFactory appcontext.Factory
	}

	FindByEmailOutput struct {
		Data SummaryOutputData `json:"data"`
	}
)

func NewFindByEmailUsecase(contextFactory appcontext.Factory) FindByEmailUsecase {
	return &findByEmailUsecase{contextFactory: contextFactory}
}

func (u *findByEmailUsecase) Execute(ctx context.Context, email string) (*FindByEmailOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	found, appErr := app.Repositories.User.Get(ctx, user.GetFilterOptions{Email: email})
	if appErr != nil {
		return nil, appErr
	}
	if found == nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetNotFoundError, nil)
	}

	return &FindByEmailOutput{Data: toSummaryOutputData(found)}, nil
}
