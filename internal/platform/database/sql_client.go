package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tapiaw38/auth-api-be/internal/platform/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var sqlClient *sql.DB

func initSQLClient() error {
	cfg := config.GetConfigService()
	newSQLClient, err := sql.Open("postgres", cfg.DBConfig.DatabaseURL)
	if err != nil {
		return err
	}

	if err = newSQLClient.Ping(); err != nil {
		return err
	}

	sqlClient = newSQLClient
	return nil
}

func GetSQLClientInstance() (*sql.DB, error) {
	if sqlClient == nil {
		if err := initSQLClient(); err != nil {
			return nil, err
		}
	}

	return sqlClient, nil
}

func getRelativePathToMigrationsDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absMigrationsDirPath := filepath.Join(cwd, "migrations")

	relMigrationsDirPath, err := filepath.Rel(cwd, absMigrationsDirPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("file://%s", relMigrationsDirPath), nil
}

func Makemigration() error {
	cfg := config.GetConfigService()
	migrationPath, err := getRelativePathToMigrationsDirectory()
	if err != nil {
		return err
	}

	m, err := migrate.New(migrationPath, cfg.DBConfig.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migrations: error closing source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migrations: error closing database: %v", dbErr)
		}
	}()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	log.Printf("migrations: current version is %v (dirty: %v)", version, dirty)

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("migrations: no new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("migrations: database migrated successfully")

	return nil
}
