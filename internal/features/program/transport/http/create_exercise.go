package program_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *ProgramHTTPHandler) CreateExercise(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	trainingDayID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get traingDayID path value", http.StatusBadRequest)
		return
	}

	var req CreateExerciseDTORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	exercise, err := h.programService.CreateExercise(
		r.Context(),
		userID,
		trainingDayID,
		req.ExerciseID,
		req.Sets,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(exercise); err != nil {
		h.log.Error(
			"failed to encode exercise",
			zap.Error(err),
		)

		return
	}
}
