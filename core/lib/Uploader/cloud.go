// File: uploader/cloud.go
package myUploader

import (
	"agenda-kaki-go/core/lib"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type CloudUploader struct {
	Entity    string
	EntityID  string
	Client    *s3.Client
	Bucket    string
	PublicURL string
}

func NewCloudUploader(entity, id string) (*CloudUploader, error) {
	driver := os.Getenv("STORAGE_DRIVER")
	bucket := os.Getenv("AWS_BUCKET")
	publicURL := os.Getenv("AWS_PUBLIC_URL")
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")

	if bucket == "" || publicURL == "" || region == "" {
		return nil, lib.Error.General.InternalError.WithError(fmt.Errorf("missing required env vars"))
	}

	var cfg aws.Config
	var err error

	if driver == "r2" {
		if endpoint == "" {
			return nil, lib.Error.General.InternalError.WithError(fmt.Errorf("missing AWS_ENDPOINT for R2"))
		}
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	}

	if err != nil {
		return nil, err
	}

	return &CloudUploader{
		Entity:    entity,
		EntityID:  id,
		Client:    s3.NewFromConfig(cfg),
		Bucket:    bucket,
		PublicURL: publicURL,
	}, nil
}

func (c *CloudUploader) Save(fileType string, file []byte, originalFilename string) (string, error) {
	scopedPath := GenerateUniqueFilename(c.Entity, c.EntityID, originalFilename)
	strategy, err := getStrategy(c, fileType)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return strategy(file, scopedPath)
}

func (c *CloudUploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid file URL: %s", fileURL))
	}
	scopedPath := filepath.Join(c.Entity, c.EntityID, filename)
	_, err := c.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &c.Bucket,
		Key:    &scopedPath,
	})
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (c *CloudUploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if err := c.Delete(oldURL); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return c.Save(fileType, newFile, originalFilename)
}

func (c *CloudUploader) save(file []byte, scopedPath string) (string, error) {
	_, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &c.Bucket,
		Key:    &scopedPath,
		Body:   bytes.NewReader(file),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	url := fmt.Sprintf("%s/%s", c.PublicURL, filepath.ToSlash(scopedPath))
	return url, nil
}