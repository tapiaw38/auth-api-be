# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Development

- `make run` - Run the API server locally (port 8082)
- `make run-dev` - Run in development mode with Air live reloader
- `make install-deps` - Install required Go tools (Delve debugger, Air live reloader)
- `make init-docker` - Start Docker services (PostgreSQL and RabbitMQ)

### Testing & Quality

- `make test` - Run all tests with coverage
- `make cover` - Display test coverage in function mode
- `make cover-html` - Generate HTML coverage report
- `make gen-mocks` - Generate mocks using go generate

### Database Management

- `make migrate-up` - Apply all database migrations
- `make migrate-down` - Revert the last migration
- `make migrate-create` - Create a new migration (prompts for name)
- Database runs on port 54321 (PostgreSQL 14)

### Build & Deployment

- `go build -o ./build/app.sh ./cmd/api/` - Build the application
- `docker-compose up -d` - Start all services via Docker

## Architecture Overview

This is a **Clean Architecture** Go REST API for authentication with the following structure:

### Core Layers

- **Domain** (`internal/domain/`) - Core business entities (User, Role)
- **Use Cases** (`internal/usecases/`) - Business logic layer with user and role operations
- **Adapters** (`internal/adapters/`) - External interfaces and implementations

### Adapter Structure

- **Web Layer** (`internal/adapters/web/`)

  - `handlers/` - HTTP request handlers organized by domain (user/, role/)
  - `middlewares/` - Authentication and other HTTP middlewares
  - `integrations/` - External service integrations (SSO, notifications)
  - `routes.go` - Main routing configuration

- **Data Layer** (`internal/adapters/datasources/`)

  - `repositories/` - Database access layer organized by domain
  - Uses PostgreSQL with direct SQL queries

- **Infrastructure**
  - `queue/` - RabbitMQ message queue integration
  - `workers/` - Background job processors

### Key Components

- **App Context Factory** (`internal/platform/appcontext/`) - Dependency injection container
- **Configuration** (`internal/platform/config/`) - Environment-based config management
- **Database Platform** (`internal/platform/database/`) - Database connection and migration management
- **Error Handling** (`internal/platform/errors/`) - Structured error management system

### Error Handling System

The application uses a structured error handling system based on best practices from `fury_taxes-workflow-be`:

#### Structure

```
internal/platform/errors/
├── application_error.go    # ApplicationError interface & implementation
├── logging.go              # Error logging functionality
└── mappings/
    ├── error_details.go   # Base ErrorDetails struct
    ├── auth.go            # Authentication errors
    ├── user.go            # User domain errors
    ├── role.go            # Role domain errors
    └── common.go          # Common/generic errors
```

#### Key Features

- **Unique error codes**: Each error has a unique identifier (e.g., `"user:login:invalid-credentials"`)
- **Appropriate HTTP status codes**: 400, 401, 403, 404, 409, 500, etc.
- **Centralized error messages**: All error messages defined in mapping files
- **Automatic logging**: Errors log with structured context
- **JSON serialization**: Consistent error response format
- **Original error tracking**: Preserves technical error details for debugging

#### Usage in Use Cases

```go
import (
    apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
    "github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

// Update interface signature
type LoginUsecase interface {
    Execute(context.Context, LoginInput) (*LoginOutput, apperrors.ApplicationError)
}

// Create application errors
func (u *loginUsecase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, apperrors.ApplicationError) {
    user, err := app.Repositories.User.Get(ctx, filters)
    if err != nil {
        return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, err)
    }

    if user == nil {
        return nil, apperrors.NewApplicationError(
            mappings.UserLoginUserNotFoundError,
            errors.New("user not found"),
        )
    }

    return &LoginOutput{Data: data}, nil
}
```

#### Usage in Handlers

```go
import (
    apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
    "github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

func NewLoginHandler(usecase user.LoginUsecase) gin.HandlerFunc {
    return func(c *gin.Context) {
        var login user.LoginInput
        if err := c.ShouldBindJSON(&login); err != nil {
            appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
            appErr.Log(c)  // Log with context
            c.JSON(appErr.StatusCode(), appErr)  // Return appropriate status
            return
        }

        output, appErr := usecase.Execute(c, login)
        if appErr != nil {
            appErr.Log(c)
            c.JSON(appErr.StatusCode(), appErr)
            return
        }

        c.JSON(http.StatusOK, output)
    }
}
```

