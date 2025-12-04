package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

// NewMigrate creates a new migrate instance
func NewMigrate(migrationsPath string) (*migrate.Migrate, error) {
	config := GetMigrationConfig()

	// Open database connection with explicit sslmode
	dbURL := config.GetDatabaseURL()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Convert Windows path to proper format
	// golang-migrate needs forward slashes
	migrationsPath = strings.ReplaceAll(migrationsPath, "\\", "/")

	// Use file:// with relative-style path (no triple slash for drive letter)
	// This works better on Windows: file://C:/path instead of file:///C:/path
	sourceURL := "file://" + migrationsPath

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// RunMigrations runs all pending migrations
func RunMigrations(migrationsPath string) error {
	m, err := NewMigrate(migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations have been applied yet")
	} else {
		log.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
	}

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(migrationsPath string, steps int) error {
	m, err := NewMigrate(migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	if steps <= 0 {
		steps = 1
	}

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Printf("Successfully rolled back %d migration(s)\n", steps)
	return nil
}

// MigrationVersion returns the current migration version
func MigrationVersion(migrationsPath string) (uint, bool, error) {
	m, err := NewMigrate(migrationsPath)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	return m.Version()
}

// ForceMigrationVersion sets the migration version without running migrations
func ForceMigrationVersion(migrationsPath string, version int) error {
	m, err := NewMigrate(migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Printf("Successfully forced migration version to %d\n", version)
	return nil
}
