package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *WorkoutService) ExecuteSet(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
	repsDone *int,
	weight *float64,
) (core_domain.WorkoutSet, error) {
	workoutSet, err := s.workoutRepository.ExecuteSet(
		ctx,
		userID,
		trainingDayExerciseID,
		repsDone,
		weight,
	)
	if err != nil {

		return core_domain.WorkoutSet{}, err
	}

	return workoutSet, nil
}
