// File: uploader/s3_client.go
package myUploader

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader interface {
	PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}
