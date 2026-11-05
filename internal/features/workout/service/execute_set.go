package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *WorkoutService) ExecuteSet(
	ctx context.Context,
	userID uuid.UUID,
	workoutExerciseID int,
	repsDone *int,
	weight *float64,
) (core_domain.WorkoutSet, error) {
	workoutSet, err := s.workoutRepository.ExecuteSet(
		ctx,
		userID,
		workoutExerciseID,
		repsDone,
		weight,
	)
	if err != nil {

		return core_domain.WorkoutSet{}, err
	}

	s.log.Info(
		"training approach was completed",
		zap.Int("workout_exercise_id", workoutSet.WorkoutExerciseID),
		zap.Int("set_number", workoutSet.SetNumber),
	)

	return workoutSet, nil
}
