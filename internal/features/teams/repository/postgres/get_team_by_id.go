package teams_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	teams_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/teams/repository"
	"github.com/google/uuid"
)

func (r *TeamsRepository) GetTeamByID(ctx context.Context, teamID uuid.UUID) (domain.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT t.id, t.version, t.owner_id, t.game_id, t.platform_id, t.title, t.description,
    	t.is_rating_required, t.desired_rating, t.contact_link, t.slots_total,
    	(SELECT COUNT(*) FROM dps.team_members tm WHERE tm.team_id = t.id) AS slots_taken,
    	t.created_at, t.is_active
	FROM dps.teams t
	WHERE t.id = $1;
	`

	row := r.pool.QueryRow(ctx, query, teamID)

	var teamModel teams_repository.TeamModel
	err := row.Scan(
		&teamModel.ID, &teamModel.Version, &teamModel.OwnerID, &teamModel.GameID, &teamModel.PlatformID,
		&teamModel.Title, &teamModel.Description, &teamModel.IsRatingRequired, &teamModel.DesiredRating,
		&teamModel.ContactLink, &teamModel.SlotsTotal, &teamModel.SlotsTaken, &teamModel.CreatedAt, &teamModel.IsActive,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.Team{}, fmt.Errorf("find team with id='%s': %w", teamID, core_errors.ErrNotFound)
		}

		return domain.Team{}, fmt.Errorf("scan error: %w", err)
	}

	teamDomain := teams_repository.TeamDomainFromModel(teamModel)

	return teamDomain, nil
}
