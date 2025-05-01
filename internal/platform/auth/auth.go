package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tapiaw38/auth-api-be/internal/domain"
)

type (
	CustomClaims struct {
		UserID       string `json:"user_id"`
		TokenVersion uint   `json:"token_version"`
		jwt.StandardClaims
	}
)

func GenerateToken(user *domain.User, expiration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID:       user.Username,
		TokenVersion: user.TokenVersion,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string, secret []byte, getUserTokenVersion func(userID string) (uint, error)) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	dbVersion, err := getUserTokenVersion(claims.UserID)
	if err != nil {
		return nil, err
	}

	if dbVersion != claims.TokenVersion {
		return nil, errors.New("invalid token version")
	}

	return claims, nil
}
