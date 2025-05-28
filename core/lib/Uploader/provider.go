// File: uploader/provider.go
package myUploader

import (
	"agenda-kaki-go/core/config/cloud"
	"fmt"
	"os"
)

var factory func(entity, entityID string) (Uploader, error)

func StartProvider() error {
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
	return nil
}

func get_s3_client() (*myS3Client, error) {
	env := os.Getenv("APP_ENV")
	if env == "dev" || env == "test" {
		minIO := &cloud.MinIO{}
		client, err := minIO.Client()
		if err != nil {
			return nil, fmt.Errorf("failed to init MIN-IO client: %w", err)
		}
		return &myS3Client{Client: client}, nil // Local client using MIN-IO
	} else if env != "prod" {
		return nil, fmt.Errorf("unsupported APP_ENV: %s", env)
	}
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "R2" {
		cloudflare := &cloud.CloudFlare{}
		client, err := cloudflare.R2()
		if err != nil {
			return nil, fmt.Errorf("failed to init R2 client: %w", err)
		}
		return &myS3Client{Client: client}, nil
	}
	return nil, fmt.Errorf("unsupported storage driver: %s", driver)
}

func get_s3_info() (string, string, error) {
	env := os.Getenv("APP_ENV")
	if env == "dev" || env == "test" {
		return "local-bucket", "http://localhost:9000", nil // Local bucket for MIN-IO
	} else if env != "prod" {
		return "", "", fmt.Errorf("unsupported APP_ENV: %s", env)
	}
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
