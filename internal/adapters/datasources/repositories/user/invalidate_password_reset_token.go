
package user

import (
	"context"
	"log"
)

func (r *repository) InvalidatePasswordResetToken(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET password_reset_token = NULL, password_reset_token_expiry = NULL WHERE id = $1;")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
