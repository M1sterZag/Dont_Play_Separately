package teams_postgres_repository

import core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"

type TeamsRepository struct {
	pool core_repository.Pool
}

func NewTeamsRepository(pool core_repository.Pool) *TeamsRepository {
	return &TeamsRepository{
		pool: pool,
	}
}
