package providers

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appConfig "github.com/aboloredev/armory/internal/config"
)

type s3Provider struct {
	client   *s3.Client
	uploader *manager.Uploader
	endpoint string
	bucket   string
}

func NewS3Provider(cfg *appConfig.Config) *s3Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSServices.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSServices.KeyId,
			cfg.AWSServices.Key,
			"",
		)))
	if err != nil {
		panic("Failed to create AWS config" + err.Error())
	}

	// configure for localstack
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.AWSServices.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.AWSServices.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	return &s3Provider{
		client:   client,
		uploader: manager.NewUploader(client),
		bucket:   cfg.AWSServices.S3Bucket,
		endpoint: cfg.AWSServices.S3Endpoint,
	}
}

// Upload file
func (p *s3Provider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	result, err := p.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
		Body:   src,
	})

	if err != nil {
		return "", err
	}

	return *result.Key, nil
}

// Delete file
func (p *s3Provider) DeleteFile(path string) error {
	_, err := p.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}
