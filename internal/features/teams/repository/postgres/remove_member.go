package teams_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	"github.com/google/uuid"
)

func (r *TeamsRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM dps.team_members
	WHERE team_id=$1 AND user_id=$2
		AND user_id <> (SELECT owner_id FROM dps.teams WHERE id=$1)
	RETURNING team_id;
	`

	var deletedTeamMemberID uuid.UUID
	err := r.pool.QueryRow(ctx, query, teamID, userID).Scan(&deletedTeamMemberID)
	if err != nil {
		if errors.Is(err, core_repository.ErrNoRows) {
			return fmt.Errorf("user with id='%s' is not a member of team with id='%s': %w", userID, teamID, core_errors.ErrNotFound)
		}
		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}
