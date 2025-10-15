package user

import (
	"encoding/json"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/domain"
)

func unmarshalUser(
	id string,
	firstName string,
	lastName string,
	username string,
	email string,
	password string,
	phoneNumber *string,
	picture *string,
	address *string,
	isActive bool,
	verifiedEmail bool,
	verifiedEmailToken string,
	verifiedEmailTokenExpiry time.Time,
	passwordResetToken *string,
	passwordResetTokenExpiry *time.Time,
	tokenVersion uint,
	authMethod string,
	createdAt time.Time,
	updatedAt time.Time,
	roles []domain.Role,
) *domain.User {
	return &domain.User{
		ID:                       id,
		FirstName:                firstName,
		LastName:                 lastName,
		Username:                 username,
		Email:                    email,
		Password:                 password,
		PhoneNumber:              phoneNumber,
		Picture:                  picture,
		Address:                  address,
		IsActive:                 isActive,
		VerifiedEmail:            verifiedEmail,
		VerifiedEmailToken:       verifiedEmailToken,
		VerifiedEmailTokenExpiry: verifiedEmailTokenExpiry,
		PasswordResetToken:       passwordResetToken,
		PasswordResetTokenExpiry: passwordResetTokenExpiry,
		TokenVersion:             tokenVersion,
		AuthMethod:               authMethod,
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
		Roles:                    roles,
	}
}

func unmarshalRoles(raw json.RawMessage) ([]domain.Role, error) {
	var rolesJSON []RoleJSON
	if err := json.Unmarshal(raw, &rolesJSON); err != nil {
		return nil, err
	}

	roles := make([]domain.Role, 0, len(rolesJSON))
	for _, roleJSON := range rolesJSON {
		roles = append(roles, domain.Role{
			ID:   roleJSON.ID,
			Name: domain.RoleName(roleJSON.Name),
		})
	}

	return roles, nil
}
