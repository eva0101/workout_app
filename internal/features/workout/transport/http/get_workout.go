package workout_transport_http

import (
	"encoding/json"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetActiveWorkout godoc
// @Summary     Получение активной тренировки
// @Description Возвращает активную тренировку пользователя со snapshot-упражнениями и подходами
// @Tags        workout
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} WorkoutDTOResponse "Активная тренировка"
// @Failure     401 {string} string "Unauthorized"
// @Failure     404 {string} string "Активная тренировка не найдена"
// @Failure     500 {string} string "Internal server error"
// @Router      /workouts/active [get]
func (h *WorkoutHTTPHandler) GetActiveWorkout(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	workout, err := h.workoutService.GetActiveWorkout(r.Context(), userID)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}
	h.writeWorkoutResponse(rw, workout)
}

// GetWorkout godoc
// @Summary     Получение тренировки
// @Description Возвращает тренировку пользователя со snapshot-упражнениями и подходами
// @Tags        workout
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "ID тренировки"
// @Success     200 {object} WorkoutDTOResponse "Тренировка"
// @Failure     400 {string} string "Bad request"
// @Failure     401 {string} string "Unauthorized"
// @Failure     404 {string} string "Тренировка не найдена"
// @Failure     500 {string} string "Internal server error"
// @Router      /workouts/{id} [get]
func (h *WorkoutHTTPHandler) GetWorkout(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get workoutID path value", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutService.GetWorkout(r.Context(), userID, workoutID)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}
	h.writeWorkoutResponse(rw, workout)
}

func (h *WorkoutHTTPHandler) writeWorkoutResponse(
	rw http.ResponseWriter,
	workout core_domain.WorkoutDetails,
) {
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(ToWorkoutResponse(workout)); err != nil {
		h.log.Error("failed to encode workout response", zap.Error(err))
	}
}
