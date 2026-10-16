package workout_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// StartTraining godoc
// @Summary      Начать тренировку
// @Description  Запускает начало тренировочного дня
// @Tags         workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID тренировочного дня"
// @Success      201 {object} StartWorkoutDTOResponse "Успешно начатый тренировочный день"
// @Failure      400 {string} string "Bad request"
// @Failure      401 {string} string "Unauthorized"
// @Failure      409 {string} string "Conflict"
// @Failure      500 {string} string "Internal server error"
// @Router       /days/{id}/workouts [post]
func (h *WorkoutHTTPHandler) StartTraining(rw http.ResponseWriter, r *http.Request) {
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

	startTraining, err := h.workoutService.StartTraining(
		r.Context(),
		userID,
		dayID,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(startTraining); err != nil {
		h.log.Error(
			"failed to encode start training",
			zap.Error(err),
		)

		return
	}
}
