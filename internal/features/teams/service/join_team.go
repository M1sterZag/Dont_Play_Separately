package teams_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TeamsService) JoinTeam(ctx context.Context, teamID, userID uuid.UUID) error {
	if err := s.teamsRepository.JoinTeam(ctx, teamID, userID); err != nil {
		return fmt.Errorf("joining user '%s' in team '%s': %w", userID, teamID, err)
	}

	return nil
}
