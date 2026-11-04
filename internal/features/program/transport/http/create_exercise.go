package program_transport_http

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

// CreateExercise godoc
// @Summary       Создание упражнения
// @Description   Добавление упражнения и целевых подходов с повторениями и весом в тренировочный день. Пустой массив допустим.
// @Tags          program
// @Accept        json
// @Produce       json
// @Security      BearerAuth
// @Param         id path int true "ID тренировочного дня"
// @Param         request body CreateExerciseDTORequest true "CreateExercise тело запроса"
// @Success       201 {object} CreateExerciseDTOResponse "Упражнение добавлено"
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

	plannedSets := make([]core_domain.PlannedSet, 0, len(req.Sets))
	for index, set := range req.Sets {
		plannedSets = append(plannedSets, core_domain.PlannedSet{
			SetNumber:   index + 1,
			Repetitions: set.Repetitions,
			Weight:      set.Weight,
		})
	}

	exercise, err := h.programService.CreateExercise(
		r.Context(),
		userID,
		trainingDayID,
		req.ExerciseID,
		plannedSets,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	response := CreateExerciseDTOResponse{
		ID:           exercise.ID,
		ExerciseID:   exercise.ExerciseID,
		NameExercise: exercise.NameExercise,
		Sets:         exercise.Sets,
		PlannedSets:  ToPlannedSetResponse(exercise.PlannedSets),
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error(
			"failed to encode exercise",
			zap.Error(err),
		)

		return
	}
}
