package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewVerifyEmailHandler(usecase user.VerifyEmailUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		VerifiedEmailToken := c.Query("token")
		if VerifiedEmailToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "token is required",
			})
			return
		}

		redirectURL, err := usecase.Execute(c, VerifiedEmailToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}
