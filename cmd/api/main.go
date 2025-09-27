package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources"
	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/auth-api-be/internal/adapters/workers"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
	"github.com/tapiaw38/auth-api-be/internal/platform/database"
	"github.com/tapiaw38/auth-api-be/internal/usecases"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

func main() {
	scope := config.GetScope()

	log.Printf("scope identifier: %s", scope)

	if err := initConfig(); err != nil {
		panic(err)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	configService := config.GetConfigService()

	db, err := database.GetSQLClientInstance()
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	err = database.Makemigration()
	if err != nil {
		return err
	}

	mq, err := queue.NewRabbitMQ(&configService)
	if err != nil {
		return err
	}

	if configService.ServerConfig.GinMode == config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowOrigins = []string{"*"}
	ginConfig.AllowCredentials = true
	ginConfig.AllowMethods = []string{"*"}
	ginConfig.AllowHeaders = []string{"*"}
	ginConfig.ExposeHeaders = []string{"*"}
	app.Use(cors.New(ginConfig))

	if err := bootstrap(app, db, mq, &configService); err != nil {
		return err
	}

	return app.Run(":" + configService.ServerConfig.Port)
}

func bootstrap(
	app *gin.Engine,
	db *sql.DB,
	mq *queue.RabbitMQ,
	configService *config.ConfigurationService,
) error {
	datasources := datasources.CreateDatasources(db)
	integrations := integrations.CreateIntegration(configService)

	contextFactory := appcontext.NewFactory(datasources, integrations, mq, configService)
	useCases := usecases.CreateUsecases(contextFactory)

	if !configService.InitConfig.EnsureDefaultRoles {
		log.Println("Skipping default roles initialization")
	}

	if err := ensureDefaultRoles(context.Background(), useCases.Role.EnsureUsecase); err != nil {
		log.Printf("Failed to ensure default roles: %v", err)
		return err
	}

	web.RegisterApplicationRoutes(app, useCases)

	if err := workers.RegisterWorkers(context.Background(), mq, contextFactory); err != nil {
		log.Fatalf("Failed to register workers: %v", err)
		return err
	}

	return nil
}

func ensureDefaultRoles(ctx context.Context, ensureUsecase role.EnsureUseCase) error {
	log.Println("Ensuring default roles exist...")
	err := ensureUsecase.Execute(ctx)
	if err != nil {
		return err
	}

	log.Println("Default roles ensured successfully")
	return nil
}
