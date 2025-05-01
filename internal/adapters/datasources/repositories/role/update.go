package role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, role *domain.Role) (string, error) {
	row, err := r.executeUpdateQuery(ctx, id, role)
	if err != nil {
		return "", err
	}

	var updatedID string
	if err := row.Scan(&updatedID); err != nil {
		return "", err
	}

	return updatedID, nil
}

func (r *repository) executeUpdateQuery(ctx context.Context, id string, role *domain.Role) (*sql.Row, error) {
	query := `UPDATE roles
		SET 
			name = COALESCE($1, name)
		WHERE id = $2
		RETURNING id;`

	args := []any{
		role.Name,
		id,
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
