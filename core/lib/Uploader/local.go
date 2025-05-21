package uploader

import (
	"fmt"
	"os"
	"path/filepath"
)

type LocalUploader struct{}

func (l *LocalUploader) Save(fileType string, file []byte, originalFilename string) (string, error) {
	strategy, err := getStrategy(l, fileType)
	if err != nil {
		return "", err
	}
	return strategy(file, originalFilename)
}

func (l *LocalUploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return fmt.Errorf("invalid file URL: %s", fileURL)
	}
	path := filepath.Join("./uploads", filename)
	return os.Remove(path)
}

func (l *LocalUploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if err := l.Delete(oldURL); err != nil {
		return "", err
	}
	return l.Save(fileType, newFile, originalFilename)
}

// usada internamente pela strategy
func (l *LocalUploader) save(file []byte, originalFilename string) (string, error) {
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	uniqueName := GenerateUniqueFilename(originalFilename)
	path := filepath.Join(uploadDir, uniqueName)

	if err := os.WriteFile(path, file, 0644); err != nil {
		return "", err
	}

	return "/static/uploads/" + uniqueName, nil
}