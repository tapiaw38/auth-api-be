package user_role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Delete(ctx context.Context, userID, roleID string) (*domain.UserRole, error) {
	result, err := r.executeDeleteQuery(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}

	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}

	return &domain.UserRole{
		UserID: userID,
		RoleID: roleID,
	}, nil
}

func (r *repository) executeDeleteQuery(ctx context.Context, userID, roleID string) (sql.Result, error) {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	args := []any{
		userID,
		roleID,
	}

	return r.db.ExecContext(ctx, query, args...)
}
