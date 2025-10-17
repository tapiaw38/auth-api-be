package auth

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/tapiaw38/auth-api-be/internal/platform/utils"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) ([]byte, error) {
	hashCost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ComparePassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	if !hasSpecial {
		return errors.New("password must contain at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)")
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	if len(email) > 254 {
		return errors.New("email must not exceed 254 characters")
	}

	return nil
}

func GenerateUsername(firstName, lastName string) string {
	normalizedFirst := strings.ToLower(strings.TrimSpace(firstName))
	normalizedLast := strings.ToLower(strings.TrimSpace(lastName))

	reg := regexp.MustCompile("[^a-z0-9]+")
	normalizedFirst = reg.ReplaceAllString(normalizedFirst, "")
	normalizedLast = reg.ReplaceAllString(normalizedLast, "")

	randomSuffix := utils.RandomString(10)
	username := normalizedFirst + "." + normalizedLast + "." + strings.ToLower(randomSuffix)

	return username
}
