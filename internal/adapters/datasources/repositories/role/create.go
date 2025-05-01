package role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, role domain.Role) (string, error) {
	row, err := r.executeCreateQuery(ctx, role)
	if err != nil {
		return "", err
	}

	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) executeCreateQuery(ctx context.Context, role domain.Role) (*sql.Row, error) {
	query := `INSERT INTO roles (id, name) VALUES ($1, $2) RETURNING id;`

	args := []any{
		role.ID,
		role.Name,
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
