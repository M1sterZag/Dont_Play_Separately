package users_transport_http

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
	core_http_types "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/types"
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

type PatchProfileRequest struct {
	Nickname  core_http_types.Nullable[string] `json:"nickname"`
	Bio       core_http_types.Nullable[string] `json:"bio"`
	AvatarKey core_http_types.Nullable[string] `json:"avatar_key"`
}

func userProfilePatchFromRequest(request PatchProfileRequest) domain.UserProfilePatch {
	return domain.NewUserProfilePatch(
		request.Nickname.ToDomain(),
		request.Bio.ToDomain(),
		request.AvatarKey.ToDomain(),
	)
}
