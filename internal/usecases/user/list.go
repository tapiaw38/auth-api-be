package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListFilterOptions) ([]UserOutputData, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []UserOutputData `json:"data"`
	}

	ListFilterOptions user.ListFilterOptions
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{
		contextFactory: contextFactory,
	}
}

func (u *listUsecase) Execute(ctx context.Context, filters ListFilterOptions) ([]UserOutputData, apperrors.ApplicationError) {
	app := u.contextFactory()

	users, err := app.Repositories.User.List(ctx, user.ListFilterOptions(filters))
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserListQueryError, err)
	}

	if users == nil {
		return nil, nil
	}

	outputUsers := make([]UserOutputData, 0, len(users))
	for _, user := range users {
		outputUsers = append(outputUsers, toUserOutputData(user))
	}

	return outputUsers, nil
}
