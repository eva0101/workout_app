package workout_transport_http

import (
	"encoding/json"
	"net/http"
	"strconv"
	core_domain "workout_app/internal/core/domain"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetWorkoutHistory godoc
// @Summary     Получение истории завершённых тренировок
// @Description Возвращает завершённые тренировки пользователя, рассчитанные по снимку тренировки
// @Tags        workout
// @Produce     json
// @Security    BearerAuth
// @Param       limit query int false "Максимальное количество тренировок (по умолчанию 20, максимум 100)"
// @Param       offset query int false "Количество пропускаемых тренировок"
// @Success     200 {array} WorkoutHistoryDTOResponse "История тренировок"
// @Failure     400 {string} string "Некорректные параметры пагинации"
// @Failure     401 {string} string "Unauthorized"
// @Failure     500 {string} string "Internal server error"
// @Router      /workouts/history [get]
func (h *WorkoutHTTPHandler) GetWorkoutHistory(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit, err := queryInt(r, "limit", core_domain.WorkoutHistoryDefaultLimit)
	if err != nil {
		http.Error(rw, "invalid limit", http.StatusBadRequest)
		return
	}
	offset, err := queryInt(r, "offset", 0)
	if err != nil {
		http.Error(rw, "invalid offset", http.StatusBadRequest)
		return
	}
	if limit < 1 || limit > core_domain.WorkoutHistoryMaxLimit || offset < 0 {
		http.Error(rw, "invalid pagination", http.StatusBadRequest)
		return
	}

	history, err := h.workoutService.GetWorkoutHistory(r.Context(), userID, limit, offset)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	response := make([]WorkoutHistoryDTOResponse, 0, len(history))
	for _, item := range history {
		response = append(response, ToWorkoutHistoryResponse(item))
	}
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error("failed to encode workout history", zap.Error(err))
	}
}

func queryInt(r *http.Request, key string, defaultValue int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}
