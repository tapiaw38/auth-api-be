package domain

import (
	"time"
)

const (
	SsoTypeGoogle SsoType = "google"
)

const (
	AuthMethodPassword AuthMethod = "password"
	AuthMethodGoogle   AuthMethod = "google"
	AuthMethodHybrid   AuthMethod = "hybrid"
)

type (
	SsoType    string
	AuthMethod string

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
		AuthMethod               string
		Roles                    []Role
		CreatedAt                time.Time
		UpdatedAt                time.Time
	}

	UserRole struct {
		UserID string
		RoleID string
	}
)
