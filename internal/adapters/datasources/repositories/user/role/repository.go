package user_role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

type (
	Repository interface {
		Create(context.Context, domain.UserRole) (*domain.UserRole, error)
		Delete(context.Context, string, string) (*domain.UserRole, error)
	}

	repository struct {
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
