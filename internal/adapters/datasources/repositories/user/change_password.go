
package user

import (
	"context"
	"log"
)

func (r *repository) ChangePassword(ctx context.Context, id string, password string) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET password = $1 WHERE id = $2;")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, password, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
