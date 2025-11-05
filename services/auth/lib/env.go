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

	// Try loading from service directory first, then from root
	err := godotenv.Load("services/auth/.env")
	if err != nil {
		err = godotenv.Load()
	}

	// In Docker environments, .env file won't exist
	// Variables will be injected by Docker Compose
	if err != nil {
		log.Println("INFO: .env file not found, proceeding with system-provided environment variables. This is expected in a container environment.")
	} else {
		log.Println("INFO: .env file loaded successfully. This is expected in a local environment.")
	}
}
