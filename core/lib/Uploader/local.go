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
	filename := ExtractFilenameFromURL(fileURL) // Extracts the base unique filename, e.g., "original_uuid.ext"
	if filename == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid file URL: %s", fileURL)) // Changed to BadRequest
	}

	root_path, err := lib.FindProjectRoot()
	if err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to find project root for deletion: %w", err))
	}

	// Construct the full path consistent with how `save` stores it.
	// `scopedPath` in `save` is `entity/entityID/name_uuid.ext`.
	// `filename` here is `name_uuid.ext`.
	// So, we need to join `root_path, namespace.UploadsFolder, l.Entity, l.EntityID, filename`.
	fullPath := filepath.Join(root_path, namespace.UploadsFolder, l.Entity, l.EntityID, filename)

	err = os.Remove(fullPath)
	if err != nil {
		// Return a wrapped error, potentially checking os.IsNotExist if desired
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to remove file %s: %w", fullPath, err))
	}
	return nil
}

func (l *LocalUploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if oldURL != "" {
		if err := l.Delete(oldURL); err != nil {
			// If Delete fails (e.g., file not found, or other FS error), wrap and return.
			// Consider if os.IsNotExist(err) from Delete should be handled differently (e.g., ignored).
			// For now, strict error propagation.
			return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to delete old file for replacement: %w", err))
		}
	}
	return l.Save(fileType, newFile, originalFilename)
}

func (l *LocalUploader) save(file []byte, scopedPath string) (string, error) {
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

	// Ensure forward slashes for URL compatibility
	urlPath := filepath.ToSlash(filepath.Join(namespace.UploadsFolder, scopedPath))
	return "/static/" + urlPath, nil
}