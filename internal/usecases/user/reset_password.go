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
		Execute(context.Context, ResetPasswordInput) (*ResetPasswordOutput, error)
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

func (u *resetPasswordUsecase) Execute(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, error) {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(
		ctx,
		user_repo.GetFilterOptions{
			PasswordResetToken: input.Token,
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

	if err := auth.ValidatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := auth.HashedPassword(input.Password)
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
