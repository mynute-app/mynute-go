// This file has a function that will create a file inside /debug/output folder
// must create /debug/output folder in case it does not exist
package debug

import (
	"fmt"
	"os"
	"path/filepath"
)

// Output creates a file inside debug/output folder
// with the given name and writes the content into it
func Output(name string, content any) error {
	if os.Getenv("APP_ENV") != "test" {
		return nil
	}

	dir := filepath.Join("debug", "output")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, name+".txt")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%+v\n", content)); err != nil {
		return err
	}

	return nil
}

// Remove all files inside debug/output folder
func Clear() error {
	if os.Getenv("APP_ENV") != "test" {
		return nil
	}
	dir := filepath.Join("debug", "output")

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}
