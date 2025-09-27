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
