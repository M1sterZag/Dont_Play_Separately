package games_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	games_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/games/repository"
)

func (r *GamesRepository) SearchGames(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, slug, icon_url, checksum, updated_at, synced_at
	FROM dps.games
	WHERE title ILIKE '%' || $1 || '%'
	ORDER BY title ASC
	OFFSET $2
	LIMIT $3;
	`

	rows, err := r.pool.Query(ctx, query, searchQuery, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("query search games: %w", err)
	}
	defer rows.Close()

	var gameModels []games_repository.GameModel

	for rows.Next() {
		var gameModel games_repository.GameModel
		err := rows.Scan(
			&gameModel.ID,
			&gameModel.Title,
			&gameModel.Slug,
			&gameModel.IconURL,
			&gameModel.Checksum,
			&gameModel.UpdatedAt,
			&gameModel.SyncedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan games: %w", err)
		}

		gameModels = append(gameModels, gameModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	gameDomains := games_repository.GameDomainsFromModels(gameModels)

	return gameDomains, nil
}
