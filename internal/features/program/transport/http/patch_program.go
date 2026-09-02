package program_transport_http

import (
	"encoding/json"
	"net/http"
	"strings"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *ProgramHTTPHandler) PatchProgram(rw http.ResponseWriter, r *http.Request) {
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

	newNameProgram, err := h.programService.PatchProgram(
		r.Context(),
		userID,
		programID,
		req.Name,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(newNameProgram); err != nil {
		h.log.Error(
			"failed to encode new name program",
			zap.Error(err),
		)

		return
	}
}
