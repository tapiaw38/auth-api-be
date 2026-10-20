package user

import (
	"context"

	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) IncrementTokenVersion(ctx context.Context, id string) apperrors.ApplicationError {
	query := `UPDATE users SET token_version = token_version + 1, updated_at = NOW() WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return apperrors.NewApplicationError(mappings.UserUpdateQueryError, err)
	}

	return nil
}
