package users_repository

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type UserModel struct {
	ID      uuid.UUID
	Version int

	Email          string
	HashedPassword string
	Nickname       string
	Bio            *string
	AvatarKey      string
	CreatedAt      time.Time
}

func UserDomainFromModel(userModel UserModel) domain.User {
	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Email,
		userModel.HashedPassword,
		userModel.Nickname,
		userModel.Bio,
		userModel.AvatarKey,
		userModel.CreatedAt,
	)
}

type UserProfileModel struct {
	ID        uuid.UUID
	Version   int
	Nickname  string
	Bio       *string
	AvatarKey string
	CreatedAt time.Time
}

func UserProfileFromModel(profileModel UserProfileModel) domain.UserProfile {
	return domain.UserProfile{
		ID:        profileModel.ID,
		Version:   profileModel.Version,
		Nickname:  profileModel.Nickname,
		Bio:       profileModel.Bio,
		AvatarKey: profileModel.AvatarKey,
		CreatedAt: profileModel.CreatedAt,
	}
}

type PlatformModel struct {
	ID           int
	Title        string
	Abbreviation string
	Slug         *string
	IconURL      *string
	Checksum     *string
	UpdatedAt    *int
	SyncedAt     *time.Time
}

func PlatformDomainFromModel(platformModel PlatformModel) domain.Platform {
	return domain.NewPlatform(
		platformModel.ID,
		platformModel.Title,
		platformModel.Abbreviation,
		platformModel.Slug,
		platformModel.IconURL,
		platformModel.Checksum,
		platformModel.UpdatedAt,
		platformModel.SyncedAt,
	)
}

func PlatformDomainsFromModels(platformModels []PlatformModel) []domain.Platform {
	platformDomains := make([]domain.Platform, len(platformModels))
	for i, platformModel := range platformModels {
		platformDomains[i] = PlatformDomainFromModel(platformModel)
	}

	return platformDomains
}
