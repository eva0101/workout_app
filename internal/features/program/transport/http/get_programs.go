package program_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"github.com/google/uuid"
)

func (h *ProgramHTTPHandler) GetPrograms(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	programs, err := h.programService.GetPrograms(
		r.Context(),
		userID,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	response := ToProgramResponse(programs)

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		return
	}
}
