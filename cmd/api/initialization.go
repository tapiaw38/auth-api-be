package main

import (
	"context"
	"log"

	"github.com/tapiaw38/auth-api-be/internal/platform/config"
	"github.com/tapiaw38/auth-api-be/internal/usecases"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

// initializeDefaults runs all startup initialization tasks based on config flags
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
	// Example:
	// if configService.InitConfig.EnsureDefaultUsers {
	//     if err := ensureDefaultUsers(ctx, useCases.User.EnsureUsecase); err != nil {
	//         return err
	//     }
	// }

	log.Println("Application defaults initialized successfully")
	return nil
}

// ensureDefaultRoles creates default system roles (superadmin, admin, user) if they don't exist
func ensureDefaultRoles(ctx context.Context, ensureUsecase role.EnsureUseCase) error {
	log.Println("Ensuring default roles exist...")
	err := ensureUsecase.Execute(ctx)
	if err != nil {
		return err
	}

	log.Println("Default roles ensured successfully")
	return nil
}
