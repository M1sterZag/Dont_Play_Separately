package games_service

import (
	"context"
	"time"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"golang.org/x/sync/singleflight"
)

type GamesRepository interface {
	GetGameByID(ctx context.Context, gameID int64) (domain.Game, error)
	SearchGames(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Game, error)
	UpsertGame(ctx context.Context, game domain.Game) error
}

type GamesProvider interface {
	FetchGamesByIDs(ctx context.Context, gameIDs []int64) ([]domain.Game, error)
	SearchGames(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Game, error)
}

type GamesService struct {
	gamesRepository GamesRepository
	cache           core_cache.Cache
	provider        GamesProvider

	singleFlight singleflight.Group
	cacheTTL     time.Duration
}

func NewGameService(
	gamesRepository GamesRepository,
	cache core_cache.Cache,
	provider GamesProvider,
	cacheTTL time.Duration,
) *GamesService {
	return &GamesService{
		gamesRepository: gamesRepository,
		cache:           cache,
		provider:        provider,
		cacheTTL:        cacheTTL,
	}
}

func markSynced(game domain.Game) domain.Game {
	now := time.Now()
	game.SyncedAt = &now
	return game
}
