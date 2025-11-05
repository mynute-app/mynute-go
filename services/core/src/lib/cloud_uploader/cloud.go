// File: uploader/cloud.go
package myUploader

import (
	"bytes"
	"context"
	"fmt"
	"mynute-go/services/core/src/lib"
	"net/http"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type cloudUploader struct {
	Entity    string
	EntityID  string
	Client    S3Uploader
	Bucket    string
	PublicURL string
}

func NewCloudUploader(entity, id string, client S3Uploader, bucket, publicURL string) *cloudUploader {
	return &cloudUploader{
		Entity:    entity,
		EntityID:  id,
		Client:    client,
		Bucket:    bucket,
		PublicURL: publicURL,
	}
}

func (c *cloudUploader) Save(fileType string, file []byte, originalFilename string) (string, error) {
	scopedPath := GenerateUniqueFilename(c.Entity, c.EntityID, originalFilename)
	strategy, err := getStrategy(c, fileType)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return strategy(file, scopedPath)
}

func (c *cloudUploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid file URL: %s", fileURL))
	}
	scopedPath := path.Join(c.Entity, c.EntityID, filename)
	_, err := c.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &c.Bucket,
		Key:    &scopedPath,
	})
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (c *cloudUploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if err := c.Delete(oldURL); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return c.Save(fileType, newFile, originalFilename)
}

func (c *cloudUploader) save(file []byte, scopedPath string) (string, error) {
	contentType := http.DetectContentType(file)
	_, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:             &c.Bucket,
		Key:                &scopedPath,
		Body:               bytes.NewReader(file),
		ACL:                types.ObjectCannedACLPublicRead,
		ContentDisposition: aws.String("inline"),
		ContentType:        aws.String(contentType),
	})
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	url := fmt.Sprintf("%s/%s/%s", c.PublicURL, c.Bucket, filepath.ToSlash(scopedPath))
	return url, nil
}
