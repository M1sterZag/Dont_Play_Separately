package users_transport_http

import (
	"time"

	"github.com/google/uuid"
)

type UserDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Bio       *string   `json:"bio,omitempty"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}
