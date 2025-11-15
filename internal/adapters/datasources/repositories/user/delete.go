package user

import (
	"context"
	"database/sql"

	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) Delete(ctx context.Context, id string) apperrors.ApplicationError {
	result, err := r.executeDeleteQuery(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserDeleteQueryError, err)
	}

	if _, err := result.RowsAffected(); err != nil {
		return apperrors.NewApplicationError(mappings.UserDeleteQueryError, err)
	}

	return nil
}

func (r *repository) executeDeleteQuery(ctx context.Context, id string) (sql.Result, error) {
	query := `DELETE FROM users WHERE id = $1`

	args := []any{
		id,
	}

	return r.db.ExecContext(ctx, query, args...)
}
