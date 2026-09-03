package auth_repository

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID      uuid.UUID
	Version int

	Email          string
	HashedPassword string
	Nickname       string
	Bio            *string
	AvatarKey      string
	CreatedAt      time.Time
}

type RefreshSessionModel struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
