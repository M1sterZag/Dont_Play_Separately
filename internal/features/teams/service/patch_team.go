package teams_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (s *TeamsService) PatchTeam(ctx context.Context, userID, teamID uuid.UUID, patch domain.TeamPatch) (domain.Team, error) {
	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("get team: %w", err)
	}

	if team.OwnerID != userID {
		return domain.Team{}, fmt.Errorf("user '%s' is not the owner of team '%s': %w", userID, teamID, core_errors.ErrForbidden)
	}

	if err := team.ApplyPatch(patch); err != nil {
		return domain.Team{}, fmt.Errorf("apply patch: %w", err)
	}

	patchedTeam, err := s.teamsRepository.PatchTeam(ctx, teamID, team)
	if err != nil {
		return domain.Team{}, fmt.Errorf("patch team: %w", err)
	}

	return patchedTeam, nil
}
