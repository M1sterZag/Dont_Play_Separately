package teams_service

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (s *TeamsService) DeleteTeam(ctx context.Context, userID, teamID uuid.UUID) error {
	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	if team.OwnerID != userID {
		return fmt.Errorf("user '%s' is not the owner of team '%s': %w", userID, teamID, core_errors.ErrForbidden)
	}

	if err := s.teamsRepository.DeleteTeam(ctx, teamID); err != nil {
		return fmt.Errorf("delete team: %w", err)
	}

	return nil
}
