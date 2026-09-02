package user

import "github.com/tapiaw38/auth-api-be/internal/domain"

type (
	UserOutputData struct {
		ID            string           `json:"id"`
		Username      string           `json:"username"`
		FirstName     string           `json:"first_name"`
		LastName      string           `json:"last_name"`
		Email         string           `json:"email"`
		PhoneNumber   *string          `json:"phone_number"`
		Picture       *string          `json:"picture"`
		Address       *string          `json:"address"`
		IsActive      bool             `json:"is_active"`
		VerifiedEmail bool             `json:"verified_email"`
		TokenVersion  uint             `json:"token_version"`
		AuthMethod    string           `json:"auth_method"`
		Roles         []RoleOutputData `json:"roles"`
	}

	RoleOutputData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// SummaryOutputData.ID is the username, not the UUID primary key — every
	// other service stores that username as its own "user_id" foreign
	// reference (it's what the JWT's user_id claim carries), so batch lookups
	// keyed by that same value need it echoed back, not the UUID.
	SummaryOutputData struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	ResetPasswordOutputData struct {
		Email   string `json:"email"`
		Message string `json:"message"`
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
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Picture:       user.Picture,
		Address:       user.Address,
		IsActive:      user.IsActive,
		VerifiedEmail: user.VerifiedEmail,
		TokenVersion:  user.TokenVersion,
		AuthMethod:    user.AuthMethod,
		Roles:         roles,
	}
}

func toSummaryOutputData(user *domain.User) SummaryOutputData {
	return SummaryOutputData{
		ID:        user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