#### Error Response Format

```json
{
  "code": "user:login:invalid-credentials",
  "message": "invalid email or password"
}
```

#### Adding New Errors

1. Define the error in the appropriate mapping file:

```go
// In mappings/user.go
var (
    UserProfileUpdateFailedError = ErrorDetails{
        "user:profile:update-failed",
        http.StatusInternalServerError,
        "failed to update user profile",
    }
)
```

2. Use it in your code:

```go
if err != nil {
    return nil, apperrors.NewApplicationError(mappings.UserProfileUpdateFailedError, err)
}
```

#### Error Code Naming Convention

Format: `{domain}:{action}:{error-type}`

Examples:
- `user:login:invalid-credentials`
- `role:create:name-required`
- `user:reset-password:invalid-token`
- `auth:token-version-mismatch`

#### Best Practices

1. **Always log errors** before returning to client: `appErr.Log(c)`
2. **Use appropriate HTTP status codes**: 400 (validation), 401 (auth), 404 (not found), 500 (server)
3. **User-friendly messages**: `Message` field for users, `OriginalMessage` for debugging
4. **Never expose sensitive data**: Keep error messages generic for security
5. **Add context when needed**: Use `AddExtraFields()` for additional logging context

### Authentication & Authorization

- JWT-based authentication with token versioning
- Role-based access control (RBAC)
- Google OAuth2 SSO integration
- Email verification and password reset workflows
- Authorization middleware protects routes requiring authentication

### Message Queue Architecture

- RabbitMQ integration for async processing
- Worker pattern for background jobs
- Email notifications and other async tasks

## Key Features

- User registration/login with email verification
- JWT token management with version control for secure logout
- Role management and assignment
- Password reset via email
- Google OAuth2 integration
- Background job processing for emails
- PostgreSQL database with migrations
- Docker containerization

## Database Schema

- **users** - Main user entity with auth fields and profile data
- **roles** - Role definitions
- **user_roles** - Many-to-many relationship between users and roles

## Environment Configuration

Configure via `.env` file with:

- Database connection (PostgreSQL on port 54321)
- JWT secrets
- Google OAuth2 credentials
- AWS S3 for file storage
- Email server (SMTP)
- RabbitMQ connection

## Development Notes

- Entry point: `cmd/api/main.go`
- Uses Gin web framework
- Air for live reloading in development
- Database migrations in `migrations/` directory
- All handlers follow Clean Architecture patterns with dependency injection

---

## Development Guidelines

### Creating New Features (Complete Flow)

When implementing a new feature, follow this order:

1. **Domain Model** (if needed) - Define entities in `internal/domain/`
2. **Repository Interface & Implementation** - Data access in `internal/adapters/datasources/repositories/`
3. **Use Case** - Business logic in `internal/usecases/`
4. **Handler** - HTTP layer in `internal/adapters/web/handlers/`
5. **Routes** - Register in `internal/adapters/web/routes.go`
6. **Tests** - Unit tests for each layer
7. **Database Migration** - Schema changes via `make migrate-create`

### 1. Domain Layer (`internal/domain/`)

Define core business entities and constants.

**Example: `internal/domain/role.go`**

```go
package domain

const (
    RoleSuperAdmin RoleName = "superadmin"
    RoleAdmin      RoleName = "admin"
    RoleUser       RoleName = "user"
)

type (
    RoleName string

    Role struct {
        ID   string
        Name RoleName
    }
)
```

### 2. Repository Layer (`internal/adapters/datasources/repositories/`)

#### Repository Structure

Each domain gets its own package with:

- `repository.go` - Interface definition
- `create.go`, `get.go`, `update.go`, `delete.go`, `list.go` - Method implementations
- `*_test.go` - Tests for each method

**Example: `internal/adapters/datasources/repositories/role/repository.go`**

```go
package role

import (
    "context"
    "database/sql"
    "github.com/tapiaw38/auth-api-be/internal/domain"
)

type (
    Repository interface {
        Create(context.Context, domain.Role) (string, error)
        Get(context.Context, GetFilterOptions) (*domain.Role, error)
        Update(context.Context, string, *domain.Role) (string, error)
        Delete(context.Context, string) error
        List(context.Context, ListFilterOptions) ([]domain.Role, error)
    }

    repository struct {
        db *sql.DB
    }

    GetFilterOptions struct {
        ID   string
        Name string
    }

    ListFilterOptions struct {
        Name string
    }
)

func NewRepository(db *sql.DB) Repository {
    return &repository{db: db}
}
```

