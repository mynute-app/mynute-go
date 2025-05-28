// File: uploader/provider.go
package myUploader

import (
	"agenda-kaki-go/core/config/cloud"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var factory func(entity, entityID string) (Uploader, error)

func StartProvider() error {
	switch os.Getenv("APP_ENV") {
	case "prod":
		client, err := get_s3_client()
		if err != nil {
			return fmt.Errorf("failed to get S3 client: %w", err)
		}
		bucket, publicURL, err := get_s3_info()
		if err != nil {
			return fmt.Errorf("failed to get S3 info: %w", err)
		}
		factory = func(entity, entityID string) (Uploader, error) {
			return NewCloudUploader(entity, entityID, client, bucket, publicURL), nil
		}
	case "test", "dev":
		factory = func(entity, entityID string) (Uploader, error) {
			return NewLocalUploader(entity, entityID), nil
		}
	default:
		return fmt.Errorf("unknown APP_ENV for uploader")
	}

	return nil
}

func get_s3_client() (*s3.Client, error) {
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "R2" {
		cloudflare := &cloud.CloudFlare{}
		client, err := cloudflare.R2()
		if err != nil {
			return nil, fmt.Errorf("failed to init R2 client: %w", err)
		}
		return client, nil
	}
	return nil, fmt.Errorf("unsupported storage driver: %s", driver)
}

func get_s3_info() (string, string, error) {
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "R2" {
		bucket := os.Getenv("R2_BUCKET")
		publicURL := os.Getenv("R2_PUBLIC_URL")
		if bucket == "" || publicURL == "" {
			return "", "", fmt.Errorf("R2_BUCKET or R2_PUBLIC_URL is not set")
		}
		return bucket, publicURL, nil
	}
	return "", "", fmt.Errorf("unsupported storage driver: %s", driver)
}
