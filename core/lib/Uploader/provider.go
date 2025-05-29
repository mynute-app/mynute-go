// File: uploader/provider.go
package myUploader

import (
	"agenda-kaki-go/core/config/cloud"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	if err := ensure_bucket_exists(client.Client, bucket); err != nil {
		return fmt.Errorf("failed to ensure bucket exists: %w", err)
	}
	if err := make_bucket_public(client.Client, bucket); err != nil {
		return fmt.Errorf("failed to make bucket public: %w", err)
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
		public_url := os.Getenv("MINIO_PUBLIC_URL")
		bucket := os.Getenv("MINIO_BUCKET")
		if bucket == "" || public_url == "" {
			return "", "", fmt.Errorf("MINIO_BUCKET or MINIO_PUBLIC_URL is not set")
		}
		return bucket, public_url, nil // Local bucket for MIN-IO
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

func ensure_bucket_exists(client *s3.Client, bucket string) error {
	_, err := client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucket,
	})
	if err != nil {
		_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: &bucket,
		})
	}
	return err
}

func make_bucket_public(client *s3.Client, bucket string) error {
	if os.Getenv("APP_ENV") == "prod" {
		return nil // Avoid using this function in production
	}
	// Check if the bucket already has a public policy
	policyOutput, err := client.GetBucketPolicy(context.TODO(), &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	})
	if err != nil {
		return fmt.Errorf("failed to get bucket policy: %w", err)
	}
	policy := cloud.PublicReadPolicy(bucket)
	if policyOutput.Policy != nil && *policyOutput.Policy != "" {
		if *policyOutput.Policy == policy {
			return nil // Bucket is already public
		}
	}
	_, err = client.PutBucketPolicy(context.TODO(), &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: &policy,
	})
	if err != nil {
		return fmt.Errorf("failed to make bucket public: %w", err)
	}
	return nil
}
