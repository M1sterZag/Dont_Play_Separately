package core_storage_minio

import (
	"context"
	"fmt"
	"io"

	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	config core_storage.Config
	client *minio.Client
}

func NewStorage(ctx context.Context, config core_storage.Config) (*MinioStorage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return &MinioStorage{}, fmt.Errorf("create minio client: %w", err)
	}

	if _, err := client.BucketExists(ctx, config.Bucket); err != nil {
		return nil, fmt.Errorf("minio unavalible: %w", err)
	}

	return &MinioStorage{
		config: config,
		client: client,
	}, nil
}

func (s *MinioStorage) PublicURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.config.PublicBaseURL, s.config.Bucket, key)
}

func (s *MinioStorage) PutObject(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	if _, err := s.client.PutObject(ctx, s.config.Bucket, key, r, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return fmt.Errorf("failed put object in s3: %w", err)
	}

	return nil
}

func (s *MinioStorage) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.config.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from s3: %w", err)
	}

	return object, nil
}
