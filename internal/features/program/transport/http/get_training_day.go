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

// GetTrainingDay godoc
// @Summary       Получение тренировочного дня
// @Description   Получение конкретного тренировочного дня по ID
// @Tags          program
// @Produce       json
// @Param         id path int true "ID получаймого дня"
// @Security      BearerAuth
// @Success       200 {object} TrainingDayDTOResponse "Тренировочный день успешно найден"
// @Failure       400 {string} string "Bad request"
// @Failure       404 {string} string "Not found"
// @Failure       500 {string} string "Internal server error"
// @Router        /days/{id} [get]
func (h *ProgramHTTPHandler) GetTrainingDay(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	dayID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get dayID path value", http.StatusBadRequest)
		return
	}

	trainingDay, err := h.programService.GetTrainingDay(
		r.Context(),
		userID,
		dayID,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	response := ToTrainingDayResponse(trainingDay)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error(
			"failed to encode training day response",
			zap.Error(err),
		)

		return
	}
}
