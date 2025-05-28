package myUploader

import (
	"errors"
	"fmt"
)

type UploadStrategy func(file []byte, filename string) (string, error)

var strategies = map[string]func(Uploader) UploadStrategy{
	"image": makeImageStrategy,
	"pdf":   makePDFStrategy,
}

func getStrategy(u Uploader, fileType string) (UploadStrategy, error) {
	if fn, ok := strategies[fileType]; ok {
		return fn(u), nil
	}
	return nil, fmt.Errorf("unsupported file type: %s", fileType)
}

func makeImageStrategy(u Uploader) UploadStrategy {
	return func(file []byte, originalFilename string) (string, error) {
		if err := validateImage(file); err != nil {
			return "", err
		}
		switch up := u.(type) {
		case *cloudUploader:
			return up.save(file, originalFilename)
		default:
			return "", errors.New("unknown uploader type")
		}
	}
}

func makePDFStrategy(u Uploader) UploadStrategy {
	return func(file []byte, originalFilename string) (string, error) {
		if err := validatePDF(file); err != nil {
			return "", err
		}
		switch up := u.(type) {
		case *cloudUploader:
			return up.save(file, originalFilename)
		default:
			return "", errors.New("unknown uploader type")
		}
	}
}
