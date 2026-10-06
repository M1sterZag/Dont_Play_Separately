package auth_service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
)

const verificationCacheKeyPrefix = "email_verification:"

type verificationCodeValue struct {
	CodeHash  string    `json:"code_hash"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}

func verificationCacheKey(email string) string {
	return verificationCacheKeyPrefix + strings.ToLower(strings.TrimSpace(email))
}

func generateVerificationCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	const digits = "0123456789"
	max := big.NewInt(int64(len(digits)))

	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("generate verification code digit: %w", err)
		}
		code[i] = digits[n.Int64()]
	}

	return string(code), nil
}

func hashVerificationCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) saveVerificationCode(ctx context.Context, email, code string, now time.Time) error {
	value := verificationCodeValue{
		CodeHash:  hashVerificationCode(code),
		Attempts:  0,
		CreatedAt: now,
	}

	if err := s.cache.Set(ctx, verificationCacheKey(email), value, s.config.VerificationCodeTTL); err != nil {
		return fmt.Errorf("save verification code for '%s': %w", email, err)
	}

	return nil
}

func (s *AuthService) sendVerificationEmail(ctx context.Context, email, code string) error {
	subject := "Подтверждение аккаунта"
	body := fmt.Sprintf(
		"Здравствуйте!\n\nВаш код подтверждения аккаунта: %s\n\nКод действителен в течение %s.",
		code,
		s.config.VerificationCodeTTL,
	)

	if err := s.mailer.Send(ctx, email, subject, body); err != nil {
		return fmt.Errorf("send verification email to '%s': %w", email, err)
	}

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) (Tokens, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return Tokens{}, core_errors.ErrInvalidArgument
		}
		return Tokens{}, fmt.Errorf("get user by email: %w", err)
	}

	if user.IsEmailVerified {
		return Tokens{}, core_errors.ErrConflict
	}

	key := verificationCacheKey(email)

	var value verificationCodeValue
	if err := s.cache.Get(ctx, key, &value); err != nil {
		if errors.Is(err, core_cache.ErrNotFound) {
			return Tokens{}, core_errors.ErrInvalidArgument
		}
		return Tokens{}, fmt.Errorf("get verification code for '%s': %w", email, err)
	}

	if value.Attempts >= s.config.MaxVerificationAttempts {
		_ = s.cache.Delete(ctx, key)
		return Tokens{}, core_errors.ErrInvalidArgument
	}

	if subtle.ConstantTimeCompare(
		[]byte(value.CodeHash),
		[]byte(hashVerificationCode(code)),
	) != 1 {
		value.Attempts++
		if err := s.cache.Set(ctx, key, value, s.config.VerificationCodeTTL); err != nil {
			return Tokens{}, fmt.Errorf("update verification attempts for '%s': %w", email, err)
		}

		return Tokens{}, core_errors.ErrInvalidArgument
	}

	if err := s.authRepository.MarkEmailVerified(ctx, user.ID); err != nil {
		return Tokens{}, fmt.Errorf("mark email verified: %w", err)
	}

	if err := s.cache.Delete(ctx, key); err != nil {
		return Tokens{}, fmt.Errorf("delete verification code for '%s': %w", email, err)
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

func (s *AuthService) ResendVerificationCode(ctx context.Context, email string) error {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_errors.ErrNotFound
		}
		return fmt.Errorf("get user by email: %w", err)
	}

	if user.IsEmailVerified {
		return core_errors.ErrConflict
	}

	now := time.Now()

	var value verificationCodeValue
	err = s.cache.Get(ctx, verificationCacheKey(email), &value)
	switch {
	case err == nil:
		if now.Sub(value.CreatedAt) < s.config.MinResendInterval {
			return core_errors.ErrConflict
		}
	case errors.Is(err, core_cache.ErrNotFound):
	default:
		return fmt.Errorf("get verification code for '%s': %w", email, err)
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