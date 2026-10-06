package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *AuthService) Register(ctx context.Context, email, password, nickname string, favoritePlatformIDs []int) (Tokens, error) {
	now := time.Now()

	if err := domain.ValidateFavoritePlatformIDs(favoritePlatformIDs); err != nil {
		return Tokens{}, fmt.Errorf("validate favorite platform ids: %w", err)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return Tokens{}, fmt.Errorf("hash password: %w", err)
	}

	userID := uuid.New()
	user := domain.NewUser(
		userID,
		1,
		email,
		hashedPassword,
		nickname,
		nil,
		domain.SelectDefaultAvatar(userID),
		now,
	)

	if _, err := s.authRepository.CreateUser(ctx, user, favoritePlatformIDs); err != nil {
		return Tokens{}, fmt.Errorf("create user: %w", err)
	}

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
