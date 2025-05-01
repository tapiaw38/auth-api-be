package domain

import (
	"time"
)

const (
	SsoTypeGoogle SsoType = "google"
)

type (
	SsoType string

	User struct {
		ID                       string
		FirstName                string
		LastName                 string
		Username                 string
		Email                    string
		Password                 string
		PhoneNumber              *string
		Picture                  *string
		Address                  *string
		IsActive                 bool
		VerifiedEmail            bool
		VerifiedEmailToken       string
		VerifiedEmailTokenExpiry time.Time
		PasswordResetToken       *string
		PasswordResetTokenExpiry *time.Time
		TokenVersion             uint
		Roles                    []Role
		CreatedAt                time.Time
		UpdatedAt                time.Time
	}

	UserRole struct {
		UserID string
		RoleID string
	}
)
