package role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) Get(ctx context.Context, filters GetFilterOptions) (*domain.Role, apperrors.ApplicationError) {
	row, err := r.executeGetQuery(ctx, filters)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.RoleGetQueryError, err)
	}

	var (
		id, name string
	)

	err = row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, apperrors.NewApplicationError(mappings.RoleGetQueryError, err)
	}

	return &domain.Role{
		ID:   id,
		Name: domain.RoleName(name),
	}, nil
}

func (r *repository) executeGetQuery(ctx context.Context, filters GetFilterOptions) (*sql.Row, error) {
	query := `SELECT id, name FROM roles`

	query += ` WHERE 1=1 `

	args := []any{}

	if filters.ID != "" {
		query += ` AND id = $1`
		args = append(args, filters.ID)
	}

	if filters.Name != "" {
		query += ` AND name = $1`
		args = append(args, filters.Name)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
