package role

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListFilterOptions) ([]RoleOutputData, apperrors.ApplicationError)
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

func (u *listUsecase) Execute(ctx context.Context, filters ListFilterOptions) ([]RoleOutputData, apperrors.ApplicationError) {
	app := u.contextFactory()

	roles, appErr := app.Repositories.Role.List(ctx, role.ListFilterOptions(filters))
	if appErr != nil {
		return nil, appErr
	}

	if roles == nil {
		return nil, nil
	}

	outputRoles := make([]RoleOutputData, 0, len(roles))
	for _, role := range roles {
		outputRoles = append(outputRoles, toRoleOutputData(role))
	}

	return outputRoles, nil
}
