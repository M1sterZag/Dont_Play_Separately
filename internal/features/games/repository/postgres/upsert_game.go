package games_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

func (r *GamesRepository) UpsertGame(ctx context.Context, game domain.Game) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO dps.games
	(id, title, slug, icon_url, checksum, updated_at, synced_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		slug = EXCLUDED.slug,
		icon_url = EXCLUDED.icon_url,
		checksum = EXCLUDED.checksum,
		updated_at = EXCLUDED.updated_at,
		synced_at = EXCLUDED.synced_at;
	`

	_, err := r.pool.Exec(ctx, query,
		game.ID,
		game.Title,
		game.Slug,
		game.IconURL,
		game.Checksum,
		game.UpdatedAt,
		game.SyncedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert game id='%d': %w", game.ID, err)
	}

	return nil
}
