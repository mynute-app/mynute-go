// File: uploader/local.go
package uploader

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"fmt"
	"os"
	"path/filepath"
)

type LocalUploader struct {
	Entity   string
	EntityID string
}

func (l *LocalUploader) Save(fileType string, file []byte, originalFilename string) (string, error) {
	scopedPath := GenerateUniqueFilename(l.Entity, l.EntityID, originalFilename)
	strategy, err := getStrategy(l, fileType)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return strategy(file, scopedPath)
}

func (l *LocalUploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("invalid file URL: %s", fileURL))
	}
	fullPath := filepath.Join("./"+namespace.UploadsFolder, l.Entity, l.EntityID, filename)
	return os.Remove(fullPath)
}

func (l *LocalUploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if oldURL != "" {
		if err := l.Delete(oldURL); err != nil {
			return "", lib.Error.General.InternalError.WithError(err)
		}
	}
	return l.Save(fileType, newFile, originalFilename)
}

func (l *LocalUploader) save(file []byte, scopedPath string) (string, error) {
	root_path, err := lib.FindProjectRoot();
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	uploadDir := filepath.Join(root_path, namespace.UploadsFolder, filepath.Dir(scopedPath))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}

	fullPath := filepath.Join(root_path, namespace.UploadsFolder, scopedPath)
	if err := os.WriteFile(fullPath, file, 0644); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}

	return "/static/" + namespace.UploadsFolder + "/" + scopedPath, nil
}
