package users_service

import (
	users_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/repository"
)

type UsersService struct {
	usersRepository users_repository.UsersRepository
}

func NewUsersService(usersRepository users_repository.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
