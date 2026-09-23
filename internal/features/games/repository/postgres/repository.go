package games_postgres_repository

import core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"

type GamesRepository struct {
	pool core_repository.Pool
}

func NewGamesRepository(pool core_repository.Pool) *GamesRepository {
	return &GamesRepository{
		pool: pool,
	}
}
