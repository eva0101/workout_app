package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (s *WorkoutService) CompleteTraining(
	ctx context.Context,
	userID uuid.UUID,
	workoutID int,
	fatigueScore int,
) (core_domain.Workout, error) {
	if fatigueScore < 0 || fatigueScore > 10 {

		return core_domain.Workout{}, core_errors.ErrFatigueScoreNotCorrect
	}

	workout, err := s.workoutRepository.CompleteTraining(
		ctx,
		userID,
		workoutID,
		fatigueScore,
	)
	if err != nil {

		return core_domain.Workout{}, err
	}

	return workout, nil
}
