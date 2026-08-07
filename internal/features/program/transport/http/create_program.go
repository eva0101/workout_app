package program_transport_http

import (
	"encoding/json"
	"net/http"
	"strings"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"github.com/google/uuid"
)

func (h *ProgramHTTPHandler) CreateProgram(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateProgramDTORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(rw, "the program name cannot be empty", http.StatusBadRequest)
		return
	}

	program, err := h.programService.CreateProgram(
		r.Context(),
		userID,
		req.Name,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(program); err != nil {
		return
	}
}
