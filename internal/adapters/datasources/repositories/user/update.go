package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, user *domain.User) (string, error) {
	row, err := r.executeUpdateQuery(ctx, id, user)
	if err != nil {
		return "", err
	}

	var updatedID string
	if err := row.Scan(&updatedID); err != nil {
		return "", err
	}

	return updatedID, nil
}

func (r *repository) executeUpdateQuery(ctx context.Context, id string, user *domain.User) (*sql.Row, error) {
	query := `UPDATE users
		SET 
			first_name = COALESCE($1, first_name), 
			last_name = COALESCE($2, last_name), 
			email = COALESCE($3, email),
			password = COALESCE($4, password), 
			picture = COALESCE($5, picture), 
			phone_number = COALESCE($6, phone_number), 
			address = COALESCE($7, address), 
			is_active = COALESCE($8, is_active), 
			verified_email = COALESCE($9, verified_email), 
			verified_email_token = COALESCE($10, verified_email_token), 
			verified_email_token_expiry = COALESCE($11, verified_email_token_expiry),
			password_reset_token = COALESCE($12, password_reset_token), 
			password_reset_token_expiry = COALESCE($13, password_reset_token_expiry),
			updated_at = $14
		WHERE id = $15
		RETURNING id;`

	var (
		picture, phoneNumber, address, verifiedEmailToken, passwordResetToken *string
		verifiedEmailTokenExpiry, passwordResetTokenExpiry                    *time.Time
	)

	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		picture,
		phoneNumber,
		address,
		user.IsActive,
		user.VerifiedEmail,
		verifiedEmailToken,
		verifiedEmailTokenExpiry,
		passwordResetToken,
		passwordResetTokenExpiry,
		user.UpdatedAt,
		id,
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
