package user

import (
	"context"

	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) ChangePassword(ctx context.Context, id string, password string) apperrors.ApplicationError {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET password = $1 WHERE id = $2;")
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserChangePasswordUpdateError, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, password, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserChangePasswordUpdateError, err)
	}

	return nil
}
