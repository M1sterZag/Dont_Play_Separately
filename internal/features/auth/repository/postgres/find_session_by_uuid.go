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

func (r *AuthRepository) FindSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (domain.RefreshSession, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, user_id, refresh_token_hash, expires_at, revoked_at, created_at
	FROM dps.refresh_sessions
	WHERE id = $1;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		sessionUUID,
	)

	var refreshSessionModel auth_repository.RefreshSessionModel
	err := row.Scan(
		&refreshSessionModel.ID,
		&refreshSessionModel.UserUUID,
		&refreshSessionModel.TokenHash,
		&refreshSessionModel.ExpiresAt,
		&refreshSessionModel.RevokedAt,
		&refreshSessionModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.RefreshSession{}, fmt.Errorf("session with id='%s': %w", sessionUUID, core_errors.ErrNotFound)
		}

		return domain.RefreshSession{}, fmt.Errorf("scan error: %w", err)
	}

	refreshSessionDomain := domain.NewRefreshSession(
		refreshSessionModel.ID,
		refreshSessionModel.UserUUID,
		refreshSessionModel.TokenHash,
		refreshSessionModel.ExpiresAt,
		refreshSessionModel.RevokedAt,
		refreshSessionModel.CreatedAt,
	)

	return refreshSessionDomain, nil
}
