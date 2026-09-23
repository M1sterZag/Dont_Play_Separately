package core_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client  *redis.Client
	timeout time.Duration
}

func NewRedisCache(ctx context.Context, config Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisCache{
		client:  client,
		timeout: config.Timeout,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	raw, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return core_cache.ErrNotFound
		}
		return fmt.Errorf("get %s: %w", key, err)
	}

	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("unmarshal %s: %w", key, err)
	}

	return nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", key, err)
	}

	if err := c.client.Set(ctx, key, raw, ttl).Err(); err != nil {
		return fmt.Errorf("set %s: %w", key, err)
	}

	return nil
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("delete %s: %w", keys, err)
	}

	return nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
