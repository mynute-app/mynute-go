package lib

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if !IsRootDir() {
		root, err := FindProjectRoot()
		if err != nil {
			log.Fatalf("Failed to find project root: %v", err)
		}
		if err := ChangeWorkDirectoryTo(root); err != nil {
			log.Fatalf("Failed to change working directory: %v", err)
		}
	}

	// First try to load root .env for common variables
	_ = godotenv.Load()

	// Then override with service-specific .env
	err := godotenv.Overload("services/auth/.env")

	// In Docker environments, .env file won't exist
	// Variables will be injected by Docker Compose
	if err != nil {
		log.Println("INFO: .env file not found, proceeding with system-provided environment variables. This is expected in a container environment.")
	} else {
		log.Println("INFO: .env file loaded successfully. This is expected in a local environment.")
	}
}
