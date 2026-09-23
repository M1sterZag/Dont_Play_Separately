package platforms_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

func (r *PlatformsRepository) UpsertPlatform(ctx context.Context, platform domain.Platform) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO dps.platforms
	(id, title, abbreviation, slug, icon_url, checksum, updated_at, synced_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		abbreviation = EXCLUDED.abbreviation,
		slug = EXCLUDED.slug,
		icon_url = EXCLUDED.icon_url,
		checksum = EXCLUDED.checksum,
		updated_at = EXCLUDED.updated_at,
		synced_at = EXCLUDED.synced_at;
	`

	_, err := r.pool.Exec(ctx, query,
		platform.ID,
		platform.Title,
		platform.Abbreviation,
		platform.Slug,
		platform.IconURL,
		platform.Checksum,
		platform.UpdatedAt,
		platform.SyncedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert platform id='%d': %w", platform.ID, err)
	}

	return nil
}