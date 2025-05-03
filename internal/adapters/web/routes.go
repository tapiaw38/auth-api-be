package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/role"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/user"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/auth-api-be/internal/usecases"
)

func RegisterApplicationRoutes(app *gin.Engine, useCases *usecases.Usecases) {
	routeGroup := app.Group("/")

	routeGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routeGroup.POST("role/ensure", role.NewEnsureHandler(useCases.Role.EnsureUsecase))
	routeGroup.POST("auth/register", user.NewRegisterHandler(useCases.User.RegisterUsecase))
	routeGroup.POST("auth/login", user.NewLoginHandler(useCases.User.LoginUsecase))

	routeGroup.Use(middlewares.AuthorizationMiddleware(useCases.User.GetTokenVersionUsecase))

	routeGroup.GET("role/list", role.NewListHandler(useCases.Role.ListUsecase))
}
