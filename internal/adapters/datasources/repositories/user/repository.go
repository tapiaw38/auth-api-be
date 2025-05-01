package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

type (
	Repository interface {
		Create(context.Context, domain.User) (string, error)
		Get(context.Context, GetFilterOptions) (*domain.User, error)
		Update(context.Context, string, *domain.User) (string, error)
		Delete(context.Context, string) error
		List(context.Context, ListFilterOptions) ([]*domain.User, error)
	}

	repository struct {
		db *sql.DB
	}

	GetFilterOptions struct {
		ID       string
		Username string
		Email    string
	}

	ListFilterOptions struct {
		IsActive      *bool
		VerifiedEmail *bool
		RoleID        string
		CreatedAt     time.Time
		Limit         int
		Offset        int
	}

	RoleJSON struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
