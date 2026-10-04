package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewGoogleMobilePollHandler(usecase user.GoogleMobileLoginUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, found := usecase.Poll(c.Param("state"))
		if !found {
			c.JSON(http.StatusOK, gin.H{"status": user.MobileSessionPending})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
