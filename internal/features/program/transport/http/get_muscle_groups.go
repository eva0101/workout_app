package program_transport_http

import (
	"encoding/json"
	"net/http"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"go.uber.org/zap"
)

// GetMuscleGroups godoc
// @Summary     Получение групп мышц
// @Description Получение всех групп мышц каталога упражнений
// @Tags        program
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} MuscleGroupDTOResponse "Группы мышц"
// @Failure     500 {string} string "Internal server error"
// @Router      /muscle-groups [get]
func (h *ProgramHTTPHandler) GetMuscleGroups(rw http.ResponseWriter, r *http.Request) {
	groups, err := h.programService.GetMuscleGroups(r.Context())
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(ToMuscleGroupResponse(groups)); err != nil {
		h.log.Error("failed to encode muscle groups response", zap.Error(err))
	}
}
