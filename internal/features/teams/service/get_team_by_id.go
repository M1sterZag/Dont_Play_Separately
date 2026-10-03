package teams_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *TeamsService) GetTeamByID(ctx context.Context, teamID uuid.UUID) (domain.Team, error) {
	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("get team: %w", err)
	}

	return team, nil
}
