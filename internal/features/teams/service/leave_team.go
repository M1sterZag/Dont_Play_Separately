package teams_service

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (s *TeamsService) LeaveTeam(ctx context.Context, teamID, userID uuid.UUID) error {
	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	if team.OwnerID == userID {
		return fmt.Errorf("owner cannot leave the team '%s': %w", teamID, core_errors.ErrConflict)
	}

	if err := s.teamsRepository.LeaveTeam(ctx, teamID, userID); err != nil {
		return fmt.Errorf("delete user '%s' from team '%s': %w", userID, teamID, err)
	}

	return nil
}
