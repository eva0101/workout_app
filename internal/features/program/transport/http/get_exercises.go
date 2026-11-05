package program_transport_http

import (
	"encoding/json"
	"net/http"
	"strconv"
	core_errors "workout_app/internal/core/errors"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"go.uber.org/zap"
)

// GetExercises godoc
// @Summary     Получение упражнений каталога
// @Description Получение всех упражнений или упражнений выбранной группы мышц
// @Tags        program
// @Produce     json
// @Security    BearerAuth
// @Param       muscle_group_id query int false "ID группы мышц"
// @Success     200 {array} CatalogExerciseDTOResponse "Упражнения"
// @Failure     400 {string} string "Invalid muscle_group_id"
// @Failure     500 {string} string "Internal server error"
// @Router      /exercises [get]
func (h *ProgramHTTPHandler) GetExercises(rw http.ResponseWriter, r *http.Request) {
	var muscleGroupID *int
	if value := r.URL.Query().Get("muscle_group_id"); value != "" {
		id, err := strconv.Atoi(value)
		if err != nil || id < 1 {
			core_http_errors.WriteError(rw, core_errors.ErrInvalidArgument)
			return
		}
		muscleGroupID = &id
	}

	exercises, err := h.programService.GetExercises(r.Context(), muscleGroupID)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(ToCatalogExerciseResponse(exercises)); err != nil {
		h.log.Error("failed to encode exercises response", zap.Error(err))
	}
}
