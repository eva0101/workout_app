package workout_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

func (h *WorkoutHTTPHandler) CompleteTraining(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)

		return
	}

	var req DTOFatigueScoreRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	workoutID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get workoutID path value", http.StatusBadRequest)

		return
	}

	workout, err := h.workoutService.CompleteTraining(
		r.Context(),
		userID,
		workoutID,
		req.FatigueScore,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)

		return
	}

	response := DTOCompleteWorkoutResponse{
		TrainingDayID: workout.TrainingDayID,
		Status:        workout.Status,
		FatigueScore:  workout.FatigueScore,
		TotalTime:     workout.TotalTime.String(),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(response); err != nil {

		return
	}
}
