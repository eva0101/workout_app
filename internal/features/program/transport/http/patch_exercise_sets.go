package program_transport_http

import (
	"encoding/json"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

// PatchExerciseSets godoc
// @Summary       Обновление плановых подходов упражнения
// @Description   Атомарно заменяет список плановых подходов упражнения в тренировочном дне
// @Tags          program
// @Accept        json
// @Security      BearerAuth
// @Param         id path int true "ID записи упражнения в тренировочном дне"
// @Param         request body PatchExerciseSetsDTORequest true "Новый список плановых подходов"
// @Success       204 "Плановые подходы обновлены"
// @Failure       400 {string} string "Bad request"
// @Failure       401 {string} string "Unauthorized"
// @Failure       404 {string} string "Exercise not found"
// @Failure       500 {string} string "Internal server error"
// @Router        /days/exercises/{id} [patch]
func (h *ProgramHTTPHandler) PatchExerciseSets(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	trainingDayExerciseID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(rw, "failed to get trainingDayExerciseID path value", http.StatusBadRequest)
		return
	}

	var req PatchExerciseSetsDTORequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil || req.Sets == nil {
		core_http_errors.WriteError(rw, core_errors.ErrInvalidArgument)
		return
	}

	sets := make([]core_domain.PlannedSet, 0, len(req.Sets))
	for index, set := range req.Sets {
		sets = append(sets, core_domain.PlannedSet{
			SetNumber:   index + 1,
			Repetitions: set.Reps,
			Weight:      set.Weight,
		})
	}

	if err := h.programService.PatchExerciseSets(
		r.Context(),
		userID,
		trainingDayExerciseID,
		sets,
	); err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
