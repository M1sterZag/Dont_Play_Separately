package domain

import (
	"fmt"
	"time"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Version        int
	Email          string
	HashedPassword string
	Nickname       string
	Bio            *string
	AvatarKey      string
	CreatedAt      time.Time
}

func NewUser(
	ID uuid.UUID,
	version int,
	email string,
	hashedPassword string,
	nickname string,
	bio *string,
	avatarKey string,
	createdAt time.Time) User {
	return User{
		ID:             ID,
		Version:        version,
		Email:          email,
		HashedPassword: hashedPassword,
		Nickname:       nickname,
		Bio:            bio,
		AvatarKey:      avatarKey,
		CreatedAt:      createdAt,
	}
}

type UserProfile struct {
	ID        uuid.UUID
	Version   int
	Nickname  string
	Bio       *string
	AvatarKey string
	CreatedAt time.Time
}

func NewUserProfile(
	ID uuid.UUID,
	version int,
	nickname string,
	bio *string,
	avatarKey string,
	createdAt time.Time,
) UserProfile {
	return UserProfile{
		ID:        ID,
		Version:   version,
		Nickname:  nickname,
		Bio:       bio,
		AvatarKey: avatarKey,
		CreatedAt: createdAt,
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
	Nickname  Nullable[string]
	Bio       Nullable[string]
	AvatarKey Nullable[string]
}

func NewUserProfilePatch(nickname Nullable[string], bio Nullable[string], avatarKey Nullable[string]) UserProfilePatch {
	return UserProfilePatch{
		Nickname:  nickname,
		Bio:       bio,
		AvatarKey: avatarKey,
	}
}

func (p *UserProfilePatch) Validate() error {
	if p.Nickname.Set && p.Nickname.Value == nil {
		return fmt.Errorf("`Nickname` can`t be patched to `NULL`: %w", core_errors.ErrInvalidArgument)
	}

	if p.AvatarKey.Set && p.AvatarKey.Value == nil {
		return fmt.Errorf("`AvatarKey` can`t be patched to `NULL`: %w", core_errors.ErrInvalidArgument)
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
