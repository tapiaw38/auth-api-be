package user

import (
	"context"
	"errors"

	userRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	GoogleNativeLoginUsecase interface {
		Execute(context.Context, GoogleNativeLoginInput) (*GoogleNativeLoginOutput, apperrors.ApplicationError)
	}

	googleNativeLoginUsecase struct {
		contextFactory appcontext.Factory
	}

	GoogleNativeLoginInput struct {
		IDToken string `json:"id_token"`
	}

	GoogleNativeLoginOutput struct {
		Data         UserOutputData `json:"data"`
		Token        string         `json:"token"`
		RefreshToken string         `json:"refresh_token"`
	}
)

func NewGoogleNativeLoginUsecase(contextFactory appcontext.Factory) GoogleNativeLoginUsecase {
	return &googleNativeLoginUsecase{contextFactory: contextFactory}
}

func (u *googleNativeLoginUsecase) Execute(ctx context.Context, input GoogleNativeLoginInput) (*GoogleNativeLoginOutput, apperrors.ApplicationError) {
	if input.IDToken == "" {
		return nil, apperrors.NewApplicationError(mappings.UserLoginGoogleInvalidIDTokenError, errors.New("Google ID token is required"))
	}

	app := u.contextFactory()
	userInfo, err := app.Integrations.SSO.VerifyIDToken(ctx, input.IDToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginGoogleInvalidIDTokenError, err)
	}

	userID, appErr := resolveGoogleUser(ctx, app, userInfo)
	if appErr != nil {
		return nil, appErr
	}

	if userID == nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found after Google native login"))
	}

	resolvedUser, appErr := app.Repositories.User.Get(ctx, userRepo.GetFilterOptions{ID: *userID})
	if appErr != nil {
		return nil, appErr
	}

	if resolvedUser == nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found"))
	}

	issued, appErr := issueSession(ctx, app, resolvedUser)
	if appErr != nil {
		return nil, appErr
	}

	return &GoogleNativeLoginOutput{
		Data:         toUserOutputData(resolvedUser),
		Token:        issued.accessToken,
		RefreshToken: issued.refreshToken,
	}, nil
}
