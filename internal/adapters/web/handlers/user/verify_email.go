package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewVerifyEmailHandler(usecase user.VerifyEmailUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		VerifiedEmailToken := c.Query("token")
		if VerifiedEmailToken == "" {
			appErr := apperrors.NewApplicationError(mappings.UserVerifyEmailTokenRequiredError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		redirectURL, appErr := usecase.Execute(c, VerifiedEmailToken)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}
