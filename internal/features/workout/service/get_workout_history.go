package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (s *WorkoutService) GetWorkoutHistory(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]core_domain.WorkoutHistoryItem, error) {
	if limit < 1 || limit > core_domain.WorkoutHistoryMaxLimit || offset < 0 {
		return nil, core_errors.ErrInvalidArgument
	}

	return s.workoutRepository.GetWorkoutHistory(ctx, userID, limit, offset)
}
