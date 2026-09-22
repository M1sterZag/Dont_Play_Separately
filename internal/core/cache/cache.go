package core_cache

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("cache: value not found")

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}
