package lib

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var root_dir = "mynute-go"

// FindProjectRoot searches for the `mynute-go` directory from the current working directory upwards.
func FindProjectRoot() (string, error) {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree until we find `mynute-go`
	for {
		if filepath.Base(dir) == root_dir {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// If we reach the root directory and didn't find it, return an error
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func ChangeWorkDirectoryTo(dir string) error {
	return os.Chdir(dir)
}

func IsRootDir() bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	return strings.HasSuffix(dir, root_dir)
}

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
