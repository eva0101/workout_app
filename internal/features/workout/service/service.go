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

	ExecuteSet(
		ctx context.Context,
		userID uuid.UUID,
		trainingDayExerciseID int,
		repsDone *int,
		weight *float64,
	) (core_domain.WorkoutSet, error)

	CompleteTraining(
		ctx context.Context,
		userID uuid.UUID,
		workoutID int,
		fatigueScore int,
	) (core_domain.Workout, error)

	GetStatistic(
		ctx context.Context,
		userID uuid.UUID,
	) ([]core_domain.Statistic, error)
}

func NewWorkoutService(
	workoutRepository WorkoutRepository,
) *WorkoutService {
	return &WorkoutService{
		workoutRepository: workoutRepository,
	}
}
