package teams_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	teams_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/teams/repository"
)

func (r *TeamsRepository) ListTeams(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT t.id, t.version, t.owner_id, t.game_id, t.platform_id, t.title, t.description,
	      t.is_rating_required, t.desired_rating, t.contact_link, t.slots_total,
	   (SELECT COUNT(*) FROM dps.team_members tm WHERE tm.team_id = t.id) AS slots_taken, t.created_at, t.is_active
	FROM dps.teams t
	JOIN dps.games g ON g.id = t.game_id
	JOIN dps.platforms p ON p.id = t.platform_id
	WHERE t.is_active = TRUE
		AND (
			$1 = ''
			OR t.title ILIKE '%' || $1 || '%'
			OR g.title ILIKE '%' || $1 || '%'
			OR p.title ILIKE '%' || $1 || '%'
			OR p.abbreviation ILIKE '%' || $1 || '%'
		)
		AND ($2::INTEGER IS NULL OR t.game_id = $2)
		AND ($3::INTEGER IS NULL OR t.platform_id = $3)
	ORDER BY t.created_at DESC, t.id ASC
	OFFSET COALESCE($4, 0)
	LIMIT COALESCE($5, 20);
	`

	rows, err := r.pool.Query(ctx, query, filter.Search, filter.GameID, filter.PlatformID, filter.Offset, filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("query search teams: %w", err)
	}
	defer rows.Close()

	var teamModels []teams_repository.TeamModel
	for rows.Next() {
		var teamModel teams_repository.TeamModel
		err := rows.Scan(
			&teamModel.ID, &teamModel.Version, &teamModel.OwnerID, &teamModel.GameID, &teamModel.PlatformID, &teamModel.Title, &teamModel.Description,
			&teamModel.IsRatingRequired, &teamModel.DesiredRating, &teamModel.ContactLink, &teamModel.SlotsTotal, &teamModel.SlotsTaken, &teamModel.CreatedAt, &teamModel.IsActive,
		)

		if err != nil {
			return nil, fmt.Errorf("scan teams: %w", err)
		}

		teamModels = append(teamModels, teamModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	teamDomains := teams_repository.TeamDomainsFromModels(teamModels)

	return teamDomains, nil
}
