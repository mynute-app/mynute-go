package lib

import (

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
	if err != nil {
		log.Fatal(err)
	}
}
