package role

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func TestRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		role          domain.Role
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedID    string
		expectedError bool
		errorMsg      string
	}{
		{
			name: "successful role creation",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("role-123")
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnRows(rows)
			},
			expectedID:    "role-123",
			expectedError: false,
		},
		{
			name: "successful role creation with user role",
			role: domain.Role{
				ID:   "role-456",
				Name: domain.RoleUser,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("role-456")
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-456", domain.RoleUser).
					WillReturnRows(rows)
			},
			expectedID:    "role-456",
			expectedError: false,
		},
		{
			name: "successful role creation with superadmin role",
			role: domain.Role{
				ID:   "role-789",
				Name: domain.RoleSuperAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("role-789")
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-789", domain.RoleSuperAdmin).
					WillReturnRows(rows)
			},
			expectedID:    "role-789",
			expectedError: false,
		},
		{
			name: "database connection error on query execution",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:    "",
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "query execution error",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnError(sql.ErrTxDone)
			},
			expectedID:    "",
			expectedError: true,
			errorMsg:      "sql: transaction has already been committed or rolled back",
		},
		{
			name: "scan error - no rows returned",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnRows(rows)
			},
			expectedID:    "",
			expectedError: true,
			errorMsg:      "sql: no rows in result set",
		},
		{
			name: "scan error - incorrect column type",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Note: sqlmock converts types automatically, so we expect success
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(123) // This gets converted to string "123"
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnRows(rows)
			},
			expectedID:    "123",
			expectedError: false,
		},
		{
			name: "role with empty ID",
			role: domain.Role{
				ID:   "",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("")
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("", domain.RoleAdmin).
					WillReturnRows(rows)
			},
			expectedID:    "",
			expectedError: false,
		},
		{
			name: "role with empty name",
			role: domain.Role{
				ID:   "role-123",
				Name: "",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("role-123")
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", "").
					WillReturnRows(rows)
			},
			expectedID:    "role-123",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewRepository(db)
			tt.mockSetup(mock)

			ctx := context.Background()

			// Act
			id, err := repo.Create(ctx, tt.role)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_executeCreateQuery(t *testing.T) {
	tests := []struct {
		name          string
		role          domain.Role
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError bool
		errorMsg      string
	}{
		{
			name: "successful query execution",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("role-123")
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "query execution with database error",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "query execution with constraint violation",
			role: domain.Role{
				ID:   "role-123",
				Name: domain.RoleAdmin,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO roles \(id, name\) VALUES \(\$1, \$2\) RETURNING id;`).
					WithArgs("role-123", domain.RoleAdmin).
					WillReturnError(driver.ErrBadConn)
			},
			expectedError: true,
			errorMsg:      "expected a connection to be available", // Actual error message from sqlmock
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := &repository{db: db}
			tt.mockSetup(mock)

			ctx := context.Background()

			// Act
			row, err := repo.executeCreateQuery(ctx, tt.role)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, row)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, row)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}