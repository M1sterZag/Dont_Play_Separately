package users_repository

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID      uuid.UUID
	Version string

	Email          string
	HashedPassword string
	Nickname       string
	Bio            *string
	AvatarURL      string
	CreatedAt      time.Time
}
