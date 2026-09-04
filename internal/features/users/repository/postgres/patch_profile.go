package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	users_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/repository"
	"github.com/google/uuid"
)

func (r *UsersRepository) PatchProfile(ctx context.Context, userID uuid.UUID, profile domain.UserProfile) (domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dps.users
	SET nickname=$3, bio=$4, avatar_key=$5, version=version+1
	WHERE id=$1 AND version=$2
	RETURNING id, version, nickname, bio, avatar_key, created_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		userID,
		profile.Version,
		profile.Nickname,
		profile.Bio,
		profile.AvatarKey,
	)

	var userProfileModel users_repository.UserProfileModel
	if err := row.Scan(
		&userProfileModel.ID,
		&userProfileModel.Version,
		&userProfileModel.Nickname,
		&userProfileModel.Bio,
		&userProfileModel.AvatarKey,
		&userProfileModel.CreatedAt,
	); err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.UserProfile{}, fmt.Errorf("user with id='%s' concurently accessed: %w", userID, core_errors.ErrConflict)
		}
		return domain.UserProfile{}, fmt.Errorf("scan error: %w", err)
	}

	profileDomain := users_repository.UserProfileFromModel(userProfileModel)

	return profileDomain, nil
}
