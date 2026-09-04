package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *UsersService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	if err := s.usersRepository.DeleteProfile(ctx, userID); err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	return nil
}
