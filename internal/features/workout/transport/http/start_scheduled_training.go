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

// StartScheduledTraining godoc
// @Summary      Запустить тренировку из расписания
// @Description  Создаёт snapshot тренировки и удаляет запись расписания атомарно
// @Tags         workout
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID записи расписания"
// @Success      201 {object} StartWorkoutDTOResponse "Тренировка запущена"
// @Failure      400 {string} string "Bad request"
// @Failure      401 {string} string "Unauthorized"
// @Failure      404 {string} string "Запись расписания не найдена"
// @Failure      409 {string} string "Уже есть активная тренировка"
// @Failure      500 {string} string "Internal server error"
// @Router       /workouts/schedule/{id}/start [post]
func (h *WorkoutHTTPHandler) StartScheduledTraining(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	scheduleID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(w, "invalid schedule id", http.StatusBadRequest)
		return
	}
	workout, err := h.workoutService.StartScheduledTraining(r.Context(), userID, scheduleID)
	if err != nil {
		core_http_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(workout); err != nil {
		h.log.Error("failed to encode scheduled workout response", zap.Error(err))
	}
}
