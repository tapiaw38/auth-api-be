package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

// RequireRoles corta el pedido si quien lo hace no tiene ninguno de los roles
// esperados. Se usa después de AuthorizationMiddleware, que es el que deja los
// roles del token en el contexto.
func RequireRoles(expected ...domain.RoleName) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, expected...) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

// HasRole responde si el token trae alguno de los roles esperados.
func HasRole(c *gin.Context, expected ...domain.RoleName) bool {
	claims, _ := c.Request.Context().Value("userRoles").([]auth.RoleClaim)
	for _, claim := range claims {
		for _, want := range expected {
			if claim.Name == string(want) {
				return true
			}
		}
	}

	return false
}
