package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type (
	RoleClaim struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	CustomClaims struct {
		UserID       string      `json:"user_id"`
		TokenVersion uint        `json:"token_version"`
		Roles        []RoleClaim `json:"roles"`
		jwt.StandardClaims
	}
)

func (r RoleClaim) GetID() string {
	return r.ID
}

func (r RoleClaim) GetName() string {
	return r.Name
}

func GenerateToken(user *domain.User, expiration time.Duration) (string, error) {
	roles := make([]RoleClaim, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, RoleClaim{
			ID:   role.ID,
			Name: string(role.Name),
		})
	}

	claims := CustomClaims{
		UserID:       user.Username,
		TokenVersion: user.TokenVersion,
		Roles:        roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetConfigService().ServerConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		secret := config.GetConfigService().ServerConfig.JWTSecret
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
