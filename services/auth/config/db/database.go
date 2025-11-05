package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	AuthDB *gorm.DB // Auth database connection
	Error  error
}

type Test struct {
	*Database
	name string
}

// Connect establishes a connection to the auth database
func Connect() *Database {
	// Get environment variables
	host := os.Getenv("POSTGRES_AUTH_HOST")
	user := os.Getenv("POSTGRES_AUTH_USER")
	password := os.Getenv("POSTGRES_AUTH_PASSWORD")
	dbName := os.Getenv("POSTGRES_AUTH_DB")
	port := os.Getenv("POSTGRES_AUTH_PORT")

	// Fallback to main DB connection settings if auth-specific vars not set
	if host == "" {
		host = os.Getenv("POSTGRES_HOST")
	}
	if user == "" {
		user = os.Getenv("POSTGRES_USER")
	}
	if password == "" {
		password = os.Getenv("POSTGRES_PASSWORD")
	}
	if port == "" {
		port = os.Getenv("POSTGRES_PORT")
	}

	app_env := os.Getenv("APP_ENV")
	db_log_level := os.Getenv("POSTGRES_LOG_LEVEL")
	LogLevel := logger.Warn

	// Set auth database name based on environment
	if dbName == "" {
		switch app_env {
		case "test":
			dbName = os.Getenv("POSTGRES_DB_TEST")
			if dbName == "" {
				dbName = "mynute_auth_test"
			}
			LogLevel = logger.Info
		case "dev":
			dbName = "mynute_auth_dev"
			LogLevel = logger.Warn
		default:
			dbName = "mynute_auth"
		}
	}

	switch db_log_level {
	case "info":
		LogLevel = logger.Info
	case "error":
		LogLevel = logger.Error
	case "silent":
		LogLevel = logger.Silent
	case "warn":
		LogLevel = logger.Warn
	}

	log.Printf("Auth Service - Running in %s environment. Database: %s\n", app_env, dbName)

	// Build the DSN (Data Source Name)
	sslmode := "disable"
	timeZone := "UTC"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslmode, timeZone)

	customGormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  LogLevel,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		},
	)

	gormConfig := &gorm.Config{
		Logger: customGormLogger,
	}

	// Connect to the auth database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the auth database: ", err)
	}

	// Set the connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection pool: ", err)
	}

	sqlDB.SetMaxIdleConns(10)                  // Max number of idle connections in the pool
	sqlDB.SetMaxOpenConns(50)                  // Max number of open connections to the database
	sqlDB.SetConnMaxLifetime(15 * time.Minute) // Max lifetime of a connection in the pool
	sqlDB.SetConnMaxIdleTime(2 * time.Second)  // Max idle time for a connection in the pool

	dbWrapper := &Database{
		AuthDB: db,
		Error:  nil,
	}

	if app_env == "test" {
		dbWrapper.Test().Clear()
	}

	return dbWrapper
}

// Migrate runs database migrations for auth models
func (d *Database) Migrate(models []interface{}) error {
	return d.AuthDB.AutoMigrate(models...)
}

// WithDB allows using a specific database connection
func (d *Database) WithDB(db *gorm.DB) *Database {
	return &Database{
		AuthDB: db,
		Error:  d.Error,
	}
}

// Disconnect closes the database connection
func (d *Database) Disconnect() {
	sqlDB, err := d.AuthDB.DB()
	if err != nil {
		log.Println("Failed to get database connection for closing:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
	}
}

// Test returns a Test instance for testing utilities
func (d *Database) Test() *Test {
	return &Test{
		Database: d,
		name:     "auth-test",
	}
}

// Clear clears all data from auth tables (for testing)
func (t *Test) Clear() {
	log.Println("Clearing auth test database...")

	// Delete all records from auth tables
	t.AuthDB.Exec("DELETE FROM admin_roles")
	t.AuthDB.Exec("DELETE FROM admins")
	t.AuthDB.Exec("DELETE FROM policy_rules")
	t.AuthDB.Exec("DELETE FROM resources")
	t.AuthDB.Exec("DELETE FROM endpoints")

	log.Println("Auth test database cleared")
}

// InitialSeed runs initial seeding for auth database
func (d *Database) InitialSeed() {
	// Seed admin-related data
	// This should be implemented based on your auth seeding requirements
	log.Println("Auth database seeding completed (implement as needed)")
}
