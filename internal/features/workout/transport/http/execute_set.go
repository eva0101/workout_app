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

// ExecuteSet   godoc
// @Summary     Выполнить подход
// @Description Отмечает выполнение подхода, записывая кол-во повторений и вес
// @Description ***Кол-во повторений и вес необязательны, для выполнения подхода***
// @Tags        workout
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "ID тренировочного дня с упражнением"
// @Param       request body ExecuteSetDTORequest true "ExecuteSet тело запроса"
// @Success     201 {object} ExecuteSetDTOResponse "Успешно выполненный подход"
// @Failure     400 {string} string "Bad request"
// @Failure     401 {string} string "Unauthorized"
// @Failure     409 {string} string "Conflict"
// @Failure     500 {string} string "Internal server error"
// @Router      /workouts/exercises/{id}/sets [post]
func (h *WorkoutHTTPHandler) ExecuteSet(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)

		return
	}

	var req ExecuteSetDTORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	trainingDayExerciseID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get trainingDayExerciseID path value", http.StatusBadRequest)

		return
	}

	workoutSet, err := h.workoutService.ExecuteSet(
		r.Context(),
		userID,
		trainingDayExerciseID,
		req.RepsDone,
		req.Weight,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(workoutSet); err != nil {
		h.log.Error(
			"failed to encode workout set",
			zap.Error(err),
		)

		return
	}
}
