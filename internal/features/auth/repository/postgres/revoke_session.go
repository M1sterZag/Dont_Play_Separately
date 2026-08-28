package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *AuthRepository) RevokeSession(ctx context.Context, sessionUUID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dps.refresh_sessions
	SET revoked_at = now()
	WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, sessionUUID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
