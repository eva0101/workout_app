package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (r *ProgramRepository) DeleteProgram(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM workoutapp.program p
	WHERE p.id = $2
		AND p.user_id = $1;
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		userID,
		programID,
	)
	if err != nil {

		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {

		return fmt.Errorf(
			"program with id='%d': %w",
			programID,
			core_errors.ErrProgramNotFound,
		)
	}

	return nil
}
