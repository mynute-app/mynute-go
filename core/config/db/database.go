package database

import (
	"agenda-kaki-go/core/config/db/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
}

type Test struct {
	*Database
	name string
}

var models = []any{
	&model.Sector{},
	&model.Company{}, // Must be migrated before Service
	&model.Branch{},
	&model.Appointment{},
	&model.Holidays{},
	&model.Client{},
	&model.Employee{},
	&model.Service{},
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

	fmt.Printf("Running in %s environment. Database: %s\n", app_env, dbName)

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

// Migrate the database schema
func (db *Database) Migrate() {
	for _, model := range models {
		log.Printf("Migrating: %T", model)
		if err := db.Gorm.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate %T: %v", model, err)
		}
	}
	log.Println("Migration completed successfully")
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

// Clear the database. Only in test environment
func (t *Test) Clear() {
	if os.Getenv("APP_ENV") != "test" {
		return
	}
	// for _, model := range models {
	// 	log.Printf("Clearing: %T", model)
	// 	if err := t.Gorm.Migrator().DropTable(model); err != nil {
	// 		log.Fatalf("Failed to clear %T: %v", model, err)
	// 	}
	// }
	query := `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`
	if err := t.Gorm.Exec(query).Error; err != nil {
		log.Fatalf("Failed to clear database: %v", err)
	}
	fmt.Printf("Erased all tables on %s database.\n", t.name)
}
