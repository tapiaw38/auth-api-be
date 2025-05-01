package auth

import (
	"errors"

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
