package teams_service

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (s *TeamsService) RemoveMember(ctx context.Context, teamID, userID, memberUserID uuid.UUID) error {
	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	if team.OwnerID != userID {
		return fmt.Errorf("check ownerID: %w", core_errors.ErrForbidden)
	}

	if team.OwnerID == memberUserID {
		return fmt.Errorf("cannot remove the owner of team '%s': %w", teamID, core_errors.ErrConflict)
	}

	if err := s.teamsRepository.RemoveMember(ctx, teamID, memberUserID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}

	return nil
}
