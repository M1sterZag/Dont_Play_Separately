package teams_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	"github.com/google/uuid"
)

func (r *TeamsRepository) JoinTeam(ctx context.Context, teamID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH eligible AS (
    SELECT t.id
    FROM dps.teams t
    WHERE t.id = $1
      AND t.is_active = TRUE
      AND t.slots_total > (SELECT COUNT(*) FROM dps.team_members tm WHERE tm.team_id = t.id)
	), inserted AS (
		INSERT INTO dps.team_members (team_id, user_id)
		SELECT id, $2 FROM eligible
		ON CONFLICT (team_id, user_id) DO NOTHING
		RETURNING team_id, user_id
	)
	SELECT team_id, user_id FROM inserted;
	`

	var insertedTeamID uuid.UUID
	if err := r.pool.QueryRow(ctx, query, teamID, userID).Scan(&insertedTeamID); err == nil {
		return nil
	} else if !errors.Is(err, core_repository.ErrNoRows) {
		return fmt.Errorf("join team: %w", err)
	}

	var exists bool
	if err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM dps.team_members WHERE team_id=$1 AND user_id=$2);`,
		teamID, userID,
	).Scan(&exists); err != nil {
		return fmt.Errorf("check membership: %w", err)
	}
	if exists {
		return fmt.Errorf("user with id='%s' already in team with id='%s': %w", userID, teamID, core_errors.ErrConflict)
	}

	team, err := r.GetTeamByID(ctx, teamID)
	if err != nil {
		return err
	}
	if !team.IsActive {
		return fmt.Errorf("team with id='%s' is not active: %w", teamID, core_errors.ErrConflict)
	}
	if team.SlotsTaken >= team.SlotsTotal {
		return fmt.Errorf("team with id='%s' is full: %w", teamID, core_errors.ErrConflict)
	}

	return fmt.Errorf("join team: unexpected: %w", core_errors.ErrConflict)
}
