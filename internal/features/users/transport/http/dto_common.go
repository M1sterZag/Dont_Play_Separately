package users_transport_http

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type UserProfileDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Version   int       `json:"version"`
	Nickname  string    `json:"nickname"`
	Bio       *string   `json:"bio,omitempty"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

func userProfileDTOFromDomain(profile domain.UserProfile) UserProfileDTOResponse {
	return UserProfileDTOResponse{
		ID:        profile.ID,
		Version:   profile.Version,
		Nickname:  profile.Nickname,
		Bio:       profile.Bio,
		AvatarURL: profile.AvatarURL,
		CreatedAt: profile.CreatedAt,
	}
}
