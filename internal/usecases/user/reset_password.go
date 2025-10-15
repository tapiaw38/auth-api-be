package user

import (
	"context"
	"errors"
	"time"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

type (
	ResetPasswordUsecase interface {
		Execute(context.Context, string, string) (*ResetPasswordOutput, error)
	}

	resetPasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	ResetPasswordOutput struct {
		Data ResetPasswordOutputData `json:"data"`
	}
)

func NewResetPasswordUsecase(contextFactory appcontext.Factory) ResetPasswordUsecase {
	return &resetPasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *resetPasswordUsecase) Execute(ctx context.Context, token, password string) (*ResetPasswordOutput, error) {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(
		ctx,
		user_repo.GetFilterOptions{
			PasswordResetToken: token,
		},
	)
	if err != nil {
		return nil, errors.New("token expired or invalid")
	}

	if user == nil {
		return nil, errors.New("token expired or invalid")
	}

	if time.Now().After(*user.PasswordResetTokenExpiry) {
		return nil, errors.New("token expired or invalid")
	}

	if err := auth.ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	hashedPassword, err := auth.HashedPassword(password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	if _, err = app.Repositories.User.Update(
		ctx,
		user.ID,
		user,
	); err != nil {
		return nil, err
	}

	if err := app.Repositories.User.InvalidatePasswordResetToken(ctx, user.ID); err != nil {
		return nil, err
	}

	return &ResetPasswordOutput{
		Data: ResetPasswordOutputData{
			Email:   user.Email,
			Message: "Password reseted successfully",
		},
	}, nil
}
