package role

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
)

type (
	Repository interface {
		Create(context.Context, domain.Role) (string, apperrors.ApplicationError)
		Get(context.Context, GetFilterOptions) (*domain.Role, apperrors.ApplicationError)
		Update(context.Context, string, *domain.Role) (string, apperrors.ApplicationError)
		Delete(context.Context, string) apperrors.ApplicationError
		List(context.Context, ListFilterOptions) ([]domain.Role, apperrors.ApplicationError)
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
