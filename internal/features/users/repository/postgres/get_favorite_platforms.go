package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	users_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/repository"
	"github.com/google/uuid"
)

func (r *UsersRepository) GetFavoritePlatforms(ctx context.Context, userID uuid.UUID) ([]domain.Platform, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT p.id, p.title, p.abbreviation, p.slug, p.icon_url, p.checksum, p.updated_at, p.synced_at
	FROM dps.user_platforms up
	JOIN dps.platforms p ON p.id = up.platform_id
	WHERE up.user_id = $1
	ORDER BY p.title ASC, p.id ASC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query favorite platforms: %w", err)
	}
	defer rows.Close()

	var platformModels []users_repository.PlatformModel
	for rows.Next() {
		var platformModel users_repository.PlatformModel
		if err := rows.Scan(
			&platformModel.ID,
			&platformModel.Title,
			&platformModel.Abbreviation,
			&platformModel.Slug,
			&platformModel.IconURL,
			&platformModel.Checksum,
			&platformModel.UpdatedAt,
			&platformModel.SyncedAt,
		); err != nil {
			return nil, fmt.Errorf("scan favorite platform: %w", err)
		}

		platformModels = append(platformModels, platformModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	platformDomains := users_repository.PlatformDomainsFromModels(platformModels)

	return platformDomains, nil
}
