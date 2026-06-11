package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
	defer func() {
		if err := mq.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

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

	if err := bootstrap(ctx, app, db, mq, &configService); err != nil {
		return err
	}

	return app.Run(":" + configService.ServerConfig.Port)
}

func bootstrap(
	ctx context.Context,
	app *gin.Engine,
	db *sql.DB,
	mq *queue.RabbitMQ,
	configService *config.ConfigurationService,
) error {
	datasources := datasources.CreateDatasources(db)
	integrations := integrations.CreateIntegration(configService)

	contextFactory := appcontext.NewFactory(datasources, integrations, mq, configService)
	useCases := usecases.CreateUsecases(contextFactory)

	if err := initializeDefaults(context.Background(), configService, useCases); err != nil {
		log.Printf("Failed to initialize defaults: %v", err)
		return err
	}

	web.RegisterApplicationRoutes(app, useCases)

	if err := workers.RegisterWorkers(ctx, mq, contextFactory); err != nil {
		log.Printf("Failed to register workers: %v", err)
		return err
	}

	return nil
}

func initializeDefaults(ctx context.Context, configService *config.ConfigurationService, useCases *usecases.Usecases) error {
	log.Println("Initializing application defaults...")

	if configService.InitConfig.EnsureDefaultRoles {
		if err := ensureDefaultRoles(ctx, useCases.Role.EnsureUsecase); err != nil {
			return err
		}
	} else {
		log.Println("Skipping default roles initialization")
	}

	// Add more initialization tasks here in the future
	// if configService.InitConfig.SomeOtherInit {
	//     if err := someOtherInit(ctx, ...); err != nil {
	//         return err
	//     }
	// }

	log.Println("Application defaults initialized successfully")
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
