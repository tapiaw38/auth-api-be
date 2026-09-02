package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func (r *repository) List(ctx context.Context, filters ListFilterOptions) ([]*domain.User, apperrors.ApplicationError) {
	rows, err := r.executeListQuery(ctx, filters)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var (
			id, firstName, lastName, username, email, password, authMethod string
		)
		var phoneNumber, picture, address, passwordResetToken *string
		var isActive, verifiedEmail bool
		var createdAt, updatedAt time.Time
		var verifiedEmailToken sql.NullString
		var verifiedEmailTokenExpiry sql.NullTime
		var passwordResetTokenExpiry sql.NullTime
		var tokenVersion uint
		var rolesJSON json.RawMessage

		err2 := rows.Scan(
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
			return nil, apperrors.NewApplicationError(mappings.UserListQueryError, err2)
		}

		roles, err3 := unmarshalRoles(rolesJSON)
		if err3 != nil {
			return nil, apperrors.NewApplicationError(mappings.UserListQueryError, err3)
		}

		var verifiedEmailTokenValue string
		if verifiedEmailToken.Valid {
			verifiedEmailTokenValue = verifiedEmailToken.String
		}

		var verifiedEmailTokenExpiryValue time.Time
		if verifiedEmailTokenExpiry.Valid {
			verifiedEmailTokenExpiryValue = verifiedEmailTokenExpiry.Time
		}

		var passwordResetTokenExpiryValue *time.Time
		if passwordResetTokenExpiry.Valid {
			passwordResetTokenExpiryValue = &passwordResetTokenExpiry.Time
		}

		users = append(users, unmarshalUser(
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
			verifiedEmailTokenValue,
			verifiedEmailTokenExpiryValue,
			passwordResetToken,
			passwordResetTokenExpiryValue,
			tokenVersion,
			authMethod,
			createdAt,
			updatedAt,
			roles,
		))
	}

	return users, nil
}

func (r *repository) executeListQuery(ctx context.Context, filters ListFilterOptions) (*sql.Rows, apperrors.ApplicationError) {
	query := `SELECT
                u.id, u.first_name, u.last_name, u.username,
                u.email, u.password, u.phone_number, u.picture, u.address,
                u.is_active, u.verified_email, u.verified_email_token,
                u.verified_email_token_expiry, u.password_reset_token,
                u.password_reset_token_expiry, u.token_version, u.auth_method,
                u.created_at, u.updated_at,
                COALESCE(
                    jsonb_agg(
                        jsonb_build_object(
                            'id', r.id,
                            'name', r.name
                        )
                    ) FILTER (WHERE r.id IS NOT NULL), '[]'
                ) AS roles
            FROM users u
            LEFT JOIN user_roles ur ON ur.user_id = u.id
            LEFT JOIN roles r ON r.id = ur.role_id
            WHERE 1=1`

	argIndex := 1
	var args []any

	if len(filters.IDs) > 0 {
		query += ` AND u.username = ANY($` + fmt.Sprintf("%d", argIndex) + `)`
		args = append(args, pq.Array(filters.IDs))
		argIndex++
	}

	if filters.IsActive != nil {
		query += ` AND u.is_active = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.IsActive)
		argIndex++
	}

	if filters.VerifiedEmail != nil {
		query += ` AND u.verified_email = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.VerifiedEmail)
		argIndex++
	}

	if filters.RoleID != "" {
		query += ` AND ur.role_id = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, filters.RoleID)
		argIndex++
	}

	if filters.RoleName != "" {
		query += ` AND r.name = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, filters.RoleName)
		argIndex++
	}

	if !filters.CreatedAt.IsZero() {
		query += ` AND u.created_at = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, filters.CreatedAt)
		argIndex++
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

	if filters.Limit > 0 {
		query += ` LIMIT $` + fmt.Sprintf("%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += ` OFFSET $` + fmt.Sprintf("%d", argIndex)
		args = append(args, filters.Offset)
		argIndex++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserListQueryError, err)
	}

	return rows, nil
}
