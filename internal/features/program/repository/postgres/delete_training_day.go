package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *ProgramRepository) DeleteTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM workoutapp.training_days td
	USING workoutapp.program p
	WHERE td.id = $2
		AND p.id = td.program_id
		AND p.user_id = $1;
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		userID,
		trainingDayID,
	)
	if err != nil {
		r.log.Error(
			"failed to delete training day",
			zap.Error(err),
		)

		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {

		return fmt.Errorf(
			"training day with id='%d': %w",
			trainingDayID,
			core_errors.ErrTrainingDayNotFound,
		)
	}

	return nil
}
