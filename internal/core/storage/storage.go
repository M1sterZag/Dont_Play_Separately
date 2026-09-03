package core_storage

import (
	"context"
	"io"
)

type Storage interface {
	PublicURL(key string) string
	PutObject(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
}
