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

func NewRefreshSession(
	id uuid.UUID,
	userUUID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
) RefreshSession {
	return RefreshSession{
		ID: id,
		UserUUID: userUUID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		RevokedAt: revokedAt,
		CreatedAt: createdAt,
	}
}