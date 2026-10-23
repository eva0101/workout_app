package schedule_transport_http

import (
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
)

// DeleteSchedule godoc
// @Summary Delete a scheduled workout
// @Tags schedule
// @Security BearerAuth
// @Param id path int true "Schedule record ID"
// @Success 204 "Schedule deleted"
// @Failure 400 {string} string "Invalid schedule ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Schedule not found"
// @Router /workouts/schedule/{id} [delete]
func (h *ScheduleHTTPHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	scheduleID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(w, "invalid schedule id", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteSchedule(r.Context(), userID, scheduleID); err != nil {
		core_http_errors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
