package user

import (
	"context"
	"time"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	ResetPasswordUsecase interface {
		Execute(context.Context, ResetPasswordInput) (*ResetPasswordOutput, apperrors.ApplicationError)
	}

	resetPasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	ResetPasswordOutput struct {
		Data ResetPasswordOutputData `json:"data"`
	}

	ResetPasswordInput struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
)

func NewResetPasswordUsecase(contextFactory appcontext.Factory) ResetPasswordUsecase {
	return &resetPasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *resetPasswordUsecase) Execute(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	user, appErr := app.Repositories.User.Get(
		ctx,
		user_repo.GetFilterOptions{
			PasswordResetToken: input.Token,
		},
	)
	if appErr != nil {
		return nil, apperrors.NewApplicationError(mappings.UserResetPasswordInvalidTokenError, nil)
	}

	if user == nil {
		return nil, apperrors.NewApplicationError(mappings.UserResetPasswordInvalidTokenError, nil)
	}

	if time.Now().After(*user.PasswordResetTokenExpiry) {
		return nil, apperrors.NewApplicationError(mappings.UserResetPasswordInvalidTokenError, nil)
	}

	if err := auth.ValidatePasswordStrength(input.Password); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterWeakPasswordError, err)
	}

	hashedPassword, err := auth.HashedPassword(input.Password)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterPasswordHashError, err)
	}

	user.Password = string(hashedPassword)

	if _, appErr := app.Repositories.User.Patch(
		ctx,
		user.ID,
		user,
	); appErr != nil {
		return nil, appErr
	}

	if appErr := app.Repositories.User.InvalidatePasswordResetToken(ctx, user.ID); appErr != nil {
		return nil, appErr
	}

	if appErr := app.Repositories.User.IncrementTokenVersion(ctx, user.ID); appErr != nil {
		return nil, appErr
	}

	return &ResetPasswordOutput{
		Data: ResetPasswordOutputData{
			Email:   user.Email,
			Message: "Password reseted successfully",
		},
	}, nil
}
