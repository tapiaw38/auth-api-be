package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, user domain.User) (string, error) {
	row, err := r.executeCreateQuery(ctx, user)
	if err != nil {
		return "", err
	}

	var id string
	err = row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) executeCreateQuery(ctx context.Context, user domain.User) (*sql.Row, error) {
	query := `INSERT INTO users (
				id, first_name, last_name, username, email,
				password, phone_number, picture, address,
				is_active, verified_email,
				verified_email_token, verified_email_token_expiry,
				auth_method, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			RETURNING id`

	var phoneNumber, picture, address *string

	args := []any{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
		phoneNumber,
		picture,
		address,
		user.IsActive,
		user.VerifiedEmail,
		user.VerifiedEmailToken,
		user.VerifiedEmailTokenExpiry,
		user.AuthMethod,
		time.Now().UTC(),
		time.Now().UTC(),
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
