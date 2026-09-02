package user

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

const maxBatchIDs = 200

type (
	BatchUsecase interface {
		Execute(context.Context, BatchInput) (*BatchOutput, apperrors.ApplicationError)
	}

	batchUsecase struct {
		contextFactory appcontext.Factory
	}

	BatchInput struct {
		IDs []string
	}

	BatchOutput struct {
		Data []SummaryOutputData `json:"data"`
	}
)

func NewBatchUsecase(contextFactory appcontext.Factory) BatchUsecase {
	return &batchUsecase{
		contextFactory: contextFactory,
	}
}

func (u *batchUsecase) Execute(ctx context.Context, input BatchInput) (*BatchOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if len(input.IDs) == 0 {
		return &BatchOutput{Data: []SummaryOutputData{}}, nil
	}

	ids := input.IDs
	if len(ids) > maxBatchIDs {
		ids = ids[:maxBatchIDs]
	}

	users, err := app.Repositories.User.List(ctx, user.ListFilterOptions{IDs: ids})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserListQueryError, err)
	}

	data := make([]SummaryOutputData, 0, len(users))
	for _, u := range users {
		data = append(data, toSummaryOutputData(u))
	}

	return &BatchOutput{Data: data}, nil
}
