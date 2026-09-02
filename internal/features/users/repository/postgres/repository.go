package users_postgres_repository

import (
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
)

type UsersRepository struct {
	pool core_repository.Pool
}

func NewUsersRepository(pool core_repository.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
