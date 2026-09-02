package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
)

func (s *AuthService) Login(ctx context.Context, email, password string) (Tokens, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return Tokens{}, core_errors.ErrUnauthenticated
		}
		return Tokens{}, fmt.Errorf("get user by email: %w", err)
	}

	if !CheckPassword(user.HashedPassword, password) {
		return Tokens{}, core_errors.ErrUnauthenticated
	}

	now := time.Now()
	session, refreshToken, err := s.newSession(user.ID, now)
	if err != nil {
		return Tokens{}, fmt.Errorf("new session: %w", err)
	}
	if err := s.authRepository.CreateSession(ctx, session); err != nil {
		return Tokens{}, fmt.Errorf("save session: %w", err)
	}

	accessToken, err := s.jwtSigner.GenerateAccessToken(user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate access token: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
