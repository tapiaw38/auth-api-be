package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type updateInput struct {
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         string  `json:"email"`
	Picture       *string `json:"picture"`
	PhoneNumber   *string `json:"phone_number"`
	Address       *string `json:"address"`
	IsActive      *bool   `json:"is_active"`
	VerifiedEmail *bool   `json:"verified_email"`
}

func NewUpdateByIDHandler(usecase user.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		username, ok := c.Request.Context().Value("userID").(string)
		if !ok || username == "" {
			appErr := apperrors.NewApplicationError(mappings.AuthUnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		roles, _ := c.Request.Context().Value("userRoles").([]auth.RoleClaim)

		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c, user.UpdateInput{
			ID:            id,
			AuthUsername:  username,
			AuthRoles:     roles,
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			Email:         input.Email,
			Picture:       input.Picture,
			PhoneNumber:   input.PhoneNumber,
			Address:       input.Address,
			IsActive:      input.IsActive,
			VerifiedEmail: input.VerifiedEmail,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
