package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type googleMobileCallbackInput struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func NewGoogleMobileCallbackHandler(usecase user.GoogleMobileLoginUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input googleMobileCallbackInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(http.StatusOK, gin.H{"message": "Ya podés volver a la app"})
			return
		}

		if appErr := usecase.Complete(c, input.Code, input.State); appErr != nil {
			appErr.Log(c)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ya podés volver a la app"})
	}
}
