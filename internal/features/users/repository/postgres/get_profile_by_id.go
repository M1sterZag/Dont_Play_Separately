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

func (r *UsersRepository) GetProfileByID(ctx context.Context, userID uuid.UUID) (domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, nickname, bio, avatar_key, created_at
	FROM dps.users
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, userID)

	var userProfileModel users_repository.UserProfileModel
	err := row.Scan(
		&userProfileModel.ID,
		&userProfileModel.Version,
		&userProfileModel.Nickname,
		&userProfileModel.Bio,
		&userProfileModel.AvatarKey,
		&userProfileModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.UserProfile{}, fmt.Errorf("find user with id='%s': %w", userID, core_errors.ErrNotFound)
		}

		return domain.UserProfile{}, fmt.Errorf("scan error: %w", err)
	}

	userProfileDomain := users_repository.UserProfileFromModel(userProfileModel)

	return userProfileDomain, nil
}
