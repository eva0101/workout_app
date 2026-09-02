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

	s.log.Info(
		"training approach was completed",
		zap.Int("training_day_exercises_id", workoutSet.TrainingDayExercisesID),
		zap.Int("set_number", workoutSet.SetNumber),
	)

	return workoutSet, nil
}
