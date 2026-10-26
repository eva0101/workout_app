package schedule_repository_postgres

import (
	"context"
	"fmt"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *ScheduleRepository) DeleteSchedule(
	ctx context.Context,
	userID uuid.UUID,
	scheduleID int,
) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM workoutapp.scheduled_workouts
		WHERE id = $1 AND user_id = $2;
	`, scheduleID, userID)
	if err != nil {
		r.log.Error("failed to delete scheduled workout", zap.Error(err))
		return fmt.Errorf("delete scheduled workout: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrScheduledWorkoutNotFound
	}
	return nil
}
