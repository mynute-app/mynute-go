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

	err := godotenv.Load()
	// NO DOCKER: err não será nil, pois o arquivo .env não existe.
	// Nós não tratamos isso como um erro fatal.
	// A aplicação continuará, usando as variáveis injetadas pelo Docker Compose.
	if err != nil {
		log.Println("INFO: .env file not found, proceeding with system-provided environment variables. This is expected in a container environment.")
	} else {
		log.Println("INFO: .env file loaded successfully. This is expected in a local environment.")
	}
}

