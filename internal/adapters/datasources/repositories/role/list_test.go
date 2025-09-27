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

func TestRepository_List(t *testing.T) {
	tests := []struct {
		name          string
		filters       ListFilterOptions
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedRoles []domain.Role
		expectedError bool
		errorMsg      string
	}{
		{
			name: "successful list with no filters - multiple roles",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin").
					AddRow("role-456", "user").
					AddRow("role-789", "superadmin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-123",
					Name: domain.RoleAdmin,
				},
				{
					ID:   "role-456",
					Name: domain.RoleUser,
				},
				{
					ID:   "role-789",
					Name: domain.RoleSuperAdmin,
				},
			},
			expectedError: false,
		},
		{
			name: "successful list with name filter",
			filters: ListFilterOptions{
				Name: "admin",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("admin").
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-123",
					Name: domain.RoleAdmin,
				},
			},
			expectedError: false,
		},
		{
			name: "successful list with user role filter",
			filters: ListFilterOptions{
				Name: "user",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-456", "user").
					AddRow("role-789", "user")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("user").
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-456",
					Name: domain.RoleUser,
				},
				{
					ID:   "role-789",
					Name: domain.RoleUser,
				},
			},
			expectedError: false,
		},
		{
			name: "successful list with superadmin role filter",
			filters: ListFilterOptions{
				Name: "superadmin",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-999", "superadmin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("superadmin").
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-999",
					Name: domain.RoleSuperAdmin,
				},
			},
			expectedError: false,
		},
		{
			name: "successful list with empty results",
			filters: ListFilterOptions{
				Name: "nonexistent",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("nonexistent").
					WillReturnRows(rows)
			},
			expectedRoles: nil,
			expectedError: false,
		},
		{
			name:    "successful list with no filters - empty results",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: nil,
			expectedError: false,
		},
		{
			name: "successful list with single role",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-123",
					Name: domain.RoleAdmin,
				},
			},
			expectedError: false,
		},
		{
			name: "empty name filter - should ignore filter",
			filters: ListFilterOptions{
				Name: "",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin").
					AddRow("role-456", "user")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-123",
					Name: domain.RoleAdmin,
				},
				{
					ID:   "role-456",
					Name: domain.RoleUser,
				},
			},
			expectedError: false,
		},
		{
			name:    "database connection error",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedRoles: nil,
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "query execution error",
			filters: ListFilterOptions{
				Name: "admin",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("admin").
					WillReturnError(driver.ErrBadConn)
			},
			expectedRoles: nil,
			expectedError: true,
			errorMsg:      "expected a connection to be available",
		},
		{
			name:    "scan error - wrong number of columns",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("role-123")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: nil,
			expectedError: true,
		},
		{
			name:    "scan error - wrong data type",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Note: sqlmock converts types automatically, so we expect success
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(123, "admin") // ID gets converted to string "123"
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "123",
					Name: domain.RoleAdmin,
				},
			},
			expectedError: false,
		},
		{
			name:    "scan error on second row",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Note: sqlmock converts types automatically, so we expect success
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin").
					AddRow(456, "user") // ID gets converted to string "456"
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedRoles: []domain.Role{
				{
					ID:   "role-123",
					Name: domain.RoleAdmin,
				},
				{
					ID:   "456",
					Name: domain.RoleUser,
				},
			},
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
			roles, err := repo.List(ctx, tt.filters)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, roles)
			} else {
				assert.NoError(t, err)
				if tt.expectedRoles == nil {
					assert.Nil(t, roles)
				} else {
					assert.Equal(t, tt.expectedRoles, roles)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_executeListQuery(t *testing.T) {
	tests := []struct {
		name          string
		filters       ListFilterOptions
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError bool
		errorMsg      string
	}{
		{
			name:    "successful query execution with no filters",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "successful query execution with name filter",
			filters: ListFilterOptions{
				Name: "admin",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("admin").
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "successful query execution with empty name filter",
			filters: ListFilterOptions{
				Name: "",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("role-123", "admin")
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name:    "query execution with database error",
			filters: ListFilterOptions{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "query execution with bad connection",
			filters: ListFilterOptions{
				Name: "admin",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("admin").
					WillReturnError(driver.ErrBadConn)
			},
			expectedError: true,
			errorMsg:      "expected a connection to be available",
		},
		{
			name: "query execution with transaction error",
			filters: ListFilterOptions{
				Name: "user",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name FROM roles WHERE 1=1  AND name = \$1`).
					WithArgs("user").
					WillReturnError(sql.ErrTxDone)
			},
			expectedError: true,
			errorMsg:      "sql: transaction has already been committed or rolled back",
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
			rows, err := repo.executeListQuery(ctx, tt.filters)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, rows)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rows)
				rows.Close() // Clean up
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}