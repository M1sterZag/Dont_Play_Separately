package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *AuthRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dps.refresh_sessions
	SET revoked_at = now()
	WHERE user_id = $1 AND revoked_at IS NULL;
	`

	if _, err := r.pool.Exec(ctx, query, userID); err != nil {
		return fmt.Errorf("revoke all user sessions: %w", err)
	}

	return nil
}