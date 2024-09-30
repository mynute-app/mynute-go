package config

import (
	"agenda-kaki-go/core/config/db/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
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
	sslmode := "disable" // You can modify this based on your setup
	timeZone := "UTC"    // Default timezone

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslmode, timeZone)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Migrate the database schema
	MigrateDB(db)

	return db
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&models.Company{}, &models.CompanyType{}, &models.Branch{}, &models.Employee{}, &models.Service{}, &models.Schedule{})
	if err != nil {
		log.Fatal("Failed to migrate the database: ", err)
	}
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to close the database: ", err)
	}
	sqlDB.Close()
}
