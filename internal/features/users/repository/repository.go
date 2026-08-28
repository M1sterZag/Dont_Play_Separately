package users_repository

import (
	"context"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type UsersRepository interface {
	GetProfileByUUID(ctx context.Context, userUUID uuid.UUID) (domain.UserProfile, error)
}
