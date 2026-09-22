package games_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	games_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/games/repository"
)

func (r *GamesRepository) GetGameByID(ctx context.Context, gameID int64) (domain.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, title, slug, icon_url, checksum, updated_at, synced_at
	FROM dps.games
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, gameID)

	var gameModel games_repository.GameModel
	err := row.Scan(
		&gameModel.ID,
		&gameModel.Title,
		&gameModel.Slug,
		&gameModel.IconURL,
		&gameModel.Checksum,
		&gameModel.UpdatedAt,
		&gameModel.SyncedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.Game{}, fmt.Errorf("find game with id='%d': %w", gameID, core_errors.ErrNotFound)
		}

		return domain.Game{}, fmt.Errorf("scan error: %w", err)
	}

	gameDomain := games_repository.GameDomainFromModel(gameModel)

	return gameDomain, nil
}
