package uploader

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

// Gera nome de arquivo com UUID e mesma extensão
func GenerateUniqueFilename(original string) string {
	ext := filepath.Ext(original)
	id := uuid.New().String()
	return fmt.Sprintf("%s%s", id, ext)
}

// Extrai o nome do arquivo da URL (para deletar do S3/local)
func ExtractFilenameFromURL(fileURL string) string {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		return ""
	}
	return path.Base(parsed.Path)
}