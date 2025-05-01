package user

import "github.com/tapiaw38/auth-api-be/internal/domain"

type (
	UserOutputData struct {
		ID            string           `json:"id"`
		FirstName     string           `json:"first_name"`
		LastName      string           `json:"last_name"`
		Email         string           `json:"email"`
		PhoneNumber   *string          `json:"phone_number"`
		Picture       *string          `json:"picture"`
		Address       *string          `json:"address"`
		IsActive      bool             `json:"is_active"`
		VerifiedEmail bool             `json:"verified_email"`
		TokenVersion  uint             `json:"token_version"`
		Roles         []RoleOutputData `json:"roles"`
	}

	RoleOutputData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func toUserOutputData(user *domain.User) UserOutputData {
	var roles []RoleOutputData
	for _, role := range user.Roles {
		roles = append(roles, RoleOutputData{
			ID:   role.ID,
			Name: string(role.Name),
		})
	}

	return UserOutputData{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Picture:       user.Picture,
		Address:       user.Address,
		IsActive:      user.IsActive,
		VerifiedEmail: user.VerifiedEmail,
		TokenVersion:  user.TokenVersion,
		Roles:         roles,
	}
}
