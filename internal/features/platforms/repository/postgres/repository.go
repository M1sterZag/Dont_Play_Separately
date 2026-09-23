package platforms_postgres_repository

import core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"

type PlatformsRepository struct {
	pool core_repository.Pool
}

func NewPlatformsRepository(pool core_repository.Pool) *PlatformsRepository {
	return &PlatformsRepository{
		pool: pool,
	}
}