package users_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *UsersService) GetProfileByID(ctx context.Context, userID uuid.UUID) (domain.UserProfile, error) {
	profile, err := s.usersRepository.GetProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("get user profile from repository: %w", err)
	}

	return profile, nil
}
