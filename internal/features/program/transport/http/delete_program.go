package program_transport_http

import (
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

// DeleteProgram godoc
// @Summary      Удаление программы
// @Description  Удаление программы тренировок по ID
// @Tags         program
// @Security     BearerAuth
// @Param        id path int true "ID программы тренировок"
// @Success      204 "Успешное удаление программы"
// @Failure      400 {string} string "Bad request"
// @Failure      401 {string} string "Unauthorized"
// @Failure      404 {string} string "Program not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /programs/{id} [delete]
func (h *ProgramHTTPHandler) DeleteProgram(rw http.ResponseWriter, r *http.Request) {
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

	if err := h.programService.DeleteProgram(
		r.Context(),
		userID,
		programID,
	); err != nil {
		core_http_errors.WriteError(rw, err)

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
