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

func (r *AuthRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO dps.users
	(id, version, email, hashed_password, nickname, bio, avatar_url, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version, email, hashed_password, nickname, bio, avatar_url, created_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Version,
		user.Email,
		user.HashedPassword,
		user.Nickname,
		user.Bio,
		user.AvatarURL,
		user.CreatedAt,
	)

	var userModel auth_repository.UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Email,
		&userModel.HashedPassword,
		&userModel.Nickname,
		&userModel.Bio,
		&userModel.AvatarURL,
		&userModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrUniqueViolation) {
			return domain.User{}, fmt.Errorf("create user: %w", core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Email,
		userModel.HashedPassword,
		userModel.Nickname,
		userModel.Bio,
		userModel.AvatarURL,
		userModel.CreatedAt,
	)

	return userDomain, nil
}
