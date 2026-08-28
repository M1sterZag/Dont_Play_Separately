package auth_postgres_repository

import core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"

type AuthRepository struct {
	pool core_repository.Pool
}

func NewAuthRepository(pool core_repository.Pool) *AuthRepository {
	return &AuthRepository{
		pool: pool,
	}
}
