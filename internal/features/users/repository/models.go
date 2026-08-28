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
	AvatarURL      string
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
		userModel.AvatarURL,
		userModel.CreatedAt,
	)
}

type UserProfileModel struct {
	ID        uuid.UUID
	Version   int
	Nickname  string
	Bio       *string
	AvatarURL string
	CreatedAt time.Time
}

func UserProfileFromModel(profileModel UserProfileModel) domain.UserProfile {
	return domain.UserProfile{
		ID:        profileModel.ID,
		Version:   profileModel.Version,
		Nickname:  profileModel.Nickname,
		Bio:       profileModel.Bio,
		AvatarURL: profileModel.AvatarURL,
		CreatedAt: profileModel.CreatedAt,
	}
}
