package auth_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (r *AuthRepository) MarkEmailVerified(ctx context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tag, err := r.pool.Exec(
		ctx,
		`UPDATE dps.users SET is_email_verified = TRUE WHERE id = $1;`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("mark email verified for user '%s': %w", userID, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%s': %w", userID, core_errors.ErrNotFound)
	}

	return nil
}
