package role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) List(ctx context.Context, filters ListFilterOptions) ([]domain.Role, apperrors.ApplicationError) {
	rows, err := r.executeListQuery(ctx, filters)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.RoleListQueryError, err)
	}

	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var (
			id, name string
		)

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.RoleListQueryError, err)
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
