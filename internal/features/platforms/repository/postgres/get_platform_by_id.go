package platforms_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	platforms_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/platforms/repository"
)

func (r *PlatformsRepository) GetPlatformByID(ctx context.Context, platformID int64) (domain.Platform, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, abbreviation, slug, icon_url, checksum, updated_at, synced_at
	FROM dps.platforms
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, platformID)

	var platformModel platforms_repository.PlatformModel
	err := row.Scan(
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
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.Platform{}, fmt.Errorf("find platform with id='%d': %w", platformID, core_errors.ErrNotFound)
		}

		return domain.Platform{}, fmt.Errorf("scan error: %w", err)
	}

	return platforms_repository.PlatformDomainFromModel(platformModel), nil
}