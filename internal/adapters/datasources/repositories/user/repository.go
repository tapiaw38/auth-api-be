package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
)

type (
	Repository interface {
		Create(context.Context, domain.User) (string, apperrors.ApplicationError)
		Get(context.Context, GetFilterOptions) (*domain.User, apperrors.ApplicationError)
		Update(context.Context, string, *domain.User) (string, apperrors.ApplicationError)
		Delete(context.Context, string) apperrors.ApplicationError
		List(context.Context, ListFilterOptions) ([]*domain.User, apperrors.ApplicationError)
		ChangePassword(ctx context.Context, id string, password string) apperrors.ApplicationError
		InvalidatePasswordResetToken(ctx context.Context, id string) apperrors.ApplicationError
	}

	repository struct {
		db *sql.DB
	}

	GetFilterOptions struct {
		ID                 string
		Username           string
		Email              string
		VerifiedEmailToken string
		PasswordResetToken string
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
