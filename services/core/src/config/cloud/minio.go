package cloud

import (
	"context"
	"fmt"
	"mynute-go/services/core/src/lib"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MinIO struct{}

func (m *MinIO) Client() (*s3.Client, error) {
	public_url := os.Getenv("MINIO_PUBLIC_URL")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")

	if public_url == "" || accessKey == "" || secretKey == "" {
		return nil, lib.Error.General.InternalError.WithError(
			fmt.Errorf("missing required env vars for MinIO"),
		)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, lib.Error.General.InternalError.WithError(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(public_url)
		o.UsePathStyle = true
	})

	return client, nil
}
