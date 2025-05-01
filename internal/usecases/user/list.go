package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListFilterOptions) ([]UserOutputData, error)
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

func (u *listUsecase) Execute(ctx context.Context, filters ListFilterOptions) ([]UserOutputData, error) {
	app := u.contextFactory()

	users, err := app.Repositories.User.List(ctx, user.ListFilterOptions(filters))
	if err != nil {
		return nil, err
	}

	outputUsers := make([]UserOutputData, 0, len(users))
	for _, user := range users {
		outputUsers = append(outputUsers, toUserOutputData(user))
	}

	return outputUsers, nil
}
