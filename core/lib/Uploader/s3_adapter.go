// File: uploader/s3_adapter.go
package myUploader

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RealS3Client struct {
	Client *s3.Client
}

func (r *RealS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return r.Client.PutObject(ctx, input)
}

func (r *RealS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return r.Client.DeleteObject(ctx, input)
}
