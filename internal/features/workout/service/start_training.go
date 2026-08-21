package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *WorkoutService) StartTraining(
	ctx context.Context,
	userID uuid.UUID,
	dayID int,
) (core_domain.StartWorkout, error) {
	startWorkout, err := s.workoutRepository.StartTraining(
		ctx,
		userID,
		dayID,
	)
	if err != nil {

		return core_domain.StartWorkout{}, err
	}

	return startWorkout, nil
}
