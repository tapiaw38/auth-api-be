package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) Get(ctx context.Context, filters GetFilterOptions) (*domain.User, apperrors.ApplicationError) {
	row, err := r.executeGetQuery(ctx, filters)
	if err != nil {
		return nil, err
	}

	var (
		id, firstName, lastName, username, email, password, verifiedEmailToken, authMethod string
	)
	var phoneNumber, picture, address, passwordResetToken *string
	var isActive, verifiedEmail bool
	var createdAt, updatedAt, verifiedEmailTokenExpiry time.Time
	var passwordResetTokenExpiry *time.Time
	var tokenVersion uint
	var rolesJSON json.RawMessage

	err2 := row.Scan(
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
		&authMethod,
		&createdAt,
		&updatedAt,
		&rolesJSON,
	)
	if err2 != nil {
		if err2 == sql.ErrNoRows {
			return nil, nil
		}

		return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, err2)
	}

	roles, err3 := unmarshalRoles(rolesJSON)
	if err3 != nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, err3)
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
		authMethod,
		createdAt,
		updatedAt,
		roles,
	), nil
}

func (r *repository) executeGetQuery(ctx context.Context, filters GetFilterOptions) (*sql.Row, apperrors.ApplicationError) {
	query := `SELECT
				u.id, u.first_name, u.last_name, u.username,
				u.email, u.password, u.phone_number, u.picture, u.address,
				u.is_active, u.verified_email, u.verified_email_token,
				u.verified_email_token_expiry, u.password_reset_token,
				u.password_reset_token_expiry, u.token_version, u.auth_method,
				u.created_at, u.updated_at,
				COALESCE(
					json_agg(
						json_build_object(
							'id', r.id,
							'name', r.name
						)
					) FILTER (WHERE r.id IS NOT NULL), '[]'
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

	if filters.VerifiedEmailToken != "" {
		query += ` AND u.verified_email_token = $1`
		args = append(args, filters.VerifiedEmailToken)
	}

	if filters.PasswordResetToken != "" {
		query += ` AND u.password_reset_token = $1`
		args = append(args, filters.PasswordResetToken)
	}

	query += ` GROUP BY
		u.id, u.first_name, u.last_name,
		u.username, u.email, u.password,
		u.phone_number, u.picture, u.address,
		u.is_active, u.verified_email,
		u.verified_email_token, u.verified_email_token_expiry,
		u.password_reset_token, u.password_reset_token_expiry,
		u.token_version, u.auth_method,
		u.created_at, u.updated_at`

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, row.Err())
	}

	return row, nil
}
