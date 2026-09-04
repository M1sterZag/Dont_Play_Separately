package users_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *UsersService) PatchProfile(ctx context.Context, userID uuid.UUID, patch domain.UserProfilePatch) (domain.UserProfile, error) {
	profile, err := s.usersRepository.GetProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("get profile: %w", err)
	}

	if err := profile.ApplyPatch(patch); err != nil {
		return domain.UserProfile{}, fmt.Errorf("apply patch: %w", err)
	}

	patchedProfile, err := s.usersRepository.PatchProfile(ctx, userID, profile)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("patch profile: %w", err)
	}

	return patchedProfile, nil
}
