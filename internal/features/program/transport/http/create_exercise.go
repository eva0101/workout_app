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

// CreateExercise godoc
// @Summary       Создание упражнения
// @Description   Создание нового упражнения в тренировочном дне
// @Tags          program
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         id path int true "ID тренировочного дня"
// @Param         request body CreateExerciseDTORequest true "CreateExercise тело запроса"
// @Success       201 {object} CreateExerciseDTOResponse "Успешно созданное упражнение"
// @Failure       400 {string} string "Bad request"
// @Failure       401 {string} string "Unauthorized"
// @Failure       500 {string} string "Internal server error"
// @Router        /days/{id}/exercises [post]
func (h *ProgramHTTPHandler) CreateExercise(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	trainingDayID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get traingDayID path value", http.StatusBadRequest)
		return
	}

	var req CreateExerciseDTORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	exercise, err := h.programService.CreateExercise(
		r.Context(),
		userID,
		trainingDayID,
		req.ExerciseID,
		req.Sets,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(exercise); err != nil {
		h.log.Error(
			"failed to encode exercise",
			zap.Error(err),
		)

		return
	}
}
