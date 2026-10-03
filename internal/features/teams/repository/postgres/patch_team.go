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

func (r *TeamsRepository) PatchTeam(ctx context.Context, teamID uuid.UUID, team domain.Team) (domain.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dps.teams
	SET title=$3, description=$4, is_rating_required=$5, desired_rating=$6,
		contact_link=$7, slots_total=$8, version=version+1
	WHERE id=$1 AND version=$2
	RETURNING id, version, owner_id, game_id, platform_id, title, description,
			is_rating_required, desired_rating, contact_link, slots_total,
			(SELECT COUNT(*) FROM dps.team_members tm WHERE tm.team_id = dps.teams.id) AS slots_taken,
			created_at, is_active;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		teamID, team.Version, team.Title, team.Description, team.IsRatingRequired, team.DesiredRating,
		team.ContactLink, team.SlotsTotal,
	)

	var teamModel teams_repository.TeamModel
	err := row.Scan(
		&teamModel.ID, &teamModel.Version, &teamModel.OwnerID, &teamModel.GameID, &teamModel.PlatformID,
		&teamModel.Title, &teamModel.Description, &teamModel.IsRatingRequired, &teamModel.DesiredRating,
		&teamModel.ContactLink, &teamModel.SlotsTotal, &teamModel.SlotsTaken, &teamModel.CreatedAt, &teamModel.IsActive,
	)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return domain.Team{}, fmt.Errorf("team with id='%s' concurently accessed: %w", teamID, core_errors.ErrConflict)
		}
		return domain.Team{}, fmt.Errorf("scan error: %w", err)
	}

	teamDomain := teams_repository.TeamDomainFromModel(teamModel)

	return teamDomain, nil
}
