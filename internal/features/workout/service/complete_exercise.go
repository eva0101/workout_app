package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *WorkoutService) CompleteExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutExerciseID int,
) (core_domain.WorkoutExerciseCompletion, error) {
	completion, err := s.workoutRepository.CompleteExercise(
		ctx,
		userID,
		workoutExerciseID,
	)
	if err != nil {
		return core_domain.WorkoutExerciseCompletion{}, err
	}

	s.log.Info(
		"workout exercise completed",
		zap.Int("workout_exercise_id", completion.ID),
	)
	return completion, nil
}
