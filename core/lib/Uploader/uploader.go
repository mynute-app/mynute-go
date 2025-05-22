package myUploader

import (
	"agenda-kaki-go/core/lib"
	"errors"
	"os"
)

type Uploader interface {
	Save(fileType string, file []byte, filename string) (string, error)
	Delete(filename string) error
	Replace(fileType string, oldFilename string, newFile []byte, newFilename string) (string, error)
}

func FileUploader(caller_entity string, caller_id string) (Uploader, error) {
	switch os.Getenv("APP_ENV") {
	case "prod":
		return NewS3Uploader(caller_entity, caller_id)
	case "test", "dev":
		return &LocalUploader{
			Entity:   caller_entity,
			EntityID: caller_id,
		}, nil
	default:
		return nil, lib.Error.General.InternalError.WithError(errors.New("unknown APP_ENV"))
	}
}
