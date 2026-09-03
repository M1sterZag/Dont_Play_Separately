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
	Nickname  string
	Bio       *string
	AvatarKey string
	CreatedAt      time.Time
}

func NewUser(
	ID uuid.UUID,
	version int,
	email string,
	hashedPassword string,
	nickname string,
	bio *string,
	avatarKey string,
	createdAt time.Time) User {
	return User{
		ID:             ID,
		Version:        version,
		Email:          email,
		HashedPassword: hashedPassword,
		Nickname:       nickname,
		Bio:            bio,
		AvatarKey:      avatarKey,
		CreatedAt:      createdAt,
	}
}

type UserProfile struct {
	ID        uuid.UUID
	Version   int
	Nickname  string
	Bio            *string
	AvatarKey      string
	CreatedAt time.Time
}

func NewUserProfile(
	ID uuid.UUID,
	version int,
	nickname string,
	bio *string,
	avatarKey string,
	createdAt time.Time,
) UserProfile {
	return UserProfile{
		ID:        ID,
		Version:   version,
		Nickname:  nickname,
		Bio:       bio,
		AvatarKey: avatarKey,
		CreatedAt: createdAt,
	}
}
