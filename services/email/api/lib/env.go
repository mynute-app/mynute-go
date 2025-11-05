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

	// Try loading from service directory first, then from current directory
	if err := godotenv.Load("services/email/.env"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("INFO: .env file not found, proceeding with system-provided environment variables. This is expected in a container environment.")
		} else {
			log.Println(".env file loaded successfully")
		}
	} else {
		log.Println(".env file loaded successfully from services/email/.env")
	}
}
