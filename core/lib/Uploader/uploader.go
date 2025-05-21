package uploader

import (
	"errors"
	"os"
)

type Uploader interface {
	Save(fileType string, file []byte, filename string) (string, error)
	Delete(filename string) error
	Replace(fileType string, oldFilename string, newFile []byte, newFilename string) (string, error)
}

func Save(fileType string, file []byte, filename string) (string, error) {
	uploader, err := FileUploader()
	if err != nil {
		return "", err
	}
	return uploader.Save(fileType, file, filename)
}

func FileUploader() (Uploader, error) {
	switch os.Getenv("APP_ENV") {
	case "prod":
		return &S3Uploader{}, nil
	case "test", "dev", "":
		return &LocalUploader{}, nil
	default:
		return nil, errors.New("unknown APP_ENV")
	}
}
