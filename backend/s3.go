package backend

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3BackendConfig struct {
	Region   string
	Endpoint string
	Bucket   string

	Id     string
	Secret string
}

type S3Backend struct {
	config S3BackendConfig
	client *s3.S3
}

func NewS3Backend(config *S3BackendConfig) *S3Backend {
	var s3backend = S3Backend{config: *config}

	aws_config := aws.NewConfig()

	aws_config.WithRegion(s3backend.config.Region)
	aws_config.WithEndpoint(s3backend.config.Endpoint)

	aws_config.WithCredentials(credentials.NewStaticCredentials(
		s3backend.config.Id,
		s3backend.config.Secret,
		""))

	s3backend.client = s3.New(session.Must(session.NewSession()), aws_config)

	return &s3backend
}

func (s3backend *S3Backend) Read(path string) ([]byte, error) {
	output, err := s3backend.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3backend.config.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s3backend *S3Backend) Write(path string, data []byte) error {
	_, err := s3backend.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3backend.config.Bucket),
		Key:    aws.String(path),
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
	})
	return err
}
