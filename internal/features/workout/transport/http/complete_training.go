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

// CompleteTraining godoc
// @Summary         Завершить тренировку
// @Description     Завершает тренировочный день и записывает результаты выполнения в статистику
// @Description     ***При завершении добавляется оценка усталости по 10-ти бальной шкале***
// @Tags            workout
// @Accept          json
// @Produce         json
// @Security        BearerAuth
// @Param           id path int true "ID тренировочного дня"
// @Param           request body CompleteTrainingDTORequest true "CompleteTraining тело запроса"
// @Success         201 {object} CompleteTrainingDTOResponse "Успешно завершенная тренировка"
// @Failure         400 {string} string "Bad request"
// @Failure         401 {string} string "Unauthorized"
// @Failure         409 {string} string "Conflict"
// @Failure         500 {string} string "Internal server error"
// @Router          /workouts/{id}/complete [post]
func (h *WorkoutHTTPHandler) CompleteTraining(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)

		return
	}

	var req CompleteTrainingDTORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	workoutID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get workoutID path value", http.StatusBadRequest)

		return
	}

	workout, err := h.workoutService.CompleteTraining(
		r.Context(),
		userID,
		workoutID,
		req.FatigueScore,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)

		return
	}

	response := DTOCompleteWorkoutResponse{
		TrainingDayID: workout.TrainingDayID,
		Status:        workout.Status,
		FatigueScore:  workout.FatigueScore,
		TotalTime:     workout.TotalTime.String(),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error(
			"failed to encode workout response",
			zap.Error(err),
		)

		return
	}
}
