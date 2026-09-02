package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
)

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	now := time.Now()
	claims, err := s.jwtSigner.ParseRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, core_errors.ErrUnauthenticated
	}

	session, err := s.authRepository.FindSessionByID(ctx, claims.SessionID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return Tokens{}, core_errors.ErrUnauthenticated
		}
		return Tokens{}, fmt.Errorf("get session: %w", err)
	}

	if session.RevokedAt != nil {
		return Tokens{}, core_errors.ErrUnauthenticated
	}

	if session.ExpiresAt.Before(now) {
		return Tokens{}, core_errors.ErrUnauthenticated
	}

	if session.UserID.String() != claims.Subject {
		return Tokens{}, core_errors.ErrUnauthenticated
	}

	if err := s.authRepository.RevokeSession(ctx, session.ID); err != nil {
		return Tokens{}, fmt.Errorf("revoke session: %w", err)
	}

	newSession, newRefreshToken, err := s.newSession(session.UserID, now)
	if err != nil {
		return Tokens{}, fmt.Errorf("new session and refresh token: %w", err)
	}

	if err := s.authRepository.CreateSession(ctx, newSession); err != nil {
		return Tokens{}, fmt.Errorf("create session: %w", err)
	}

	newAccessToken, err := s.jwtSigner.GenerateAccessToken(session.UserID)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate new access token: %w", err)
	}

	return Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
