package users_postgres_repository

import (
	"context"
	"errors"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
)

const patchUserQuery = `
	UPDATE dps.users
	SET version = $2, nickname = $3, bio = $4, avatar_url = $5
	WHERE id = $1
	RETURNING id, version, email, hashed_password, nickname, bio, avatar_url, created_at
`

func (r *UsersRepository) PatchUser(ctx context.Context, patch domain.User) (domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(
		ctx,
		patchUserQuery,
		patch.ID,
		patch.Version+1,
		patch.Nickname,
		patch.Bio,
		patch.AvatarUrl,
	).Scan(
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
			return domain.User{}, core_errors.ErrNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}
