package user

import (
	"context"
	"errors"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(context.Context, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		ID             string
		AuthUsername   string
		CanManageUsers bool
		Patch          domain.UserPatch
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

func (u *updateUsecase) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	authUser, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Username: input.AuthUsername,
	})
	if appErr != nil {
		return nil, appErr
	}
	if authUser == nil {
		return nil, apperrors.NewApplicationError(mappings.AuthUnauthorizedError,
			errors.New("authenticated user not found"))
	}

	if authUser.ID != input.ID && !input.CanManageUsers {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateUnauthorizedError,
			errors.New("user attempted to update another user without admin role"))
	}

	if err := targetPatchValidation(input.Patch); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateInvalidInputError, err)
	}

	if !input.CanManageUsers && hasRestrictedFieldUpdates(input) {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateRestrictedFieldsError,
			errors.New("user attempted to update restricted fields"))
	}

	targetUser := authUser
	if authUser.ID != input.ID {
		targetUser, appErr = app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
			ID: input.ID,
		})
		if appErr != nil {
			return nil, appErr
		}
	}
	if targetUser == nil {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateNotFoundError, nil)
	}

	payload := *targetUser
	if err := payload.ApplyPatch(input.Patch); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateInvalidInputError, err)
	}

	updatedID, repoErr := app.Repositories.User.Patch(ctx, input.ID, &payload)
	if repoErr != nil {
		return nil, repoErr
	}

	updatedUser, repoErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: updatedID,
	})
	if repoErr != nil {
		return nil, repoErr
	}

	return &UpdateOutput{
		Data: toUserOutputData(updatedUser),
	}, nil
}

func hasRestrictedFieldUpdates(input UpdateInput) bool {
	return input.Patch.IsActive.Set || input.Patch.VerifiedEmail.Set
}

func targetPatchValidation(patch domain.UserPatch) error { return (&domain.User{}).ApplyPatch(patch) }
