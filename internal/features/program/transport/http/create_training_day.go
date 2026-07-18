package program_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

func (h *ProgramHTTPHandler) CreateTrainingDay(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	programID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get programID path value", http.StatusBadRequest)
		return
	}

	trainingDay, err := h.programService.CreateTrainingDay(
		r.Context(),
		userID,
		programID,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(trainingDay); err != nil {
		return
	}
}
