package users_postgres_repository

import (
	"context"
	"errors"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	"github.com/google/uuid"
)

const getUserByUUIDQuery = `
	SELECT id, version, email, hashed_password, nickname, bio, avatar_url, created_at
	FROM dps.users
	WHERE id = $1
`

func (r *UsersRepository) GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, getUserByUUIDQuery, id).Scan(
		&user.ID,
		&user.Version,
		&user.Email,
		&user.HashedPassword,
		&user.Nickname,
		&user.Bio,
		&user.AvatarUrl,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return nil, core_errors.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}
