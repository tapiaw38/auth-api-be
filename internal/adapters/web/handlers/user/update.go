package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	platformweb "github.com/tapiaw38/auth-api-be/internal/platform/web"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type updateInput struct {
	FirstName     platformweb.RequestOptional[string] `json:"first_name"`
	LastName      platformweb.RequestOptional[string] `json:"last_name"`
	Email         platformweb.RequestOptional[string] `json:"email"`
	Picture       platformweb.RequestOptional[string] `json:"picture"`
	PhoneNumber   platformweb.RequestOptional[string] `json:"phone_number"`
	Address       platformweb.RequestOptional[string] `json:"address"`
	IsActive      platformweb.RequestOptional[bool]   `json:"is_active"`
	VerifiedEmail platformweb.RequestOptional[bool]   `json:"verified_email"`
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
			ID:             id,
			AuthUsername:   username,
			CanManageUsers: domain.CanManageUsers(domain.RolesFromClaims(roles)),
			Patch: domain.UserPatch{
				FirstName:     toDomainOptional(input.FirstName),
				LastName:      toDomainOptional(input.LastName),
				Email:         toDomainOptional(input.Email),
				Picture:       toDomainOptional(input.Picture),
				PhoneNumber:   toDomainOptional(input.PhoneNumber),
				Address:       toDomainOptional(input.Address),
				IsActive:      toDomainOptional(input.IsActive),
				VerifiedEmail: toDomainOptional(input.VerifiedEmail),
			},
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func toDomainOptional[T any](field platformweb.RequestOptional[T]) domain.Optional[T] {
	return domain.Optional[T]{
		Set:   field.Set,
		Value: field.Value,
	}
}
