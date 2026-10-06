package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	auth_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/repository"
)

func (r *AuthRepository) CreateUser(ctx context.Context, user domain.User, favoritePlatformIDs []int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
	INSERT INTO dps.users
	(id, version, email, hashed_password, nickname, bio, avatar_key, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version, email, hashed_password, nickname, bio, avatar_key, created_at;
	`

	row := tx.QueryRow(
		ctx,
		query,
		user.ID,
		user.Version,
		user.Email,
		user.HashedPassword,
		user.Nickname,
		user.Bio,
		user.AvatarKey,
		user.CreatedAt,
	)

	var userModel auth_repository.UserModel
	err = row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Email,
		&userModel.HashedPassword,
		&userModel.Nickname,
		&userModel.Bio,
		&userModel.AvatarKey,
		&userModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrUniqueViolation) {
			return domain.User{}, fmt.Errorf("create user: %w", core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	for _, platformID := range favoritePlatformIDs {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO dps.user_platforms (user_id, platform_id) VALUES ($1, $2);`,
			user.ID,
			platformID,
		); err != nil {
			if errors.Is(err, core_repository.ErrViolatesForeignKey) {
				return domain.User{}, fmt.Errorf(
					"favorite platform with id='%d' not found: %w",
					platformID,
					core_errors.ErrInvalidArgument,
				)
			}

			return domain.User{}, fmt.Errorf("insert favorite platform: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit transaction: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Email,
		userModel.HashedPassword,
		userModel.Nickname,
		userModel.Bio,
		userModel.AvatarKey,
		userModel.CreatedAt,
	)

	return userDomain, nil
}
