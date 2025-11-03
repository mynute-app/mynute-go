package myUploader

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// GenerateUniqueFilename generates a scoped path for storing an uploaded file,
// placing it in a structured directory based on entity type and ID.
// The resulting filename includes a UUID for uniqueness.
//
// Example:
//
//	entity     = "client"
//	entityID   = 123e4567-e89b-12d3-a456-426614174000
//	filename   = "photo.png"
//
//	Output:
//	  "client/123e4567-e89b-12d3-a456-426614174000/photo_dccf28b1-0ff2-44c9-bcfd-ccf05e6f6a61.png"
func GenerateUniqueFilename(entity string, entityID string, originalFilename string) string {
	ext := filepath.Ext(originalFilename)                            // .png
	name := strings.TrimSuffix(filepath.Base(originalFilename), ext) // photo
	unique := uuid.New().String()                                    // e.g., dccf28b1-0ff2-44c9-bcfd-ccf05e6f6a61
	return fmt.Sprintf("%s/%s/%s_%s%s", entity, entityID, name, unique, ext)
}

// Extrai o nome do arquivo da URL (para deletar do S3/local)
func ExtractFilenameFromURL(fileURL string) string {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		return ""
	}
	return path.Base(parsed.Path)
}

