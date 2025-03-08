package database

import (
	"agenda-kaki-go/core/config/db/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
}

func Connect() *Database {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB_NAME")
	port := os.Getenv("POSTGRES_PORT")
	app_env := os.Getenv("APP_ENV")
	sslmode := "disable" // You can modify this based on your setup
	timeZone := "UTC"    // Default timezone

	if app_env == "test" {
		dbName = fmt.Sprintf("%s-%s", dbName, app_env)
	}

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslmode, timeZone)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Migrate the database schema

	return &Database{Gorm: db}
}

func (db *Database) Migrate() {
	models := []any{
		&model.Sector{},
		&model.Company{}, // Must be migrated before Service
		&model.Branch{},
		&model.Appointment{},
		&model.Holidays{},
		&model.User{},
		&model.Employee{},
		&model.Service{},
	}

	for _, model := range models {
		log.Printf("Migrating: %T", model)
		if err := db.Gorm.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate %T: %v", model, err)
		}
	}
	log.Println("Migration completed successfully")
}

func (db *Database) CloseDB() {
	sqlDB, err := db.Gorm.DB()
	if err != nil {
		log.Fatal("Failed to close the database: ", err)
	}
	sqlDB.Close()
}

func (db *Database) ClearDB() {
	// Avoid this in production
	if os.Getenv("APP_ENV") != "test" {
		log.Fatal("ClearDB should only be used in test environment")
	}

	models := []any{
		&model.Sector{},
		&model.Company{},
		&model.Branch{},
		&model.Appointment{},
		&model.Holidays{},
		&model.User{},
		&model.Employee{},
		&model.Service{},
	}

	for _, model := range models {
		log.Printf("Clearing: %T", model)
		if err := db.Gorm.Migrator().DropTable(model); err != nil {
			log.Fatalf("Failed to clear %T: %v", model, err)
		}
	}
	log.Println("Clearing completed successfully")
}
