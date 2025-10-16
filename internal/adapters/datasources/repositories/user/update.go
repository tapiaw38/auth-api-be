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
			auth_method = COALESCE($14, auth_method),
			updated_at = $15
		WHERE id = $16
		RETURNING id;`

	var (
		verifiedEmailToken                                *string
		verifiedEmailTokenExpiry                          *time.Time
		picture, phoneNumber, address, passwordResetToken *string
		passwordResetTokenExpiry                          *time.Time
	)

	if user.VerifiedEmailToken != "" {
		verifiedEmailToken = &user.VerifiedEmailToken
	}

	if !user.VerifiedEmailTokenExpiry.IsZero() {
		verifiedEmailTokenExpiry = &user.VerifiedEmailTokenExpiry
	}

	if user.Picture != nil {
		picture = user.Picture
	}

	if user.PhoneNumber != nil {
		phoneNumber = user.PhoneNumber
	}

	if user.Address != nil {
		address = user.Address
	}

	if user.PasswordResetToken != nil {
		passwordResetToken = user.PasswordResetToken
	}

	if user.PasswordResetTokenExpiry != nil && !user.PasswordResetTokenExpiry.IsZero() {
		passwordResetTokenExpiry = user.PasswordResetTokenExpiry
	}

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
		user.AuthMethod,
		user.UpdatedAt,
		id,
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
