package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

func (r *AuthRepository) CreateSession(ctx context.Context, session domain.RefreshSession) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO dps.refresh_sessions
	(id, user_id, token_hash, expires_at, revoked_at, created_at)
	VALUES
	($1, $2, $3, $4, $5, $6);
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		session.ID,
		session.UserUUID,
		session.TokenHash,
		session.ExpiresAt,
		session.RevokedAt,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	return nil
}
