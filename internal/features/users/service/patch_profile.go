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

	var favoritePlatformIDs *[]int
	if patch.FavoritePlatformIDs.Set {
		if patch.FavoritePlatformIDs.Value == nil {
			emptyPlatformIDs := []int{}
			favoritePlatformIDs = &emptyPlatformIDs
		} else {
			favoritePlatformIDs = patch.FavoritePlatformIDs.Value
		}
	}

	patchedProfile, err := s.usersRepository.PatchProfile(ctx, userID, profile, favoritePlatformIDs)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("patch profile: %w", err)
	}

	favoritePlatforms, err := s.usersRepository.GetFavoritePlatforms(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("get favorite platforms: %w", err)
	}
	patchedProfile.FavoritePlatforms = favoritePlatforms

	return patchedProfile, nil
}
