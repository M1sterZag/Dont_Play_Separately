package domain

import (
	"fmt"
	"time"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Version         int
	Email           string
	HashedPassword  string
	Nickname        string
	Bio             *string
	AvatarKey       string
	CreatedAt       time.Time
	IsEmailVerified bool
}

func NewUser(
	ID uuid.UUID,
	version int,
	email string,
	hashedPassword string,
	nickname string,
	bio *string,
	avatarKey string,
	createdAt time.Time,
	isEmailVerified bool) User {
	return User{
		ID:              ID,
		Version:         version,
		Email:           email,
		HashedPassword:  hashedPassword,
		Nickname:        nickname,
		Bio:             bio,
		AvatarKey:       avatarKey,
		CreatedAt:       createdAt,
		IsEmailVerified: isEmailVerified,
	}
}

const MaxFavoritePlatforms = 3

type UserProfile struct {
	ID                uuid.UUID
	Version           int
	Nickname          string
	Bio               *string
	AvatarKey         string
	CreatedAt         time.Time
	FavoritePlatforms []Platform
}

func NewUserProfile(
	ID uuid.UUID,
	version int,
	nickname string,
	bio *string,
	avatarKey string,
	createdAt time.Time,
	favoritePlatforms []Platform,
) UserProfile {
	return UserProfile{
		ID:                ID,
		Version:           version,
		Nickname:          nickname,
		Bio:               bio,
		AvatarKey:         avatarKey,
		CreatedAt:         createdAt,
		FavoritePlatforms: favoritePlatforms,
	}
}

func (p *UserProfile) Validate() error {
	nicknameLen := len([]rune(p.Nickname))
	if nicknameLen < 1 || nicknameLen > 40 {
		return fmt.Errorf("invalid `Nickname` len: %d: %w", nicknameLen, core_errors.ErrInvalidArgument)
	}

	if !IsValidAvatarKey(p.AvatarKey) {
		return fmt.Errorf("invalid `AvatarKey` '%s': %w", p.AvatarKey, core_errors.ErrInvalidArgument)
	}

	return nil
}

type UserProfilePatch struct {
	Nickname            Nullable[string]
	Bio                 Nullable[string]
	AvatarKey           Nullable[string]
	FavoritePlatformIDs Nullable[[]int]
}

func NewUserProfilePatch(
	nickname Nullable[string],
	bio Nullable[string],
	avatarKey Nullable[string],
	favoritePlatformIDs Nullable[[]int],
) UserProfilePatch {
	return UserProfilePatch{
		Nickname:            nickname,
		Bio:                 bio,
		AvatarKey:           avatarKey,
		FavoritePlatformIDs: favoritePlatformIDs,
	}
}

func (p *UserProfilePatch) Validate() error {
	if p.Nickname.Set && p.Nickname.Value == nil {
		return fmt.Errorf("`Nickname` can`t be patched to `NULL`: %w", core_errors.ErrInvalidArgument)
	}

	if p.AvatarKey.Set && p.AvatarKey.Value == nil {
		return fmt.Errorf("`AvatarKey` can`t be patched to `NULL`: %w", core_errors.ErrInvalidArgument)
	}

	if p.FavoritePlatformIDs.Set && p.FavoritePlatformIDs.Value != nil {
		if err := ValidateFavoritePlatformIDs(*p.FavoritePlatformIDs.Value); err != nil {
			return fmt.Errorf("validate `FavoritePlatformIDs`: %w", err)
		}
	}

	return nil
}

func ValidateFavoritePlatformIDs(platformIDs []int) error {
	if len(platformIDs) > MaxFavoritePlatforms {
		return fmt.Errorf(
			"too many favorite platforms: %d, max allowed is %d: %w",
			len(platformIDs),
			MaxFavoritePlatforms,
			core_errors.ErrInvalidArgument,
		)
	}

	seen := make(map[int]struct{}, len(platformIDs))
	for _, id := range platformIDs {
		if id <= 0 {
			return fmt.Errorf("invalid platform id '%d': %w", id, core_errors.ErrInvalidArgument)
		}

		if _, ok := seen[id]; ok {
			return fmt.Errorf("duplicate platform id '%d': %w", id, core_errors.ErrInvalidArgument)
		}

		seen[id] = struct{}{}
	}

	return nil
}

func (p *UserProfile) ApplyPatch(patch UserProfilePatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate profile patch: %w", err)
	}

	tmp := *p
	if patch.Nickname.Set {
		tmp.Nickname = *patch.Nickname.Value
	}

	if patch.Bio.Set {
		tmp.Bio = patch.Bio.Value
	}

	if patch.AvatarKey.Set {
		tmp.AvatarKey = *patch.AvatarKey.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched profile: %w", err)
	}

	*p = tmp

	return nil
}
