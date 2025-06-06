package usecases

import (
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type Usecases struct {
	User User
	Role Role
}

type User struct {
	RegisterUsecase             user.RegisterUsecase
	LoginUsecase                user.LoginUsecase
	GetUsecase                  user.GetUsecase
	UpdateUsecase               user.UpdateUsecase
	DeleteUsecase               user.DeleteUsecase
	ListUsecase                 user.ListUsecase
	GetTokenVersionUsecase      user.GetTokenVersionUsecase
	VerifyEmailUsecase          user.VerifyEmailUsecase
	ResetPasswordUsecase        user.ResetPasswordUsecase
	RequestResetPasswordUsecase user.RequestResetPasswordUsecase
}

type Role struct {
	EnsureUsecase role.EnsureUseCase
	ListUsecase   role.ListUsecase
}

func CreateUsecases(contextFactory appcontext.Factory) *Usecases {
	return &Usecases{
		User: User{
			RegisterUsecase:             user.NewCreateUsecase(contextFactory),
			LoginUsecase:                user.NewLoginUsecase(contextFactory),
			GetUsecase:                  user.NewGetUsecase(contextFactory),
			UpdateUsecase:               user.NewUpdateUsecase(contextFactory),
			DeleteUsecase:               user.NewDeleteUsecase(contextFactory),
			ListUsecase:                 user.NewListUsecase(contextFactory),
			GetTokenVersionUsecase:      user.NewGetTokenVersionUsecase(contextFactory),
			VerifyEmailUsecase:          user.NewVerifyEmailUsecase(contextFactory),
			ResetPasswordUsecase:        user.NewResetPasswordUsecase(contextFactory),
			RequestResetPasswordUsecase: user.NewRequestResetPasswordUsecase(contextFactory),
		},
		Role: Role{
			EnsureUsecase: role.NewEnsureUseCase(contextFactory),
			ListUsecase:   role.NewListUsecase(contextFactory),
		},
	}
}
