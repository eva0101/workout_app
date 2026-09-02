package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_logger "workout_app/internal/core/logger"

	"github.com/google/uuid"
)

type WorkoutService struct {
	workoutRepository WorkoutRepository
	log               *core_logger.Logger
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
	log *core_logger.Logger,
) *WorkoutService {
	return &WorkoutService{
		workoutRepository: workoutRepository,
		log:               log,
	}
}
