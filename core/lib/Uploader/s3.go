// File: uploader/s3.go
package myUploader

import (
	"agenda-kaki-go/core/lib"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	Entity    string
	EntityID  string
	Client    *s3.Client
	Bucket    string
	PublicURL string
}

func NewS3Uploader(entity string, id string) (*S3Uploader, error) {
	bucket := os.Getenv("AWS_BUCKET")
	publicURL := os.Getenv("AWS_PUBLIC_URL")
	if bucket == "" || publicURL == "" {
		return nil, lib.Error.General.InternalError.WithError(fmt.Errorf("missing AWS_BUCKET or AWS_PUBLIC_URL env var"))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return &S3Uploader{
		Entity:    entity,
		EntityID:  id,
		Client:    s3.NewFromConfig(cfg),
		Bucket:    bucket,
		PublicURL: publicURL,
	}, nil
}

func (s *S3Uploader) Save(fileType string, file []byte, originalFilename string) (string, error) {
	scopedPath := GenerateUniqueFilename(s.Entity, s.EntityID, originalFilename)
	strategy, err := getStrategy(s, fileType)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return strategy(file, scopedPath)
}

func (s *S3Uploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid file URL: %s", fileURL))
	}
	scopedPath := filepath.Join(s.Entity, s.EntityID, filename)
	_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &scopedPath,
	})
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (s *S3Uploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if err := s.Delete(oldURL); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return s.Save(fileType, newFile, originalFilename)
}

func (s *S3Uploader) save(file []byte, scopedPath string) (string, error) {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &scopedPath,
		Body:   bytes.NewReader(file),
	})
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	url := fmt.Sprintf("%s/%s", s.PublicURL, filepath.ToSlash(scopedPath))
	return url, nil
}
