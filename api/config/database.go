package config

import (
	"agenda-kaki-company-go/api/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	MigrateDB(db);

	return db
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&models.CompanyType{}, &models.Company{}, &models.Branch{}, &models.Employee{}, &models.Service{}, &models.Schedule{})
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to close the database: ", err)
	}
	sqlDB.Close()
}
