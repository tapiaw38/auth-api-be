package user

import "github.com/tapiaw38/auth-api-be/internal/domain"

type (
	RegisterInput struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
)

func (i *RegisterInput) toDomain() *domain.User {
	return &domain.User{
		FirstName: i.FirstName,
		LastName:  i.LastName,
		Username:  i.Username,
		Email:     i.Email,
		Password:  i.Password,
	}
}
