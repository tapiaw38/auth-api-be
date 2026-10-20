package refresh_token

import (
	"context"

	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) Revoke(ctx context.Context, id string) apperrors.ApplicationError {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return apperrors.NewApplicationError(mappings.DatabaseQueryError, err)
	}

	return nil
}

func (r *repository) RevokeAllForUser(ctx context.Context, userID string) apperrors.ApplicationError {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return apperrors.NewApplicationError(mappings.DatabaseQueryError, err)
	}

	return nil
}
