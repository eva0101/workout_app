package workout_transport_http

import (
	"context"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_http_server "workout_app/internal/core/transport/http/server"

	"github.com/google/uuid"
)

type WorkoutHTTPHandler struct {
	workoutService WorkoutService
}

type WorkoutService interface {
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

func NewWorkoutHTTPHandler(
	workoutService WorkoutService,
) *WorkoutHTTPHandler {
	return &WorkoutHTTPHandler{
		workoutService: workoutService,
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
			Path:    "/workouts/statistics",
			Handler: http.HandlerFunc(h.GetStatistic),
			Auth:    true,
		},
	}
}
