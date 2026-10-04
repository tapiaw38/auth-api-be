package user

import (
	"context"
	"errors"
	"time"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	RefreshUsecase interface {
		Execute(context.Context, RefreshInput) (*RefreshOutput, apperrors.ApplicationError)
	}

	refreshUsecase struct {
		contextFactory appcontext.Factory
	}

	RefreshInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	RefreshOutput struct {
		Data         UserOutputData `json:"data"`
		Token        string         `json:"token"`
		RefreshToken string         `json:"refresh_token"`
	}
)

func NewRefreshUsecase(contextFactory appcontext.Factory) RefreshUsecase {
	return &refreshUsecase{contextFactory: contextFactory}
}

func (u *refreshUsecase) Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if input.RefreshToken == "" {
		return nil, apperrors.NewApplicationError(mappings.AuthInvalidTokenError, errors.New("refresh token is required"))
	}

	stored, appErr := app.Repositories.RefreshToken.GetByHash(ctx, auth.HashRefreshToken(input.RefreshToken))
	if appErr != nil {
		return nil, appErr
	}

	if stored == nil {
		return nil, apperrors.NewApplicationError(mappings.AuthInvalidTokenError, errors.New("refresh token not recognised"))
	}

	if stored.RevokedAt != nil {
		if appErr := app.Repositories.RefreshToken.RevokeAllForUser(ctx, stored.UserID); appErr != nil {
			return nil, appErr
		}
		return nil, apperrors.NewApplicationError(mappings.AuthInvalidTokenError, errors.New("refresh token was already used"))
	}

	if !stored.Usable(time.Now()) {
		return nil, apperrors.NewApplicationError(mappings.AuthExpiredTokenError, errors.New("refresh token expired"))
	}

	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{ID: stored.UserID})
	if appErr != nil {
		return nil, appErr
	}

	if user == nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetNotFoundError, errors.New("user not found"))
	}

	if !user.IsActive {
		return nil, apperrors.NewApplicationError(mappings.AuthUnauthorizedError, errors.New("user is not active"))
	}

	if appErr := app.Repositories.RefreshToken.Revoke(ctx, stored.ID); appErr != nil {
		return nil, appErr
	}

	issued, appErr := issueSession(ctx, app, user)
	if appErr != nil {
		return nil, appErr
	}

	return &RefreshOutput{
		Data:         toUserOutputData(user),
		Token:        issued.accessToken,
		RefreshToken: issued.refreshToken,
	}, nil
}
