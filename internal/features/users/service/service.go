package users_service

import (
	"context"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type UsersRepository interface {
	GetProfileByID(ctx context.Context, userID uuid.UUID) (domain.UserProfile, error)
	GetFavoritePlatforms(ctx context.Context, userID uuid.UUID) ([]domain.Platform, error)
	PatchProfile(ctx context.Context, userID uuid.UUID, profile domain.UserProfile, favoritePlatformIDs *[]int) (domain.UserProfile, error)
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type UsersService struct {
	usersRepository UsersRepository
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
