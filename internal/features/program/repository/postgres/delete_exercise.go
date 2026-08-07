package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (r *ProgramRepository) DeleteExercise(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM workoutapp.training_day_exercises tde
	USING workoutapp.training_days td,
    	workoutapp.program p
	WHERE tde.id = $2
		AND td.id = tde.training_day_id
		AND p.id = td.program_id
		AND p.user_id = $1;
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		userID,
		trainingDayExerciseID,
	)
	if err != nil {

		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {

		return fmt.Errorf(
			"exercise with id='%d': %w",
			trainingDayExerciseID,
			core_errors.ErrExerciseNotFound,
		)
	}

	return nil
}
