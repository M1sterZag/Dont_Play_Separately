package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID uuid.UUID
	UserUUID uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}