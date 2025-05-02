package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

func initConfig() error {
	configService, err := readConfig()
	if err != nil {
		return err
	}
	config.InitConfigService(configService)
	return nil
}

func readConfig() (*config.ConfigurationService, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error load env file")
	}

	configService := &config.ConfigurationService{
		AppName: getEnv("APP_NAME", ""),
		ServerConfig: config.ServerConfig{
			GinMode:   config.GinModeServer(getEnv("GIN_MODE", "release")),
			Port:      getEnv("PORT", "8080"),
			Host:      getEnv("HOST", "localhost"),
			JWTSecret: getEnv("JWT_SECRET", "secret"),
		},
		DBConfig: config.DBConfig{
			DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth-api-db?sslmode=disable"),
		},
		S3Config: config.S3Config{
			AWSRegion:          getEnv("AWS_REGION", ""),
			AWSBucket:          getEnv("AWS_BUCKET", ""),
			AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		},
		GCPConfig: config.GCPConfig{
			OAuth2Config: config.OAuth2Config{
				GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				FrontendURL:        getEnv("FRONTEND_URL", ""),
			},
		},
		RabbitMQ: config.RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", ""),
		},
		Notification: config.NotificationConfig{
			Email: config.EmailConfig{
				Host:     getEnv("EMAIL_HOST", ""),
				Port:     getEnv("EMAIL_PORT", ""),
				Username: getEnv("EMAIL_HOST_USER", ""),
				Password: getEnv("EMAIL_HOST_PASSWORD", ""),
			},
		},
	}

	return configService, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
