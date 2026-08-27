package users_transport_http

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
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

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}
}
