package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func AuthorizationMiddleware(usecase user.GetTokenVersionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := tokenParts[1]

		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		ctx := c.Request.Context()
		tokenVersion, err := usecase.Execute(ctx, claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to get token version"})
			return
		}

		if claims.TokenVersion != tokenVersion {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token version mismatch"})
			return
		}

		ctx = context.WithValue(ctx, "userID", claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
