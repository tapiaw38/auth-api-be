package refresh_token

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
)

type (
	Repository interface {
		Create(context.Context, domain.RefreshToken) apperrors.ApplicationError
		GetByHash(context.Context, string) (*domain.RefreshToken, apperrors.ApplicationError)
		Revoke(context.Context, string) apperrors.ApplicationError
		RevokeAllForUser(context.Context, string) apperrors.ApplicationError
	}

	repository struct {
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
