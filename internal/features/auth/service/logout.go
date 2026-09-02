package auth_service

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
)

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.jwtSigner.ParseRefreshToken(refreshToken)
	if err != nil {
		return core_errors.ErrUnauthenticated
	}

	if err := s.authRepository.RevokeSession(ctx, claims.SessionUUID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
