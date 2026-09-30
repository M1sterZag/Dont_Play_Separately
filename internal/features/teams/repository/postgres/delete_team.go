package teams_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func (r *TeamsRepository) DeleteTeam(ctx context.Context, teamID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM dps.teams WHERE id=$1;`

	cmdTag, err := r.pool.Exec(ctx, query, teamID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("team with id='%s': %w", teamID, core_errors.ErrNotFound)
	}

	return nil
}
