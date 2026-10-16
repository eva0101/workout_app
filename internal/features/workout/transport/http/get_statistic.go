package workout_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetStatistic godoc
// @Summary     Получение статистики
// @Description Получение общей статистики пользователя
// @Tags        workout
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} StatisticDTOResponse "Статистика успешно получена"
// @Failure     400 {string} string "Bad request"
// @Failure     404 {string} string "Not found"
// @Failure     500 {string} string "Internal server error"
// @Router      /workouts/statistics [get]
func (h *WorkoutHTTPHandler) GetStatistic(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	statistic, err := h.workoutService.GetStatistic(
		r.Context(),
		userID,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	if err := json.NewEncoder(rw).Encode(statistic); err != nil {
		h.log.Error(
			"failed to encode statistic",
			zap.Error(err),
		)

		return
	}
}
