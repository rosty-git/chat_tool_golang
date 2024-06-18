package s3

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	s3PresignClient *s3.PresignClient
	s3Client        *s3.Client
	s3BucketName    string
}

func NewStorage(bucketName string) *Storage {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	s3PresignClient := s3.NewPresignClient(
		s3.NewFromConfig(cfg),
	)

	s3Client := s3.NewFromConfig(cfg)

	return &Storage{
		s3Client:        s3Client,
		s3PresignClient: s3PresignClient,
		s3BucketName:    bucketName,
	}
}

func (s *Storage) Delete(key string) error {
	var objectIds []types.ObjectIdentifier
	objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})

	_, err := s.s3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(s.s3BucketName),
		Delete: &types.Delete{Objects: objectIds},
	})

	return err
}

func (s *Storage) GetPresignUrl(key string) (string, error) {
	PresignUrl, err := s.s3PresignClient.PresignPutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.s3BucketName),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(time.Minute*15))
	if err != nil {
		return "", err
	}
	return PresignUrl.URL, err
}