**Example Implementation: `create.go`**

```go
package role

import (
    "context"
    "database/sql"
    "github.com/tapiaw38/auth-api-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, role domain.Role) (string, error) {
    row, err := r.executeCreateQuery(ctx, role)
    if err != nil {
        return "", err
    }

    var id string
    if err := row.Scan(&id); err != nil {
        return "", err
    }

    return id, nil
}

func (r *repository) executeCreateQuery(ctx context.Context, role domain.Role) (*sql.Row, error) {
    query := `INSERT INTO roles (id, name) VALUES ($1, $2) RETURNING id;`

    args := []any{role.ID, role.Name}

    row := r.db.QueryRowContext(ctx, query, args...)
    if row.Err() != nil {
        return nil, row.Err()
    }

    return row, nil
}
```

**Key Patterns:**

- Use `QueryRowContext` for single row operations
- Use `QueryContext` for multiple rows with `Scan` in a loop
- Always use parameterized queries (`$1`, `$2`, etc.)
- Separate query execution from scanning for better testability
- Return domain entities, not database-specific types

### 3. Use Case Layer (`internal/usecases/`)

#### Use Case Structure

Each domain gets its own package with:

- Individual files per operation (e.g., `create.go`, `get.go`)
- `output_types.go` - Shared output data structures and mapper functions
- `*_test.go` - Tests for each use case

**Example: `internal/usecases/role/create.go`**

```go
package role

import (
    "context"
    "errors"
    roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
    "github.com/tapiaw38/auth-api-be/internal/domain"
    "github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
    CreateUsecase interface {
        Execute(context.Context, CreateInput) (*CreateOutput, error)
    }

    createUsecase struct {
        contextFactory appcontext.Factory
    }

    CreateInput struct {
        Name string `json:"name"`
    }

    CreateOutput struct {
        Data RoleOutputData `json:"data"`
    }
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
    return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
    app := u.contextFactory()

    // Validation
    if input.Name == "" {
        return nil, errors.New("role name is required")
    }

    // Business logic
    roleName := domain.RoleName(input.Name)
    switch roleName {
    case domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleUser:
    default:
        return nil, errors.New("invalid role name")
    }

    role := domain.Role{Name: roleName}

    // Repository operations
    id, err := app.Repositories.Role.Create(ctx, role)
    if err != nil {
        return nil, err
    }

    createdRole, err := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
    if err != nil {
        return nil, err
    }

    return &CreateOutput{
        Data: toRoleOutputData(*createdRole),
    }, nil
}
```

**Example: `output_types.go`**

```go
package role

import "github.com/tapiaw38/auth-api-be/internal/domain"

type (
    RoleOutputData struct {
        ID   string `json:"id"`
        Name string `json:"name"`
    }
)

func toRoleOutputData(role domain.Role) RoleOutputData {
    return RoleOutputData{
        ID:   role.ID,
        Name: string(role.Name),
    }
}
```

**Key Patterns:**

- Each use case has an interface and implementation
- Input/Output structs define the contract
- Use `contextFactory` for dependency injection
- Centralize output data structures in `output_types.go`
- Include mapper functions (`toXOutputData`) to convert domain to output
- Keep business logic isolated from HTTP and database concerns

### 4. Handler Layer (`internal/adapters/web/handlers/`)

#### Handler Structure

Each domain gets its own package with individual handler files.

**Example: `internal/adapters/web/handlers/role/create.go`**

```go
package role

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

type CreateInput struct {
    Name string `json:"name" binding:"required"`
}

func NewCreateHandler(usecase role.CreateUsecase) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input CreateInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Invalid request body",
                "error":   err.Error(),
            })
            return
        }

        output, err := usecase.Execute(c, role.CreateInput{
            Name: input.Name,
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": err.Error(),
            })
            return
        }

        c.JSON(http.StatusCreated, output)
    }
}
```

**Key Patterns:**

