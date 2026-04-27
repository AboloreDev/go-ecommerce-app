package providers

import (
	"context"

	appConfig "github.com/aboloredev/armory/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func CreateAWSConfig(ctx context.Context, appConfig *appConfig.AWSConfig) (aws.Config, error) {
	var cfg aws.Config
	var err error

	if appConfig.S3Endpoint != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(appConfig.Region),
			config.WithBaseEndpoint(appConfig.S3Endpoint),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				appConfig.KeyId,
				appConfig.Key,
				"",
			)))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(appConfig.Region))
	}

	return cfg, err
}
