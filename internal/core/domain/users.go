package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Version        int
	Email          string
	HashedPassword string
	Nickname       string
	Bio            *string
	AvatarURL      string
	CreatedAt      time.Time
}

func NewUser(
	ID uuid.UUID,
	version int,
	email string,
	hashedPassword string,
	nickname string,
	bio *string,
	avatarURL string,
	createdAt time.Time) User {
	return User{
		ID:             ID,
		Version:        version,
		Email:          email,
		HashedPassword: hashedPassword,
		Nickname:       nickname,
		Bio:            bio,
		AvatarURL:      avatarURL,
		CreatedAt:      createdAt,
	}
}
