package role

import (
	"context"
	"database/sql"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	result, err := r.executeDeleteQuery(ctx, id)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (r *repository) executeDeleteQuery(ctx context.Context, id string) (sql.Result, error) {
	query := `DELETE FROM roles WHERE id = $1`

	args := []any{
		id,
	}

	return r.db.ExecContext(ctx, query, args...)
}
