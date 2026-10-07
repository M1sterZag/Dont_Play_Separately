package auth_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (r *AuthRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dps.users
	SET hashed_password = $2, version = version + 1
	WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%s': %w", userID, core_errors.ErrNotFound)
	}

	return nil
}