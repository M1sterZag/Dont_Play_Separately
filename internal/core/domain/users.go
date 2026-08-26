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
	AvatarUrl      string
	CreatedAt      time.Time
}
