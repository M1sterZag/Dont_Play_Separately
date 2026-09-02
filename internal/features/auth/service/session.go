package auth_service

import (
	"fmt"
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *AuthService) newSession(userID uuid.UUID, now time.Time) (domain.RefreshSession, string, error) {
	sessionID := uuid.New()
	refreshToken, err := s.jwtSigner.GenerateRefreshToken(sessionID, userID)
	if err != nil {
		return domain.RefreshSession{}, "", fmt.Errorf("generate refresh token: %w", err)
	}

	return domain.NewRefreshSession(
		sessionID,
		userID,
		hashToken(refreshToken),
		now.Add(s.jwtSigner.refreshTTL),
		nil,
		now,
	), refreshToken, nil
}
