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
	BatchUsecase                user.BatchUsecase
	GetTokenVersionUsecase      user.GetTokenVersionUsecase
	VerifyEmailUsecase          user.VerifyEmailUsecase
	ResetPasswordUsecase        user.ResetPasswordUsecase
	RequestResetPasswordUsecase user.RequestResetPasswordUsecase
	ChangePasswordUsecase       user.ChangePasswordUsecase
	SetPasswordUsecase          user.SetPasswordUsecase
	UpdateRolesUsecase          user.UpdateRolesUsecase
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
			BatchUsecase:                user.NewBatchUsecase(contextFactory),
			GetTokenVersionUsecase:      user.NewGetTokenVersionUsecase(contextFactory),
			VerifyEmailUsecase:          user.NewVerifyEmailUsecase(contextFactory),
			ResetPasswordUsecase:        user.NewResetPasswordUsecase(contextFactory),
			RequestResetPasswordUsecase: user.NewRequestResetPasswordUsecase(contextFactory),
			ChangePasswordUsecase:       user.NewChangePasswordUsecase(contextFactory),
			SetPasswordUsecase:          user.NewSetPasswordUsecase(contextFactory),
			UpdateRolesUsecase:          user.NewUpdateRolesUsecase(contextFactory),
		},
		Role: Role{
			EnsureUsecase: role.NewEnsureUseCase(contextFactory),
			ListUsecase:   role.NewListUsecase(contextFactory),
		},
	}
}
