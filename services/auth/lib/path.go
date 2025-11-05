package lib

import (
	"os"
	"path/filepath"
)

// FindProjectRoot finds the root directory of the project by looking for go.mod
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root of filesystem
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

// IsRootDir checks if the current directory is the project root
func IsRootDir() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

// ChangeWorkDirectoryTo changes the current working directory
func ChangeWorkDirectoryTo(path string) error {
	return os.Chdir(path)
}
