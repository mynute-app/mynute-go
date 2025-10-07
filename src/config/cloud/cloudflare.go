package cloud

import (
	"context"
	"fmt"
	"mynute-go/src/lib"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CloudFlare struct{}

func (cf *CloudFlare) R2() (*s3.Client, error) {
	var accountId = os.Getenv("R2_ACCOUNT_ID")
	var accessKeyId = os.Getenv("R2_ACCESS_KEY_ID")
	var accessKeySecret = os.Getenv("R2_ACCESS_KEY_SECRET")

	if accountId == "" || accessKeyId == "" || accessKeySecret == "" {
		return nil, lib.Error.General.InternalError.WithError(fmt.Errorf("missing required env vars for Cloudflare R2"))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, lib.Error.General.InternalError.WithError(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return client, nil
}
