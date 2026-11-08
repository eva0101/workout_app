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

// CompleteExercise godoc
// @Summary      Завершить упражнение
// @Description  Идемпотентно отмечает snapshot-упражнение завершённым в активной тренировке
// @Tags         workout
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID упражнения в snapshot тренировки"
// @Success      200 {object} CompleteWorkoutExerciseDTOResponse "Упражнение завершено"
// @Failure      400 {string} string "Bad request"
// @Failure      401 {string} string "Unauthorized"
// @Failure      404 {string} string "Упражнение не найдено"
// @Failure      409 {string} string "Тренировка уже завершена"
// @Failure      500 {string} string "Internal server error"
// @Router       /workouts/exercises/{id}/complete [post]
func (h *WorkoutHTTPHandler) CompleteExercise(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	exerciseID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(w, "invalid workout exercise id", http.StatusBadRequest)
		return
	}
	completion, err := h.workoutService.CompleteExercise(r.Context(), userID, exerciseID)
	if err != nil {
		core_http_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CompleteWorkoutExerciseDTOResponse{
		ID:          completion.ID,
		CompletedAt: completion.CompletedAt,
	}); err != nil {
		h.log.Error("failed to encode workout exercise completion", zap.Error(err))
	}
}