- Define handler-specific input structs with validation tags (`binding:"required"`)
- Use factory functions that return `gin.HandlerFunc`
- Map handler input to use case input explicitly
- Use appropriate HTTP status codes (200, 201, 400, 404, 500)
- Return consistent error response format
- Keep handlers thin - delegate to use cases

### 5. Input/Output Handling

#### Input Structure (3 Levels)

1. **Handler Input** - HTTP validation layer

   ```go
   type CreateInput struct {
       Name string `json:"name" binding:"required"`
   }
   ```

2. **Use Case Input** - Business logic contract

   ```go
   type CreateInput struct {
       Name string `json:"name"`
   }
   ```

3. **Repository Parameters** - Data access options

   ```go
   type GetFilterOptions struct {
       ID   string
       Name string
   }
   ```

#### Output Structure

1. **Repository** - Returns domain entities

   ```go
   func Create(ctx context.Context, role domain.Role) (string, error)
   func Get(ctx context.Context, opts GetFilterOptions) (*domain.Role, error)
   ```

2. **Use Case** - Returns structured output

   ```go
   type CreateOutput struct {
       Data RoleOutputData `json:"data"`
   }
   ```

3. **Handler** - Returns JSON via Gin

   ```go
   c.JSON(http.StatusCreated, output)
   ```

**Best Practices:**

- Handler inputs use validation tags
- Use case inputs contain business data without HTTP concerns
- Always map between layers explicitly
- Use shared output types in `output_types.go`
- Repository returns domain entities, use case transforms to output data

### 6. Registering Routes (`internal/adapters/web/routes.go`)

```go
func RegisterApplicationRoutes(app *gin.Engine, useCases *usecases.Usecases) {
    routeGroup := app.Group("/")

    // Public routes
    routeGroup.POST("auth/register", user.NewRegisterHandler(useCases.User.RegisterUsecase))
    routeGroup.POST("auth/login", user.NewLoginHandler(useCases.User.LoginUsecase))

    // Protected routes (after middleware)
    routeGroup.Use(middlewares.AuthorizationMiddleware(useCases.User.GetTokenVersionUsecase))
    routeGroup.GET("user/me", user.NewMeHandler(useCases.User.GetUsecase))
    routeGroup.POST("role/create", role.NewCreateHandler(useCases.Role.CreateUsecase))
}
```

**Key Patterns:**

- Group routes by prefix
- Apply middleware with `.Use()` before protected routes
- Pass use cases via dependency injection
- Follow RESTful conventions where appropriate

### 7. Testing Structure

#### Repository Tests (`*_test.go`)

Use `sqlmock` for database testing.

**Example: `internal/adapters/datasources/repositories/role/create_test.go`**

```go
package role

import (
    "context"
    "database/sql"
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
                rows := sqlmock.NewRows([]string{"id"}).AddRow("role-123")
                mock.ExpectQuery(`INSERT INTO roles`).
                    WithArgs("role-123", domain.RoleAdmin).
                    WillReturnRows(rows)
            },
            expectedID:    "role-123",
            expectedError: false,
        },
        {
            name: "database connection error",
            role: domain.Role{ID: "role-123", Name: domain.RoleAdmin},
            mockSetup: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery(`INSERT INTO roles`).
                    WithArgs("role-123", domain.RoleAdmin).
                    WillReturnError(sql.ErrConnDone)
            },
            expectedID:    "",
            expectedError: true,
            errorMsg:      "connection is already closed",
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

            // Act
            id, err := repo.Create(context.Background(), tt.role)

            // Assert
            if tt.expectedError {
                assert.Error(t, err)
                if tt.errorMsg != "" {
                    assert.Contains(t, err.Error(), tt.errorMsg)
                }
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedID, id)
            }

            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}
```

#### Use Case Tests (`*_test.go`)

Use `gomock` for mocking repository dependencies.

**Example: `internal/usecases/role/create_test.go`**

