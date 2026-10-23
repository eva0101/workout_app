package schedule_transport_http

import (
	"encoding/json"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetSchedule godoc
// @Summary Get user's training schedule
// @Description Returns scheduled workouts in an inclusive date range.
// @Tags schedule
// @Produce json
// @Security BearerAuth
// @Param from query string true "Start date (YYYY-MM-DD)"
// @Param to query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} ScheduleDTOResponse
// @Failure 400 {string} string "Invalid date or range"
// @Failure 401 {string} string "Unauthorized"
// @Router /workouts/schedule [get]
func (h *ScheduleHTTPHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.service.GetSchedule(
		r.Context(),
		userID,
		r.URL.Query().Get("from"),
		r.URL.Query().Get("to"),
	)
	if err != nil {
		core_http_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toScheduleResponses(items)); err != nil {
		h.log.Error("failed to encode schedule response", zap.Error(err))
	}
}
