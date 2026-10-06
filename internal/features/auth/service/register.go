package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *AuthService) Register(ctx context.Context, email, password, nickname string, favoritePlatformIDs []int) error {
	now := time.Now()

	if err := domain.ValidateFavoritePlatformIDs(favoritePlatformIDs); err != nil {
		return fmt.Errorf("validate favorite platform ids: %w", err)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
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
		false,
	)

	if _, err := s.authRepository.CreateUser(ctx, user, favoritePlatformIDs); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	code, err := generateVerificationCode(s.config.VerificationCodeLength)
	if err != nil {
		return fmt.Errorf("generate verification code: %w", err)
	}

	if err := s.saveVerificationCode(ctx, email, code, now); err != nil {
		return fmt.Errorf("save verification code: %w", err)
	}

	if err := s.sendVerificationEmail(ctx, email, code); err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}

	return nil
}
