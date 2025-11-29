package s3utils

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/types"
)

func NewS3Client(ctx context.Context, region string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region
	})

	return s3client, nil
}

func PresignUrlPutObject(ctx context.Context, s3client *s3.Client, inputCfg types.PresignUrlPutObjectInput) (string, error) {
	// Check bucket existence
	if _, err := s3client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(inputCfg.BucketName),
	}); err != nil {
		return "", fmt.Errorf("bucket '%s' not exist or no permission: %w", inputCfg.BucketName, err)
	}

	// Generate presigned URL
	presignClient := s3.NewPresignClient(s3client)

	key := fmt.Sprintf(constants.LogfileNamePattern, inputCfg.Region, inputCfg.NodeName, time.Now().UTC().Format(time.RFC3339))
	req, err := presignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(inputCfg.BucketName),
			Key:    aws.String(key),
		},
		func(o *s3.PresignOptions) {
			o.Expires = time.Duration(inputCfg.ExpiredSeconds) * time.Second
		})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return req.URL, nil
}
