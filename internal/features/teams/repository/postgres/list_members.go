package teams_postgres_repository

import (
	"context"
	"fmt"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	teams_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/teams/repository"
	"github.com/google/uuid"
)

func (r *TeamsRepository) ListMembers(ctx context.Context, teamID uuid.UUID) ([]domain.TeamMember, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT tm.team_id, tm.user_id, u.nickname
	FROM dps.team_members tm
	JOIN dps.users u ON u.id = tm.user_id
	WHERE tm.team_id = $1
	ORDER BY u.nickname ASC, tm.user_id ASC;
	`

	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("query members: %w", err)
	}
	defer rows.Close()

	var teamMemberModels []teams_repository.TeamMemberModel
	for rows.Next() {
		var model teams_repository.TeamMemberModel
		if err := rows.Scan(&model.TeamID, &model.UserID, &model.Nickname); err != nil {
			return nil, fmt.Errorf("scan members: %w", err)
		}
		teamMemberModels = append(teamMemberModels, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	teamMemberDomains := teams_repository.TeamMemberDomainsFromModels(teamMemberModels)

	return teamMemberDomains, nil
}
