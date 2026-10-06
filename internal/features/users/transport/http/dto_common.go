package users_transport_http

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
	core_http_types "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type UserProfileDTOResponse struct {
	ID                uuid.UUID                     `json:"id"`
	Version           int                           `json:"version"`
	Nickname          string                        `json:"nickname"`
	Bio               *string                       `json:"bio,omitempty"`
	AvatarURL         string                        `json:"avatar_url"`
	CreatedAt         time.Time                     `json:"created_at"`
	FavoritePlatforms []FavoritePlatformDTOResponse `json:"favorite_platforms"`
}

type FavoritePlatformDTOResponse struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Abbreviation string  `json:"abbreviation"`
	Slug         *string `json:"slug"`
	IconURL      *string `json:"icon_url"`
}

func userProfileDTOFromDomain(profile domain.UserProfile, storage core_storage.Storage) UserProfileDTOResponse {
	return UserProfileDTOResponse{
		ID:                profile.ID,
		Version:           profile.Version,
		Nickname:          profile.Nickname,
		Bio:               profile.Bio,
		AvatarURL:         storage.PublicURL(profile.AvatarKey),
		CreatedAt:         profile.CreatedAt,
		FavoritePlatforms: favoritePlatformsDTOFromDomain(profile.FavoritePlatforms),
	}
}

func favoritePlatformsDTOFromDomain(platforms []domain.Platform) []FavoritePlatformDTOResponse {
	favoritePlatforms := make([]FavoritePlatformDTOResponse, 0, len(platforms))
	for _, platform := range platforms {
		favoritePlatforms = append(favoritePlatforms, FavoritePlatformDTOResponse{
			ID:           platform.ID,
			Title:        platform.Title,
			Abbreviation: platform.Abbreviation,
			Slug:         platform.Slug,
			IconURL:      platform.IconURL,
		})
	}

	return favoritePlatforms
}

type PatchProfileRequest struct {
	Nickname            core_http_types.Nullable[string] `json:"nickname"`
	Bio                 core_http_types.Nullable[string] `json:"bio"`
	AvatarKey           core_http_types.Nullable[string] `json:"avatar_key"`
	FavoritePlatformIDs core_http_types.Nullable[[]int]  `json:"favorite_platform_ids"`
}

func userProfilePatchFromRequest(request PatchProfileRequest) domain.UserProfilePatch {
	return domain.NewUserProfilePatch(
		request.Nickname.ToDomain(),
		request.Bio.ToDomain(),
		request.AvatarKey.ToDomain(),
		request.FavoritePlatformIDs.ToDomain(),
	)
}
