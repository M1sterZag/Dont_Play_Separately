package platforms_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
)

func (s *PlatformsService) SearchPlatforms(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Platform, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	platforms, err := s.platformsRepository.SearchPlatforms(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(platforms) > 0 {
		return platforms, nil
	}

	igdbPlatforms, err := s.provider.SearchPlatforms(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search platforms in provider: %w", err)
	}

	platformDomains := make([]domain.Platform, 0, len(igdbPlatforms))
	for _, platform := range igdbPlatforms {
		platformDomain := markSynced(platform)
		if err := s.platformsRepository.UpsertPlatform(ctx, platformDomain); err != nil {
			return nil, fmt.Errorf("upsert platform %d: %w", platform.ID, err)
		}
		platformDomains = append(platformDomains, platformDomain)
	}
	return platformDomains, nil
}