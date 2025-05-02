package database

import (
	"agenda-kaki-go/core/config/db/model"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Gorm *gorm.DB
}

type Test struct {
	*Database
	name string
}

// Connects to the database
func Connect() *Database {
	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	app_env := os.Getenv("APP_ENV")
	sslmode := "disable" // You can modify this based on your setup
	timeZone := "UTC"    // Default timezone

	if app_env == "test" {
		dbName = os.Getenv("POSTGRES_DB_TEST")
	} else if app_env != "production" && app_env != "dev" {
		log.Fatalf("Invalid APP_ENV: %s", app_env)
	}

	log.Printf("Running in %s environment. Database: %s\n", app_env, dbName)

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslmode, timeZone)

	customGormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		},
	)

	gormConfig := &gorm.Config{
		Logger: customGormLogger,
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Migrate the database schema

	return &Database{Gorm: db}
}

// Migrate the database schema
func (db *Database) Migrate() *Database {
	for _, model := range model.GeneralModels {
		if err := db.Gorm.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate %T: %v", model, err)
		}
	}
	log.Println("Migration finished!")
	return db
}

func (db *Database) Seed() *Database {
	return db
}

// Close connection to the database
func (db *Database) Disconnect() {
	sqlDB, err := db.Gorm.DB()
	if err != nil {
		log.Fatal("Failed to close the database: ", err)
	}
	sqlDB.Close()
}
func (db *Database) Test() *Test {
	dbName := os.Getenv("POSTGRES_DB_NAME")
	app_env := os.Getenv("APP_ENV")
	dbName = fmt.Sprintf("%s-%s", dbName, app_env)
	return &Test{Database: db, name: dbName}
}

func (t *Test) Clear() {
	if os.Getenv("APP_ENV") != "test" {
		return
	}

	// Step 1: Drop all schemas except 'public'
	dropSchemasSQL := `
		DO $$ DECLARE
			schema_name text;
		BEGIN
			FOR schema_name IN
				SELECT nspname FROM pg_namespace
				WHERE nspname NOT IN ('pg_catalog', 'information_schema', 'public')
				  AND nspname NOT LIKE 'pg_toast%'
			LOOP
				EXECUTE format('DROP SCHEMA IF EXISTS %I CASCADE', schema_name);
			END LOOP;
		END $$;
	`

	// Step 2: Drop and recreate 'public' just in case
	resetPublicSQL := `
		DROP SCHEMA IF EXISTS public CASCADE;
		CREATE SCHEMA public;
	`

	// Execute both
	if err := t.Gorm.Exec(dropSchemasSQL).Error; err != nil {
		log.Fatalf("Failed to drop non-public schemas: %v", err)
	}
	if err := t.Gorm.Exec(resetPublicSQL).Error; err != nil {
		log.Fatalf("Failed to reset public schema: %v", err)
	}

	log.Printf("Erased all schemas on %s database.\n", t.name)
}
