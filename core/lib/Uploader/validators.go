package uploader

import (
	"errors"
	"net/http"
	"strings"
)

func validateImage(file []byte) error {
	contentType := http.DetectContentType(file)
	if strings.HasPrefix(contentType, "image/") {
		return nil
	}
	return errors.New("file is not a valid image")
}

func validatePDF(file []byte) error {
	contentType := http.DetectContentType(file)
	if contentType == "application/pdf" {
		return nil
	}
	return errors.New("file is not a valid PDF")
}