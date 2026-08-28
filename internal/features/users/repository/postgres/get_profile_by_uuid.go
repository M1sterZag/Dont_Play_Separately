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

func (r *UsersRepository) GetProfileByUUID(ctx context.Context, userUUID uuid.UUID) (domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	// Учётные данные (email, hashed_password) не выбираем: они не должны покидать слой данных.
	query := `
	SELECT id, version, nickname, bio, avatar_url, created_at
	FROM dps.users
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, userUUID)

	var profileModel users_repository.UserProfileModel
	err := row.Scan(
		&profileModel.ID,
		&profileModel.Version,
		&profileModel.Nickname,
		&profileModel.Bio,
		&profileModel.AvatarURL,
		&profileModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.UserProfile{}, fmt.Errorf("find user with uuid='%s': %w", userUUID, core_errors.ErrNotFound)
		}

		return domain.UserProfile{}, fmt.Errorf("scan user with uuid='%s': %w", userUUID, err)
	}

	return users_repository.UserProfileFromModel(profileModel), nil
}
