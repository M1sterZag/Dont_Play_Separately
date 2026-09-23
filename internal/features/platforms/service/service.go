package platforms_service

import (
	"context"
	"time"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"golang.org/x/sync/singleflight"
)

type PlatformsRepository interface {
	GetPlatformByID(ctx context.Context, platformID int64) (domain.Platform, error)
	SearchPlatforms(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Platform, error)
	UpsertPlatform(ctx context.Context, platform domain.Platform) error
}

type PlatformsProvider interface {
	FetchPlatformsByIDs(ctx context.Context, platformIDs []int64) ([]domain.Platform, error)
	SearchPlatforms(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Platform, error)
}

type PlatformsService struct {
	platformsRepository PlatformsRepository
	cache               core_cache.Cache
	provider            PlatformsProvider

	singleFlight singleflight.Group
	cacheTTL     time.Duration
}

func NewPlatformsService(
	platformsRepository PlatformsRepository,
	cache core_cache.Cache,
	provider PlatformsProvider,
	cacheTTL time.Duration,
) *PlatformsService {
	return &PlatformsService{
		platformsRepository: platformsRepository,
		cache:               cache,
		provider:            provider,
		cacheTTL:            cacheTTL,
	}
}

func markSynced(platform domain.Platform) domain.Platform {
	now := time.Now()
	platform.SyncedAt = &now
	return platform
}
