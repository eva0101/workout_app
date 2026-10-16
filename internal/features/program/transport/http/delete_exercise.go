package program_transport_http

import (
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

// DeleteExercise godoc
// @Summary       Удаление упражнения
// @Description   Удаление упражнения в тренировочном дне по ID
// @Tags          program
// @Security      BearerAuth
// @Param         id path int true "ID удаляемого упражнения"
// @Success       204 "Успешное удаление упражнения"
// @Failure       400 {string} string "Bad request"
// @Failure       401 {string} string "Unauthorized"
// @Failure       404 {string} string "Exercise not found"
// @Failure       500 {string} string "Internal server error"
// @Router        /days/exercises/{id} [delete]
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
