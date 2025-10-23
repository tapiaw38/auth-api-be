package user

import (
	"context"

	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) InvalidatePasswordResetToken(ctx context.Context, id string) apperrors.ApplicationError {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET password_reset_token = NULL, password_reset_token_expiry = NULL WHERE id = $1;")
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserResetPasswordUpdateError, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserResetPasswordUpdateError, err)
	}

	return nil
}
