package auth_service

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.authRepository.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	if !CheckPassword(user.HashedPassword, oldPassword) {
		return core_errors.ErrForbidden
	}

	if oldPassword == newPassword {
		return fmt.Errorf("new password must differ from the current one: %w", core_errors.ErrInvalidArgument)
	}

	hashedNewPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	if err := s.authRepository.UpdatePassword(ctx, user.ID, hashedNewPassword); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	if err := s.authRepository.RevokeAllUserSessions(ctx, user.ID); err != nil {
		return fmt.Errorf("revoke all user sessions: %w", err)
	}

	return nil
}