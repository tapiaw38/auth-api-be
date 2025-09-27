# Auth API Backend - QWEN.md

## Project Overview

This is a **Clean Architecture** Go REST API for authentication services, providing user management, authentication, authorization, and integration with external services. The project follows modern Go practices with a clean architecture pattern, separating concerns into distinct layers.

**Module**: `github.com/tapiaw38/auth-api-be`

**Architecture Pattern**: Clean Architecture with:
- **Domain** (`internal/domain/`) - Core business entities (User, Role)
- **Use Cases** (`internal/usecases/`) - Business logic layer with user and role operations
- **Adapters** (`internal/adapters/`) - External interfaces and implementations
  - **Web Layer** (`internal/adapters/web/`) - HTTP handlers, middlewares, routes
  - **Data Layer** (`internal/adapters/datasources/`) - Database repositories
  - **Infrastructure** (`internal/adapters/queue/`, `internal/adapters/workers/`) - RabbitMQ, background workers

## Core Dependencies

- **Gin Web Framework** (`github.com/gin-gonic/gin`) - HTTP web framework
- **Gin CORS** (`github.com/gin-contrib/cors`) - Cross-origin resource sharing middleware
- **JWT** (`github.com/golang-jwt/jwt`) - JSON Web Token authentication
- **PostgreSQL** (`github.com/lib/pq`, `github.com/jackc/pgx/v5`) - Database driver
- **Database Migrations** (`github.com/golang-migrate/migrate/v4`) - Migration management
- **RabbitMQ** (`github.com/streadway/amqp`) - Message queue integration
- **OAuth2** (`golang.org/x/oauth2`) - Google OAuth2 integration
- **UUID** (`github.com/google/uuid`) - UUID generation

## Building and Running

### Development Commands
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

## Key Features

- **User Management**: Registration, login with email verification
- **JWT Authentication**: Token-based authentication with version control for secure logout
- **Role-Based Access Control (RBAC)**: Role management and assignment
- **Password Reset**: Via email workflow
- **Google OAuth2 Integration**: Single Sign-On (SSO) capability
- **Background Jobs**: RabbitMQ integration for async processing (email notifications)
- **Database Migrations**: PostgreSQL with automated schema management
- **Docker Containerization**: For consistent deployment

## Infrastructure & Services

- **PostgreSQL 14**: Primary database with migrations
- **RabbitMQ**: Message queue for background processing
- **Google OAuth2**: External authentication provider
- **AWS S3**: File storage (configuration available)
- **SMTP Email**: For notifications and password reset emails

## Configuration

The application uses environment-based configuration managed through:
- `.env` file for local development
- Environment variables in production
- Configuration service in `internal/platform/config/`

Configuration includes:
- Database connection (PostgreSQL on port 54321)
- JWT secrets for token management
- Google OAuth2 credentials
- AWS S3 for file storage
- Email server (SMTP) settings
- RabbitMQ connection details

## Project Structure

```
├── cmd/
│   └── api/                 # Application entry point
├── internal/
│   ├── domain/              # Business entities (User, Role)
│   ├── usecases/            # Business logic layer
│   └── adapters/            # External interface implementations
│       ├── datasources/     # Database repositories
│       ├── web/             # HTTP handlers, middlewares, routes
│       ├── queue/           # RabbitMQ integration
│       └── workers/         # Background job processors
├── migrations/              # Database migration files
├── rabbit-mq/               # RabbitMQ configuration
├── templates/               # Email and other templates
├── .air.toml               # Air live reloader configuration
├── docker-compose.yml      # Docker services configuration
├── Dockerfile              # Docker build configuration
└── Makefile                # Build and development commands
```

## Authentication & Authorization

- JWT-based authentication with token versioning
- Role-based access control (RBAC) system
- Google OAuth2 SSO integration
- Email verification and password reset workflows
- Authorization middleware protecting routes requiring authentication

## Development Conventions

- Clean Architecture patterns with clear separation of concerns
- Dependency injection through the app context factory
- Comprehensive error handling throughout the application
- Standard Go project layout following Go conventions
- Gin framework for HTTP routing and middleware
- PostgreSQL database with direct SQL queries
- Automated testing with Go's built-in testing tools
- Docker and Docker Compose for containerization and local development

## Key Components

- **Entry Point**: `cmd/api/main.go`
- **App Context Factory**: `internal/platform/appcontext/` - Dependency injection container
- **Configuration**: `internal/platform/config/` - Environment-based config management
- **Database Platform**: `internal/platform/database/` - Database connection and migration management
- **Routes Configuration**: `internal/adapters/web/routes.go` - Main routing setup
- **Message Queue**: `internal/adapters/queue/` - RabbitMQ integration
- **Background Workers**: `internal/adapters/workers/` - Async job processing

## Environment Variables

The application expects various environment variables (typically loaded from .env):
- Database connection details
- JWT secrets
- OAuth2 credentials (Google)
- Email SMTP settings
- RabbitMQ connection string
- AWS S3 configuration
- Server port and mode

## Testing

- Standard Go testing framework
- Test coverage reporting with `make cover` and `make cover-html`
- Mock generation using `go generate`
- Unit tests for business logic and integration tests for API endpoints