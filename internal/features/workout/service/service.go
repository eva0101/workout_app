package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

type WorkoutService struct {
	workoutRepository WorkoutRepository
}

type WorkoutRepository interface {
	StartTraining(
		ctx context.Context,
		userID uuid.UUID,
		dayID int,
	) (core_domain.StartWorkout, error)
}

func NewWorkoutService(
	workoutRepository WorkoutRepository,
) *WorkoutService {
	return &WorkoutService{
		workoutRepository: workoutRepository,
	}
}
