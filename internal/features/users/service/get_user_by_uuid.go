package users_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *UsersService) GetUserByUUID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.usersRepository.GetUserByUUID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
