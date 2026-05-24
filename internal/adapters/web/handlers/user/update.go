package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/domain"
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
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		payload := &domain.User{
			FirstName:   input.FirstName,
			LastName:    input.LastName,
			Email:       input.Email,
			Picture:     input.Picture,
			PhoneNumber: input.PhoneNumber,
			Address:     input.Address,
			UpdatedAt:   time.Now(),
		}
		if input.IsActive != nil {
			payload.IsActive = *input.IsActive
		}
		if input.VerifiedEmail != nil {
			payload.VerifiedEmail = *input.VerifiedEmail
		}

		output, appErr := usecase.Execute(c, id, payload)
		if appErr != nil {
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
