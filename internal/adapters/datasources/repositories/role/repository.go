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
	return &repository{
		db: db,
	}
}