```go
package role_test

import (
    "context"
    "errors"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
    roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
    mock_role "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role/mocks"
    "github.com/tapiaw38/auth-api-be/internal/domain"
    "github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
    usecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
    "go.uber.org/mock/gomock"
)

func TestCreateUsecase_Execute(t *testing.T) {
    type fields struct {
        repository *mock_role.MockRepository
    }

    tests := map[string]struct {
        input       usecase.CreateInput
        prepare     func(f *fields)
        expected    *usecase.CreateOutput
        expectedErr error
    }{
        "when creating role successfully": {
            input: usecase.CreateInput{Name: "admin"},
            prepare: func(f *fields) {
                f.repository.EXPECT().
                    Create(gomock.Any(), domain.Role{Name: domain.RoleAdmin}).
                    Return("role-123", nil)
                f.repository.EXPECT().
                    Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
                    Return(&domain.Role{ID: "role-123", Name: domain.RoleAdmin}, nil)
            },
            expected: &usecase.CreateOutput{
                Data: usecase.RoleOutputData{
                    ID:   "role-123",
                    Name: "admin",
                },
            },
        },
        "when creating role with empty name": {
            input:       usecase.CreateInput{Name: ""},
            prepare:     func(f *fields) {},
            expected:    nil,
            expectedErr: errors.New("role name is required"),
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            f := fields{repository: mock_role.NewMockRepository(ctrl)}
            if tc.prepare != nil {
                tc.prepare(&f)
            }

            contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
                return &appcontext.Context{
                    Repositories: &repositories.Repositories{Role: f.repository},
                }
            }

            uc := usecase.NewCreateUsecase(contextFactory)
            actual, actualErr := uc.Execute(context.Background(), tc.input)

            assert.Equal(t, tc.expected, actual)
            assert.Equal(t, tc.expectedErr, actualErr)
        })
    }
}
```

#### Handler Tests (`*_test.go`)

Use `gomock` for mocking use case dependencies.

**Example: `internal/adapters/web/handlers/role/create_test.go`**

```go
package role_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    role_handler "github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/role"
    mock_role "github.com/tapiaw38/auth-api-be/internal/usecases/role/mocks"
    roleUsecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
    "go.uber.org/mock/gomock"
)

func TestCreateHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := map[string]struct {
        body         string
        setupUsecase func(*mock_role.MockCreateUsecase)
        expectedCode int
        expectedBody string
    }{
        "when create usecase executes successfully": {
            body: `{"name":"admin"}`,
            setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {
                mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.CreateInput{Name: "admin"}).
                    Return(&roleUsecase.CreateOutput{
                        Data: roleUsecase.RoleOutputData{ID: "role-123", Name: "admin"},
                    }, nil)
            },
            expectedCode: 201,
            expectedBody: `{"data":{"id":"role-123","name":"admin"}}`,
        },
        "when name is missing": {
            body:         `{}`,
            setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {},
            expectedCode: 400,
            expectedBody: `{"message":"Invalid request body"}`,
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUsecase := mock_role.NewMockCreateUsecase(ctrl)
            if tc.setupUsecase != nil {
                tc.setupUsecase(mockUsecase)
            }

            handler := role_handler.NewCreateHandler(mockUsecase)

            req := httptest.NewRequest(http.MethodPost, "/role", bytes.NewBufferString(tc.body))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = req

            handler(c)

            assert.Equal(t, tc.expectedCode, w.Code)
            if tc.expectedBody != "" {
                assert.Contains(t, w.Body.String(), tc.expectedBody)
            }
        })
    }
}
```

### Test Patterns Summary

**Repository Layer:**

- Use `sqlmock` for database mocking
- Test query execution and error handling
- Use table-driven tests with struct slices
- Test both success and error paths
- Verify SQL expectations with `mock.ExpectationsWereMet()`

**Use Case Layer:**

- Use `gomock` for repository mocking
- Create mock context factory for dependency injection
- Test business logic validation
- Use table-driven tests with named map entries
- Test edge cases and error propagation

**Handler Layer:**

- Use `gomock` for use case mocking
- Use `gin.TestMode` and `httptest` for HTTP testing
- Test input validation and binding
- Test HTTP status codes and response bodies
- Use table-driven tests with named map entries

### Generating Mocks

After defining interfaces, generate mocks:

```bash
make gen-mocks
```

This uses `go generate` with `mockgen` to create mock implementations for testing.

### Database Migrations

Create a new migration:

```bash
make migrate-create
# Enter migration name when prompted
```

This creates two files in `migrations/`:

- `NNNNNN_name.up.sql` - Schema changes
- `NNNNNN_name.down.sql` - Rollback changes

Apply migrations:

```bash
make migrate-up
```

Rollback last migration:

```bash
make migrate-down
```
