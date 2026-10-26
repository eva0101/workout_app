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

// GetProgram   godoc
// @Summary     Получение программы
// @Description Получение конкретной программы по ее ID
// @Tags        program
// @Produce     json
// @Param       id path int true "ID получаемой программы"
// @Security    BearerAuth
// @Success     200 {object} ProgramDetailsDTOResponse "Программа успешно найдена"
// @Failure     400 {string} string "Bad request"
// @Failure     404 {string} string "Not found"
// @Failure     500 {string} string "Internal server error"
// @Router      /programs/{id} [get]
func (h *ProgramHTTPHandler) GetProgram(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get userID path value", http.StatusBadRequest)
		return
	}

	program, trainingDays, err := h.programService.GetProgram(
		r.Context(),
		userID,
		id,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	response := ProgramDetailsDTOResponse{
		ID:           program.ID,
		Name:         program.Name,
		TrainingDays: ToTrainingDaySummariesResponse(trainingDays),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error(
			"failed to encode program response",
			zap.Error(err),
		)

		return
	}
}
