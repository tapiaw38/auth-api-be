package user_role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, userRole domain.UserRole) (*domain.UserRole, error) {
	row, err := r.executeCreateQuery(ctx, userRole)
	if err != nil {
		return nil, err
	}

	var (
		userID, roleID string
	)

	err = row.Scan(&userID, &roleID)
	if err != nil {
		return nil, err
	}

	return &domain.UserRole{
		UserID: userID,
		RoleID: roleID,
	}, nil
}

func (r *repository) executeCreateQuery(ctx context.Context, userRole domain.UserRole) (*sql.Row, error) {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) RETURNING user_id, role_id;`

	args := []any{
		userRole.UserID,
		userRole.RoleID,
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
