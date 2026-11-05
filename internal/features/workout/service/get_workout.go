package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *WorkoutService) GetActiveWorkout(
	ctx context.Context,
	userID uuid.UUID,
) (core_domain.WorkoutDetails, error) {
	return s.workoutRepository.GetActiveWorkout(ctx, userID)
}

func (s *WorkoutService) GetWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID int,
) (core_domain.WorkoutDetails, error) {
	return s.workoutRepository.GetWorkout(ctx, userID, workoutID)
}
