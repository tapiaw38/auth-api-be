package user_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func TestRepository_Get(t *testing.T) {
	validDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	phoneNumber := "+1234567890"
	picture := "https://example.com/avatar.jpg"
	address := "123 Main St, City"
	passwordResetToken := "reset-token-123"
	passwordResetTokenExpiry := validDate.Add(1 * time.Hour)

	columns := []string{
		"id",
		"first_name",
		"last_name",
		"username",
		"email",
		"password",
		"phone_number",
		"picture",
		"address",
		"is_active",
		"verified_email",
		"verified_email_token",
		"verified_email_token_expiry",
		"password_reset_token",
		"password_reset_token_expiry",
		"token_version",
		"created_at",
		"updated_at",
		"roles",
	}

	type fields struct {
		mock sqlmock.Sqlmock
	}

	tests := map[string]struct {
		filters   user.GetFilterOptions
		prepare   func(f *fields)
		expect    *domain.User
		expectErr error
	}{
		"when getting user by ID successfully": {
			filters: user.GetFilterOptions{
				ID: "user-123",
			},
			prepare: func(f *fields) {
				rolesJSON := buildRolesJSON(`[{"id":"role-1","name":"user"}]`)
				rows := sqlmock.NewRows(columns)
				rows.AddRow(
					"user-123",
					"John",
					"Doe",
					"johndoe",
					"john@example.com",
					"hashedpassword",
					phoneNumber,
					picture,
					address,
					true,
					true,
					"token123",
					validDate.Add(24*time.Hour),
					passwordResetToken,
					passwordResetTokenExpiry,
					1,
					validDate,
					validDate,
					rolesJSON,
				)
				f.mock.ExpectQuery("SELECT").WithArgs("user-123").WillReturnRows(rows)
			},
			expect: &domain.User{
				ID:                       "user-123",
				FirstName:                "John",
				LastName:                 "Doe",
				Username:                 "johndoe",
				Email:                    "john@example.com",
				Password:                 "hashedpassword",
				PhoneNumber:              &phoneNumber,
				Picture:                  &picture,
				Address:                  &address,
				IsActive:                 true,
				VerifiedEmail:            true,
				VerifiedEmailToken:       "token123",
				VerifiedEmailTokenExpiry: validDate.Add(24 * time.Hour),
				PasswordResetToken:       &passwordResetToken,
				PasswordResetTokenExpiry: &passwordResetTokenExpiry,
				TokenVersion:             1,
				Roles: []domain.Role{
					{ID: "role-1", Name: domain.RoleUser},
				},
				CreatedAt: validDate,
				UpdatedAt: validDate,
			},
		},
		"when getting user by username successfully": {
			filters: user.GetFilterOptions{
				Username: "johndoe",
			},
			prepare: func(f *fields) {
				rolesJSON := buildRolesJSON(`[{"id":"role-1","name":"admin"}]`)
				rows := sqlmock.NewRows(columns)
				rows.AddRow(
					"user-456",
					"Jane",
					"Smith",
					"johndoe",
					"jane@example.com",
					"hashedpassword",
					nil,
					nil,
					nil,
					true,
					false,
					"token456",
					validDate.Add(24*time.Hour),
					nil,
					nil,
					1,
					validDate,
					validDate,
					rolesJSON,
				)
				f.mock.ExpectQuery("SELECT").WithArgs("johndoe").WillReturnRows(rows)
			},
			expect: &domain.User{
				ID:                       "user-456",
				FirstName:                "Jane",
				LastName:                 "Smith",
				Username:                 "johndoe",
				Email:                    "jane@example.com",
				Password:                 "hashedpassword",
				PhoneNumber:              nil,
				Picture:                  nil,
				Address:                  nil,
				IsActive:                 true,
				VerifiedEmail:            false,
				VerifiedEmailToken:       "token456",
				VerifiedEmailTokenExpiry: validDate.Add(24 * time.Hour),
				PasswordResetToken:       nil,
				PasswordResetTokenExpiry: nil,
				TokenVersion:             1,
				Roles: []domain.Role{
					{ID: "role-1", Name: domain.RoleAdmin},
				},
				CreatedAt: validDate,
				UpdatedAt: validDate,
			},
		},
		"when getting user by email successfully": {
			filters: user.GetFilterOptions{
				Email: "john@example.com",
			},
			prepare: func(f *fields) {
				rolesJSON := buildRolesJSON(`[]`)
				rows := sqlmock.NewRows(columns)
				rows.AddRow(
					"user-789",
					"John",
					"Johnson",
					"johnjohnson",
					"john@example.com",
					"hashedpassword",
					nil,
					nil,
					nil,
					false,
					false,
					"token789",
					validDate.Add(24*time.Hour),
					nil,
					nil,
					1,
					validDate,
					validDate,
					rolesJSON,
				)
				f.mock.ExpectQuery("SELECT").WithArgs("john@example.com").WillReturnRows(rows)
			},
			expect: &domain.User{
				ID:                       "user-789",
				FirstName:                "John",
				LastName:                 "Johnson",
				Username:                 "johnjohnson",
				Email:                    "john@example.com",
				Password:                 "hashedpassword",
				PhoneNumber:              nil,
				Picture:                  nil,
				Address:                  nil,
				IsActive:                 false,
				VerifiedEmail:            false,
				VerifiedEmailToken:       "token789",
				VerifiedEmailTokenExpiry: validDate.Add(24 * time.Hour),
				PasswordResetToken:       nil,
				PasswordResetTokenExpiry: nil,
				TokenVersion:             1,
				Roles:                    []domain.Role{},
				CreatedAt:                validDate,
				UpdatedAt:                validDate,
			},
		},
		"when getting user by verified email token successfully": {
			filters: user.GetFilterOptions{
				VerifiedEmailToken: "token123",
			},
			prepare: func(f *fields) {
				rolesJSON := buildRolesJSON(`[{"id":"role-1","name":"user"}]`)
				rows := sqlmock.NewRows(columns)
				rows.AddRow(
					"user-123",
					"John",
					"Doe",
					"johndoe",
					"john@example.com",
					"hashedpassword",
					nil,
					nil,
					nil,
					true,
					false,
					"token123",
					validDate.Add(24*time.Hour),
					nil,
					nil,
					1,
					validDate,
					validDate,
					rolesJSON,
				)
				f.mock.ExpectQuery("SELECT").WithArgs("token123").WillReturnRows(rows)
			},
			expect: &domain.User{
				ID:                       "user-123",
				FirstName:                "John",
				LastName:                 "Doe",
				Username:                 "johndoe",
				Email:                    "john@example.com",
				Password:                 "hashedpassword",
				PhoneNumber:              nil,
				Picture:                  nil,
				Address:                  nil,
				IsActive:                 true,
				VerifiedEmail:            false,
				VerifiedEmailToken:       "token123",
				VerifiedEmailTokenExpiry: validDate.Add(24 * time.Hour),
				PasswordResetToken:       nil,
				PasswordResetTokenExpiry: nil,
				TokenVersion:             1,
				Roles: []domain.Role{
					{ID: "role-1", Name: domain.RoleUser},
				},
				CreatedAt: validDate,
				UpdatedAt: validDate,
			},
		},
		"when getting user by password reset token successfully": {
			filters: user.GetFilterOptions{
				PasswordResetToken: "reset-token-123",
			},
			prepare: func(f *fields) {
				rolesJSON := buildRolesJSON(`[]`)
				rows := sqlmock.NewRows(columns)
				rows.AddRow(
					"user-123",
					"John",
					"Doe",
					"johndoe",
					"john@example.com",
					"hashedpassword",
					nil,
					nil,
					nil,
					true,
					true,
					"token123",
					validDate.Add(24*time.Hour),
					"reset-token-123",
					&passwordResetTokenExpiry,
					1,
					validDate,
					validDate,
					rolesJSON,
				)
				f.mock.ExpectQuery("SELECT").WithArgs("reset-token-123").WillReturnRows(rows)
			},
			expect: &domain.User{
				ID:                       "user-123",
				FirstName:                "John",
				LastName:                 "Doe",
				Username:                 "johndoe",
				Email:                    "john@example.com",
				Password:                 "hashedpassword",
				PhoneNumber:              nil,
				Picture:                  nil,
				Address:                  nil,
				IsActive:                 true,
				VerifiedEmail:            true,
				VerifiedEmailToken:       "token123",
				VerifiedEmailTokenExpiry: validDate.Add(24 * time.Hour),
				PasswordResetToken:       &passwordResetToken,
				PasswordResetTokenExpiry: &passwordResetTokenExpiry,
				TokenVersion:             1,
				Roles:                    []domain.Role{},
				CreatedAt:                validDate,
				UpdatedAt:                validDate,
			},
		},
		"when no rows are returned": {
			filters: user.GetFilterOptions{
				ID: "non-existent-user",
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("SELECT").WithArgs("non-existent-user").
					WillReturnRows(sqlmock.NewRows(columns).
						RowError(0, sql.ErrNoRows))
			},
			expect:    nil,
			expectErr: nil,
		},
		"when a query error occurs": {
			filters: user.GetFilterOptions{
				ID: "user-123",
			},
			prepare: func(f *fields) {
				f.mock.ExpectQuery("SELECT").WithArgs("user-123").
					WillReturnError(errors.New("database connection error"))
			},
			expectErr: errors.New("database connection error"),
		},
		"when roles JSON parsing fails": {
			filters: user.GetFilterOptions{
				ID: "user-123",
			},
			prepare: func(f *fields) {
				invalidRolesJSON := buildRolesJSON(`[{"id":"role-1","name":invalid}]`)
				rows := sqlmock.NewRows(columns)
				rows.AddRow(
					"user-123",
					"John",
					"Doe",
					"johndoe",
					"john@example.com",
					"hashedpassword",
					nil,
					nil,
					nil,
					true,
					true,
					"token123",
					validDate.Add(24*time.Hour),
					nil,
					nil,
					1,
					validDate,
					validDate,
					invalidRolesJSON,
				)
				f.mock.ExpectQuery("SELECT").WithArgs("user-123").WillReturnRows(rows)
			},
			expectErr: errors.New("invalid character 'i' looking for beginning of value"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			if tt.prepare != nil {
				tt.prepare(&fields{mock: mock})
			}

			repository := user.NewRepository(db)
			result, err := repository.Get(context.Background(), tt.filters)

			assert.Equal(t, tt.expect, result)

			if tt.expectErr != nil {
				assert.Error(t, err)
				if tt.expectErr.Error() == "invalid character 'i' looking for beginning of value" {
					assert.Contains(t, err.Error(), "invalid character 'i' looking for beginning of value")
				} else {
					assert.Equal(t, tt.expectErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func buildRolesJSON(roles string) json.RawMessage {
	return json.RawMessage([]byte(roles))
}
