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

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, hashed_password, nickname, bio, avatar_key, created_at
	FROM dps.users
	WHERE email = $1;
	`

	row := r.pool.QueryRow(ctx, query, email)

	var userModel auth_repository.UserModel
	err := row.Scan(
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
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with email='%s': %w", email, core_errors.ErrNotFound)
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
		userModel.AvatarKey,
		userModel.CreatedAt,
	)

	return userDomain, nil
}
