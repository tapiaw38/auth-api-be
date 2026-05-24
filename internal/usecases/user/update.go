package user

import (
	"context"
	"errors"
	"time"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
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
		ID            string
		AuthUsername  string
		AuthRoles     []auth.RoleClaim
		FirstName     string
		LastName      string
		Email         string
		Picture       *string
		PhoneNumber   *string
		Address       *string
		IsActive      *bool
		VerifiedEmail *bool
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

	if authUser.ID != input.ID && !hasAdminRole(input.AuthRoles) {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateUnauthorizedError,
			errors.New("user attempted to update another user without admin role"))
	}

	payload := &domain.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Picture:     input.Picture,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		UpdatedAt:   time.Now(),
	}
	if input.IsActive != nil {
		payload.IsActive = *input.IsActive
	}
	if input.VerifiedEmail != nil {
		payload.VerifiedEmail = *input.VerifiedEmail
	}

	updatedID, repoErr := app.Repositories.User.Update(ctx, input.ID, payload)
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

func hasAdminRole(roles []auth.RoleClaim) bool {
	for _, r := range roles {
		if r.Name == string(domain.RoleSuperAdmin) || r.Name == string(domain.RoleAdmin) {
			return true
		}
	}
	return false
}
