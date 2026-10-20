package main

import (
	"context"
	"log"

	"github.com/tapiaw38/auth-api-be/internal/platform/config"
	"github.com/tapiaw38/auth-api-be/internal/usecases"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

func initializeDefaults(ctx context.Context, configService *config.ConfigurationService, useCases *usecases.Usecases) error {
	log.Println("Initializing application defaults...")

	if configService.InitConfig.EnsureDefaultRoles {
		if err := ensureDefaultRoles(ctx, useCases.Role.EnsureUsecase); err != nil {
			return err
		}
	} else {
		log.Println("Skipping default roles initialization")
	}

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
