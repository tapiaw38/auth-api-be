package role

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListFilterOptions) ([]RoleOutputData, error)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []RoleOutputData `json:"data"`
	}

	ListFilterOptions role.ListFilterOptions
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{
		contextFactory: contextFactory,
	}
}

func (u *listUsecase) Execute(ctx context.Context, filters ListFilterOptions) ([]RoleOutputData, error) {
	app := u.contextFactory()

	roles, err := app.Repositories.Role.List(ctx, role.ListFilterOptions(filters))
	if err != nil {
		return nil, err
	}

	outputRoles := make([]RoleOutputData, 0, len(roles))
	for _, role := range roles {
		outputRoles = append(outputRoles, toRoleOutputData(role))
	}

	return outputRoles, nil
}
