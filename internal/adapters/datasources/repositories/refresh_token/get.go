package refresh_token

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, apperrors.ApplicationError) {
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE token_hash = $1`

	var token domain.RefreshToken
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.DatabaseQueryError, err)
	}

	return &token, nil
}
