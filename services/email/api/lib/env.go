package lib

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	app_env := os.Getenv("APP_ENV")
	if app_env == "prod" {
		log.Println("Production environment detected. Skipping .env file loading.")
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	} else {
		log.Println(".env file loaded successfully")
	}
}
