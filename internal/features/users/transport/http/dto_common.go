package users_transport_http

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
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

func userProfileDTOFromDomain(profile domain.UserProfile, storage core_storage.Storage) UserProfileDTOResponse {
	return UserProfileDTOResponse{
		ID:        profile.ID,
		Version:   profile.Version,
		Nickname:  profile.Nickname,
		Bio:       profile.Bio,
		AvatarURL: storage.PublicURL(profile.AvatarKey),
		CreatedAt: profile.CreatedAt,
	}
}
