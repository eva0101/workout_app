package program_transport_http

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateTrainingDay godoc
// @Summary          Создание тренировочного дня
// @Description      Создание тренировочного дня в программе тренировок
// @Tags             program
// @Accept           json
// @Produce          json
// @Security         BearerAuth
// @Param            id path int true "ID программы тренировок"
// @Param            request body CreateTrainingDayDTORequest true "Название тренировочного дня"
// @Success          201 {object} CreateTrainingDayDTOResponse "Успешно созданный тренировочный день"
// @Failure          400 {string} string "Bad request"
// @Failure          401 {string} string "Unauthorized"
// @Failure          500 {string} string "Internal server error"
// @Router           /programs/{id}/days [post]
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

	var req CreateTrainingDayDTORequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(req.Name)
	nameLength := utf8.RuneCountInString(name)
	if nameLength < 3 || nameLength > 100 {
		http.Error(rw, "training day name must be between 3 and 100 characters", http.StatusBadRequest)
		return
	}

	trainingDay, err := h.programService.CreateTrainingDay(
		r.Context(),
		userID,
		programID,
		name,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(ToTrainingDaySummaryResponse(trainingDay)); err != nil {
		h.log.Error(
			"failed to encode training day",
			zap.Error(err),
		)

		return
	}
}
