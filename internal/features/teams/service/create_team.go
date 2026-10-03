package teams_service

import (
	"context"
	"fmt"
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

func (s *TeamsService) CreateTeam(ctx context.Context, ownerID uuid.UUID, team domain.Team) (domain.Team, error) {
	now := time.Now()

	newTeam := domain.NewTeam(
		uuid.New(),
		1,
		ownerID,
		team.GameID,
		team.PlatformID,
		team.Title,
		team.Description,
		team.IsRatingRequired,
		team.DesiredRating,
		team.ContactLink,
		team.SlotsTotal,
		1,
		now,
		true,
	)

	if err := newTeam.Validate(); err != nil {
		return domain.Team{}, fmt.Errorf("validate team: %w", err)
	}

	team, err := s.teamsRepository.CreateTeam(ctx, newTeam)
	if err != nil {
		return domain.Team{}, fmt.Errorf("create team: %w", err)
	}

	return team, nil
}
