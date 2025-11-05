package myUploader

import (
	"errors"
	"mynute-go/services/core/api/lib"
)

type Uploader interface {
	Save(fileType string, file []byte, filename string) (string, error)
	Delete(filename string) error
	Replace(fileType string, oldURL string, newFile []byte, newFilename string) (string, error)
}

func FileUploader(caller_entity string, caller_id string) (Uploader, error) {
	if factory == nil {
		return nil, lib.Error.General.InternalError.WithError(errors.New("uploader factory not initialized"))
	}
	return factory(caller_entity, caller_id)
}
