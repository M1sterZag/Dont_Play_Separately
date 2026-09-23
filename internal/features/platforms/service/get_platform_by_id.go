package platforms_service

import (
	"context"
	"errors"
	"fmt"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	"go.uber.org/zap"
)

func (s *PlatformsService) GetPlatformByID(ctx context.Context, platformID int64) (domain.Platform, error) {
	log := core_logger.FromContext(ctx)

	key := fmt.Sprintf("catalog:platform:%d", platformID)

	var cached domain.Platform
	if err := s.cache.Get(ctx, key, &cached); err == nil {
		return cached, nil
	} else if !errors.Is(err, core_cache.ErrNotFound) {
		log.Warn("cache get failed", zap.String("key", key), zap.Error(err))
	}

	v, err, _ := s.singleFlight.Do(key, func() (any, error) {
		platform, err := s.platformsRepository.GetPlatformByID(ctx, platformID)
		if err == nil {
			if err := s.cache.Set(ctx, key, platform, s.cacheTTL); err != nil {
				log.Warn("cache set failed", zap.String("key", key), zap.Error(err))
			}
			return platform, nil
		}
		if !errors.Is(err, core_errors.ErrNotFound) {
			return domain.Platform{}, err
		}

		platforms, err := s.provider.FetchPlatformsByIDs(ctx, []int64{platformID})
		if err != nil {
			return domain.Platform{}, fmt.Errorf("fetch platform %d from provider: %w", platformID, err)
		}
		if len(platforms) == 0 {
			return domain.Platform{}, core_errors.ErrNotFound
		}

		platform = markSynced(platforms[0])
		if err := s.platformsRepository.UpsertPlatform(ctx, platform); err != nil {
			return domain.Platform{}, fmt.Errorf("upsert platform %d: %w", platformID, err)
		}
		if err := s.cache.Set(ctx, key, platform, s.cacheTTL); err != nil {
			log.Warn("cache set failed", zap.String("key", key), zap.Error(err))
		}

		return platform, nil
	})
	if err != nil {
		return domain.Platform{}, err
	}

	platform, ok := v.(domain.Platform)
	if !ok {
		return domain.Platform{}, fmt.Errorf("unexpected singleflight result type: %T", v)
	}
	return platform, nil
}