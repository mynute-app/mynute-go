// File: uploader/local.go
package myUploader

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type localUploader struct {
	Entity   string
	EntityID string
}

func NewLocalUploader(entity, entityID string) *localUploader {
	return &localUploader{
		Entity:   entity,
		EntityID: entityID,
	}
}

func (l *localUploader) Save(fileType string, file []byte, originalFilename string) (string, error) {
	scopedPath := GenerateUniqueFilename(l.Entity, l.EntityID, originalFilename)
	strategy, err := getStrategy(l, fileType)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return strategy(file, scopedPath)
}

func (l *localUploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid file URL: %s", fileURL))
	}

	root_path, err := lib.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root for deletion: %w", err) // remover o ErrorStruct aqui para facilitar o retry
	}

	fullPath := filepath.Join(root_path, namespace.UploadsFolder, l.Entity, l.EntityID, filename)

	err = os.Remove(fullPath)
	if err != nil {
		// Aqui não envolve o erro com ErrorStruct, só retorna direto pro Replace() poder inspecionar com os.IsPermission()
		return fmt.Errorf("failed to remove file %s: %w", fullPath, err)
	}

	return nil
}

func (l *localUploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if oldURL != "" {
		var err error
		const maxAttempts = 12
		for i := 1; i <= maxAttempts; i++ {
			err = l.Delete(oldURL)
			if err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond) // Espera antes da próxima tentativa, sempre
		}
		if err != nil {
			return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to delete old file for replacement: %w", err))
		}
	}

	return l.Save(fileType, newFile, originalFilename)
}

func (l *localUploader) save(file []byte, scopedPath string) (string, error) {
	root_path, err := lib.FindProjectRoot()
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	// scopedPath is already "entity/entityID/filename_uuid.ext"
	uploadDir := filepath.Join(root_path, namespace.UploadsFolder, filepath.Dir(scopedPath))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to create directory %s: %w", uploadDir, err))
	}

	fullPath := filepath.Join(root_path, namespace.UploadsFolder, scopedPath)
	if err := os.WriteFile(fullPath, file, 0644); err != nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to write file %s: %w", fullPath, err))
	}

	return "/static/" + filepath.ToSlash(scopedPath), nil
}
