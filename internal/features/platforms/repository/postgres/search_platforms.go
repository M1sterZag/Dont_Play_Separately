package platforms_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	platforms_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/platforms/repository"
)

func (r *PlatformsRepository) SearchPlatforms(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Platform, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, abbreviation, slug, icon_url, checksum, updated_at, synced_at
	FROM dps.platforms
	WHERE title ILIKE '%' || $1 || '%'
	ORDER BY title ASC, id ASC
	OFFSET $2
	LIMIT $3;
	`

	rows, err := r.pool.Query(ctx, query, searchQuery, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("query search platforms: %w", err)
	}
	defer rows.Close()

	var platformModels []platforms_repository.PlatformModel

	for rows.Next() {
		var platformModel platforms_repository.PlatformModel
		err := rows.Scan(
			&platformModel.ID,
			&platformModel.Title,
			&platformModel.Abbreviation,
			&platformModel.Slug,
			&platformModel.IconURL,
			&platformModel.Checksum,
			&platformModel.UpdatedAt,
			&platformModel.SyncedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan platforms: %w", err)
		}

		platformModels = append(platformModels, platformModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	platformsDomains := platforms_repository.PlatformDomainsFromModels(platformModels)

	return platformsDomains, nil
}
