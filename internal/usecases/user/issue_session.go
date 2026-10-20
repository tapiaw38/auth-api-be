package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type session struct {
	accessToken  string
	refreshToken string
}

func issueSession(ctx context.Context, app *appcontext.Context, user *domain.User) (*session, apperrors.ApplicationError) {
	accessToken, err := auth.GenerateToken(user, auth.AccessTokenTTL)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginTokenGenerationError, err)
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginTokenGenerationError, err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginTokenGenerationError, err)
	}

	if appErr := app.Repositories.RefreshToken.Create(ctx, domain.RefreshToken{
		ID:        id.String(),
		UserID:    user.ID,
		TokenHash: auth.HashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(auth.RefreshTokenTTL),
	}); appErr != nil {
		return nil, appErr
	}

	return &session{accessToken: accessToken, refreshToken: refreshToken}, nil
}
