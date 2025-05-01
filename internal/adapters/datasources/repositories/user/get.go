package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, filters GetFilterOptions) (*domain.User, error) {
	row, err := r.executeGetQuery(ctx, filters)
	if err != nil {
		return nil, err
	}

	var (
		id, firstName, lastName, username, email, password, verifiedEmailToken string
	)
	var phoneNumber, picture, address, passwordResetToken *string
	var isActive, verifiedEmail bool
	var createdAt, updatedAt, verifiedEmailTokenExpiry time.Time
	var passwordResetTokenExpiry *time.Time
	var tokenVersion uint
	var rolesJSON json.RawMessage

	err = row.Scan(
		&id,
		&firstName,
		&lastName,
		&username,
		&email,
		&password,
		&phoneNumber,
		&picture,
		&address,
		&isActive,
		&verifiedEmail,
		&verifiedEmailToken,
		&verifiedEmailTokenExpiry,
		&passwordResetToken,
		&passwordResetTokenExpiry,
		&tokenVersion,
		&createdAt,
		&updatedAt,
		&rolesJSON,
	)
	if err != nil {
		return nil, err
	}

	roles, err := unmarshalRoles(rolesJSON)
	if err != nil {
		return nil, err
	}

	return unmarshalUser(
		id,
		firstName,
		lastName,
		username,
		email,
		password,
		phoneNumber,
		picture,
		address,
		isActive,
		verifiedEmail,
		verifiedEmailToken,
		verifiedEmailTokenExpiry,
		passwordResetToken,
		passwordResetTokenExpiry,
		tokenVersion,
		createdAt,
		updatedAt,
		roles,
	), nil

}

func (r *repository) executeGetQuery(ctx context.Context, filters GetFilterOptions) (*sql.Row, error) {
	query := `SELECT 
                u.id, u.first_name, u.last_name, u.username, 
                u.email, u.password, u.phone_number, u.picture, u.address, 
                u.is_active, u.verified_email, u.verified_email_token,
                u.verified_email_token_expiry, u.password_reset_token,
                u.password_reset_token_expiry, u.token_version,
                u.created_at, u.updated_at,
                COALESCE(
                    jsonb_agg(
                        jsonb_build_object(
                            'id', r.id,
                            'name', r.name
                        )
                    ), '[]'
                ) AS roles
            FROM users u
            LEFT JOIN user_roles ur ON ur.user_id = u.id
            LEFT JOIN roles r ON r.id = ur.role_id
            `

	query += ` WHERE 1=1 `

	args := []any{}

	if filters.ID != "" {
		query += ` AND u.id = $1`
		args = append(args, filters.ID)
	}

	if filters.Username != "" {
		query += ` AND u.username = $1`
		args = append(args, filters.Username)
	}

	if filters.Email != "" {
		query += ` AND u.email = $1`
		args = append(args, filters.Email)
	}

	query += ` GROUP BY 
		u.id, u.first_name, u.last_name, 
		u.username, u.email, u.password, 
		u.phone_number, u.picture, u.address, 
		u.is_active, u.verified_email, 
		u.verified_email_token, u.verified_email_token_expiry, 
		u.password_reset_token, u.password_reset_token_expiry, 
		u.created_at, u.updated_at`

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}
