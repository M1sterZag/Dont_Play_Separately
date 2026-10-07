package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	auth_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/repository"
	"github.com/google/uuid"
)

func (r *AuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, hashed_password, nickname, bio, avatar_key, created_at, is_email_verified
	FROM dps.users
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, userID)

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
		&userModel.IsEmailVerified,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%s': %w", userID, core_errors.ErrNotFound)
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
		userModel.IsEmailVerified,
	)

	return userDomain, nil
}