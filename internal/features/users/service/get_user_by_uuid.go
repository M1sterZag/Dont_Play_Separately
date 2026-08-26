package users_service

import (
	"context"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *UsersService) GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {

}
