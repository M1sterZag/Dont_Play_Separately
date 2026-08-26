package users_repository

import (
	"context"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type UsersRepository interface {
	GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	PatchUser(ctx context.Context, profilePatch domain.User) (domain.User, error)
}
