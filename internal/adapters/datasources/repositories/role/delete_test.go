package role

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError bool
		errorMsg      string
	}{
		{
			name: "successful role deletion",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name: "successful role deletion - multiple rows affected",
			id:   "role-456",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-456").
					WillReturnResult(sqlmock.NewResult(0, 2))
			},
			expectedError: false,
		},
		{
			name: "successful role deletion - no rows affected",
			id:   "non-existent-role",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("non-existent-role").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "database connection error",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "query execution error",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnError(driver.ErrBadConn)
			},
			expectedError: true,
			errorMsg:      "expected a connection to be available",
		},
		{
			name: "transaction error",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnError(sql.ErrTxDone)
			},
			expectedError: true,
			errorMsg:      "sql: transaction has already been committed or rolled back",
		},
		{
			name: "rows affected error",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				result := sqlmock.NewErrorResult(sql.ErrConnDone)
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnResult(result)
			},
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "empty ID parameter",
			id:   "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "whitespace ID parameter",
			id:   "   ",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("   ").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "special characters in ID",
			id:   "role-123!@#$%^&*()",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123!@#$%^&*()").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name: "very long ID parameter",
			id:   "role-" + string(make([]byte, 1000)),
			mockSetup: func(mock sqlmock.Sqlmock) {
				longID := "role-" + string(make([]byte, 1000))
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs(longID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "SQL injection attempt in ID",
			id:   "role-123'; DROP TABLE roles; --",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123'; DROP TABLE roles; --").
					WillReturnResult(sqlmock.NewResult(0, 0))
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
			err = repo.Delete(ctx, tt.id)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_executeDeleteQuery(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError bool
		errorMsg      string
	}{
		{
			name: "successful query execution",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name: "successful query execution - no rows affected",
			id:   "non-existent-role",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("non-existent-role").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "query execution with database error",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorMsg:      "sql: connection is already closed",
		},
		{
			name: "query execution with bad connection",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnError(driver.ErrBadConn)
			},
			expectedError: true,
			errorMsg:      "expected a connection to be available",
		},
		{
			name: "query execution with transaction error",
			id:   "role-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("role-123").
					WillReturnError(sql.ErrTxDone)
			},
			expectedError: true,
			errorMsg:      "sql: transaction has already been committed or rolled back",
		},
		{
			name: "empty ID parameter",
			id:   "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "null-like ID parameter",
			id:   "null",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("null").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false,
		},
		{
			name: "UUID-like ID parameter",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnResult(sqlmock.NewResult(0, 1))
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

			repo := &repository{db: db}
			tt.mockSetup(mock)

			ctx := context.Background()

			// Act
			result, err := repo.executeDeleteQuery(ctx, tt.id)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// TestRepository_DeleteEdgeCases tests additional edge cases and scenarios
func TestRepository_DeleteEdgeCases(t *testing.T) {
	t.Run("context timeout", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewRepository(db)
		ctx := context.Background()

		mock.ExpectExec(`DELETE FROM roles WHERE id = \$1`).
			WithArgs("role-123").
			WillReturnError(context.DeadlineExceeded)

		// Act
		err = repo.Delete(ctx, "role-123")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}