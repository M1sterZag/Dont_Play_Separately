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

func (r *UsersRepository) PatchProfile(
	ctx context.Context,
	userID uuid.UUID,
	profile domain.UserProfile,
	favoritePlatformIDs *[]int,
) (domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
	UPDATE dps.users
	SET nickname=$3, bio=$4, avatar_key=$5, version=version+1
	WHERE id=$1 AND version=$2
	RETURNING id, version, nickname, bio, avatar_key, created_at;
	`

	row := tx.QueryRow(
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

	if favoritePlatformIDs != nil {
		if _, err := tx.Exec(
			ctx,
			`DELETE FROM dps.user_platforms WHERE user_id=$1;`,
			userID,
		); err != nil {
			return domain.UserProfile{}, fmt.Errorf("delete favorite platforms: %w", err)
		}

		for _, platformID := range *favoritePlatformIDs {
			if _, err := tx.Exec(
				ctx,
				`INSERT INTO dps.user_platforms (user_id, platform_id) VALUES ($1, $2);`,
				userID,
				platformID,
			); err != nil {
				if errors.Is(err, core_repository.ErrViolatesForeignKey) {
					return domain.UserProfile{}, fmt.Errorf(
						"favorite platform with id='%d' not found: %w",
						platformID,
						core_errors.ErrInvalidArgument,
					)
				}

				return domain.UserProfile{}, fmt.Errorf("insert favorite platform: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.UserProfile{}, fmt.Errorf("commit transaction: %w", err)
	}

	profileDomain := users_repository.UserProfileFromModel(userProfileModel)

	return profileDomain, nil
}
