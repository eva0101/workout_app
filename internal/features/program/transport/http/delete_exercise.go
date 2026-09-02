package program_transport_http

import (
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

func (h *ProgramHTTPHandler) DeleteExercise(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)

		return
	}

	trainingDayExerciseID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get trainingDayExerciseID path value", http.StatusBadRequest)

		return
	}

	if err := h.programService.DeleteExercise(
		r.Context(),
		userID,
		trainingDayExerciseID,
	); err != nil {
		core_http_errors.WriteError(rw, err)

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
