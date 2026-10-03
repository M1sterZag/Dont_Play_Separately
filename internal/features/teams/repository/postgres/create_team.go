package teams_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	teams_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/teams/repository"
)

func (r *TeamsRepository) CreateTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH inserted_team AS (
    INSERT INTO dps.teams (
    	id, version, owner_id, game_id, platform_id, title, description,
        is_rating_required, desired_rating, contact_link, slots_total, created_at, is_active
    )
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    	RETURNING id, version, owner_id, game_id, platform_id, title, description,
              	is_rating_required, desired_rating, contact_link, slots_total, created_at, is_active
	), inserted_member AS (
    	INSERT INTO dps.team_members (team_id, user_id)
    	SELECT id, owner_id FROM inserted_team
	)
	SELECT id, version, owner_id, game_id, platform_id, title, description,
       	is_rating_required, desired_rating, contact_link, slots_total, created_at, is_active
	FROM inserted_team;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		team.ID, team.Version, team.OwnerID, team.GameID, team.PlatformID, team.Title, team.Description,
		team.IsRatingRequired, team.DesiredRating, team.ContactLink, team.SlotsTotal, team.CreatedAt, team.IsActive,
	)

	var teamModel teams_repository.TeamModel
	err := row.Scan(
		&teamModel.ID, &teamModel.Version, &teamModel.OwnerID, &teamModel.GameID, &teamModel.PlatformID, &teamModel.Title, &teamModel.Description,
		&teamModel.IsRatingRequired, &teamModel.DesiredRating, &teamModel.ContactLink, &teamModel.SlotsTotal, &teamModel.CreatedAt, &teamModel.IsActive,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrViolatesForeignKey) {
			return domain.Team{}, fmt.Errorf("`OwnerID (%d), GameID` (%d) or `PlatformID` (%d) is not exists: %w", team.OwnerID, team.GameID, team.PlatformID, core_errors.ErrInvalidArgument)
		}

		return domain.Team{}, fmt.Errorf("scan error: %w", err)
	}

	teamModel.SlotsTaken = 1

	teamDomain := teams_repository.TeamDomainFromModel(teamModel)

	return teamDomain, nil
}
