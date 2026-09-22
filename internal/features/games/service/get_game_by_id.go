package games_service

import (
	"context"
	"errors"
	"fmt"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	"go.uber.org/zap"
)

func (s *GamesService) GetGameByID(ctx context.Context, gameID int64) (domain.Game, error) {
	log := core_logger.FromContext(ctx)

	key := fmt.Sprintf("catalog:game:%d", gameID)

	var cached domain.Game
	if err := s.cache.Get(ctx, key, &cached); err == nil {
		return cached, nil
	} else if !errors.Is(err, core_cache.ErrNotFound) {
		log.Warn("cache get failed", zap.String("key", key), zap.Error(err))
	}

	v, err, _ := s.singleFlight.Do(key, func() (any, error) {
		game, err := s.gamesRepository.GetGameByID(ctx, gameID)
		if err == nil {
			if err := s.cache.Set(ctx, key, game, s.cacheTTL); err != nil {
				log.Warn("cache set failed", zap.String("key", key), zap.Error(err))
			}
			return game, nil
		}
		if !errors.Is(err, core_errors.ErrNotFound) {
			return domain.Game{}, err
		}

		games, err := s.provider.FetchGamesByIDs(ctx, []int64{gameID})
		if err != nil {
			return domain.Game{}, fmt.Errorf("fetch game %d from provider: %w", gameID, err)
		}
		if len(games) == 0 {
			return domain.Game{}, core_errors.ErrNotFound
		}

		game = markSynced(games[0])
		if err := s.gamesRepository.UpsertGame(ctx, game); err != nil {
			return domain.Game{}, fmt.Errorf("upsert game %d: %w", gameID, err)
		}
		if err := s.cache.Set(ctx, key, game, s.cacheTTL); err != nil {
			log.Warn("cache set failed", zap.String("key", key), zap.Error(err))
		}

		return game, nil
	})
	if err != nil {
		return domain.Game{}, err
	}

	game, ok := v.(domain.Game)
	if !ok {
		return domain.Game{}, fmt.Errorf("unexpected singleflight result type: %T", v)
	}
	return game, nil
}
