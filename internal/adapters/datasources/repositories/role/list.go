package role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) List(ctx context.Context, filters ListFilterOptions) ([]domain.Role, error) {
	rows, err := r.executeListQuery(ctx, filters)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var (
			id, name string
		)

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		roles = append(roles, domain.Role{
			ID:   id,
			Name: domain.RoleName(name),
		})
	}

	return roles, nil
}

func (r *repository) executeListQuery(ctx context.Context, filters ListFilterOptions) (*sql.Rows, error) {
	query := `SELECT id, name FROM roles`

	query += ` WHERE 1=1 `

	args := []any{}

	if filters.Name != "" {
		query += ` AND name = $1`
		args = append(args, filters.Name)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
