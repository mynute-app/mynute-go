package lib

import (
	"os"
	"path/filepath"
	"strings"
)

var root_dir = "agenda-kaki-go"

// findProjectRoot searches for the `agenda-kaki-go` directory from the current working directory upwards.
func FindProjectRoot() (string, error) {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree until we find `agenda-kaki-go`
	for {
		if strings.HasSuffix(dir, root_dir) {
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

func MatchPath(rulePath, actualPath string) bool {
	ruleSegments := strings.Split(rulePath, "/")
	actualSegments := strings.Split(actualPath, "/")

	if len(ruleSegments) != len(actualSegments) {
		return false
	}

	for i := range ruleSegments {
		if strings.HasPrefix(ruleSegments[i], ":") {
			continue
		}
		if ruleSegments[i] != actualSegments[i] {
			return false
		}
	}
	return true
}