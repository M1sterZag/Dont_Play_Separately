package games_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
)

func (s *GamesService) SearchGames(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Game, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	games, err := s.gamesRepository.SearchGames(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(games) > 0 {
		return games, nil
	}

	igdbGames, err := s.provider.SearchGames(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search games in provider: %w", err)
	}

	gameDomains := make([]domain.Game, 0, len(igdbGames))
	for _, game := range igdbGames {
		gameDomain := markSynced(game)
		if err := s.gamesRepository.UpsertGame(ctx, gameDomain); err != nil {
			return nil, fmt.Errorf("upsert game %d: %w", game.ID, err)
		}
		gameDomains = append(gameDomains, gameDomain)
	}
	return gameDomains, nil
}
