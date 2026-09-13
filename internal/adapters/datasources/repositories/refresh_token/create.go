package refresh_token

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) Create(ctx context.Context, token domain.RefreshToken) apperrors.ApplicationError {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) VALUES ($1, $2, $3, $4)`

	if _, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt); err != nil {
		return apperrors.NewApplicationError(mappings.DatabaseQueryError, err)
	}

	return nil
}
