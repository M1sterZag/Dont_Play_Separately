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

func (r *UsersRepository) GetUserByUUID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, hashed_password, nickname, bio, avatar_url, created_at
	FROM dps.users
	WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel users_repository.UserModel
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
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.User{}, fmt.Errorf("find user with uuid='%d': %w", id, core_errors.ErrNotFound)
		}
	}

	userDomain := users_repository.UserDomainFromModel(userModel)

	return userDomain, nil
}
