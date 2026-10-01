package teams_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

func (s *TeamsService) ListTeams(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, error) {
	if err := filter.Validate(); err != nil {
		return []domain.Team{}, fmt.Errorf("validate filter: %w", err)
	}
	
	teams, err := s.teamsRepository.ListTeams(ctx, filter)
	if err != nil {
		return []domain.Team{}, fmt.Errorf("list teams: %w", err)
	}

	return teams, nil
}
