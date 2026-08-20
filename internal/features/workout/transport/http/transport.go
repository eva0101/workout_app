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
	}
}
