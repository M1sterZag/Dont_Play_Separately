package auth_service

import (
	"fmt"
	"time"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type AccessClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	TokenType string    `json:"token_type"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

type JWTSigner struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTSigner(secret string, accessTTL, refreshTTL time.Duration) *JWTSigner {
	return &JWTSigner{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *JWTSigner) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := AccessClaims{
		TokenType: accessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

func (s *JWTSigner) GenerateRefreshToken(sessionID, userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := RefreshClaims{
		TokenType:  refreshTokenType,
		SessionID:  sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	return signed, nil
}

func (s *JWTSigner) ParseAccessToken(token string) (*AccessClaims, error) {
	var claims AccessClaims

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w: %w", err, core_errors.ErrUnauthenticated)
	}

	if claims.TokenType != accessTokenType {
		return nil, fmt.Errorf("token type is not access: %w", core_errors.ErrUnauthenticated)
	}

	return &claims, nil
}

func (s *JWTSigner) ParseRefreshToken(token string) (*RefreshClaims, error) {
	var claims RefreshClaims

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w: %w", err, core_errors.ErrUnauthenticated)
	}

	if claims.TokenType != refreshTokenType {
		return nil, fmt.Errorf("token type is not refresh: %w", core_errors.ErrUnauthenticated)
	}

	return &claims, nil
}
