package workout_transport_http

import (
	"context"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_logger "workout_app/internal/core/logger"
	core_http_server "workout_app/internal/core/transport/http/server"

	"github.com/google/uuid"
)

type WorkoutHTTPHandler struct {
	workoutService WorkoutService
	log            *core_logger.Logger
}

type WorkoutService interface {
	StartTraining(
		ctx context.Context,
		userID uuid.UUID,
		dayID int,
	) (core_domain.StartWorkout, error)

	StartScheduledTraining(
		ctx context.Context,
		userID uuid.UUID,
		scheduleID int,
	) (core_domain.StartWorkout, error)

	CompleteExercise(
		ctx context.Context,
		userID uuid.UUID,
		workoutExerciseID int,
	) (core_domain.WorkoutExerciseCompletion, error)

	ExecuteSet(
		ctx context.Context,
		userID uuid.UUID,
		workoutExerciseID int,
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

	GetStatisticsDashboard(
		ctx context.Context,
		userID uuid.UUID,
		period string,
	) (core_domain.WorkoutStatisticsDashboard, error)

	GetWorkoutHistory(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
		offset int,
	) ([]core_domain.WorkoutHistoryItem, error)

	GetActiveWorkout(
		ctx context.Context,
		userID uuid.UUID,
	) (core_domain.WorkoutDetails, error)

	GetWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID int,
	) (core_domain.WorkoutDetails, error)
}

func NewWorkoutHTTPHandler(
	workoutService WorkoutService,
	log *core_logger.Logger,
) *WorkoutHTTPHandler {
	return &WorkoutHTTPHandler{
		workoutService: workoutService,
		log:            log,
	}
}

func (h *WorkoutHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/days/{id}/workouts",
			Handler: http.HandlerFunc(h.StartTraining),
			Auth:    true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/workouts/schedule/{id}/start",
			Handler: http.HandlerFunc(h.StartScheduledTraining),
			Auth:    true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/workouts/exercises/{id}/complete",
			Handler: http.HandlerFunc(h.CompleteExercise),
			Auth:    true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/workouts/exercises/{id}/sets",
			Handler: http.HandlerFunc(h.ExecuteSet),
			Auth:    true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/workouts/{id}/complete",
			Handler: http.HandlerFunc(h.CompleteTraining),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/history",
			Handler: http.HandlerFunc(h.GetWorkoutHistory),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/active",
			Handler: http.HandlerFunc(h.GetActiveWorkout),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/{id}",
			Handler: http.HandlerFunc(h.GetWorkout),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/statistics",
			Handler: http.HandlerFunc(h.GetStatistic),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/statistics/dashboard",
			Handler: http.HandlerFunc(h.GetStatisticsDashboard),
			Auth:    true,
		},
	}
}
