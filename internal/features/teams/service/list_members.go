package teams_service

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *TeamsService) ListMembers(ctx context.Context, teamID uuid.UUID) ([]domain.TeamMember, error) {
	teamMembers, err := s.teamsRepository.ListMembers(ctx, teamID)
	if err != nil {
		return []domain.TeamMember{}, fmt.Errorf("get team members: %w", err)
	}

	return teamMembers, nil
}
