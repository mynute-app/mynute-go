package lib

import (
	"fmt"
	"log"
	"os"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDialector creates a GORM Postgres dialector from DSN
func PostgresDialector(dsn string) gorm.Dialector {
	return gormpostgres.Open(dsn)
}

type MigrationConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetMigrationConfig returns the migration configuration from environment variables
// IMPORTANT: Migrations always use POSTGRES_DB_PROD (production database).
// For dev/test migrations, set POSTGRES_DB_PROD=devdb or POSTGRES_DB_PROD=testdb in your .env
func GetMigrationConfig() *MigrationConfig {
	dbName := os.Getenv("POSTGRES_DB_PROD")

	if dbName == "" {
		log.Fatal("POSTGRES_DB_PROD environment variable is required for migrations")
	}

	return &MigrationConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   dbName,
		SSLMode:  "disable",
	}
}

// GetDatabaseURL returns the database connection URL
func (c *MigrationConfig) GetDatabaseURL() string {
	// Ensure sslmode is properly set in the connection string
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, sslMode)
}

// Note: Migration management is now handled by Atlas (https://atlasgo.io/)
// Use the following commands:
//   - make migrate-up: Apply pending migrations
//   - make migrate-diff NAME=<name>: Generate new migration
//   - make migrate-status: Check migration status
//
// For manual Atlas usage:
//   - atlas migrate apply --env dev/prod
//   - atlas migrate diff <name> --env dev
//   - atlas migrate status --env dev/prod
