package schedule_transport_http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	core_middleware "workout_app/internal/core/middleware"
	core_http_errors "workout_app/internal/core/transport/http/errors"
	core_http_request "workout_app/internal/core/transport/http/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateSchedule godoc
// @Summary Create training-day schedule
// @Description Schedules a training day on one or more calendar dates.
// @Tags schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Training day ID"
// @Param request body CreateScheduleDTORequest true "Dates to schedule (YYYY-MM-DD)"
// @Success 201 {array} ScheduleDTOResponse
// @Failure 400 {string} string "Invalid request or past date"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Training day not found"
// @Failure 409 {string} string "Schedule already exists"
// @Router /days/{id}/schedule [post]
func (h *ScheduleHTTPHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(core_middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dayID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		http.Error(w, "invalid training day id", http.StatusBadRequest)
		return
	}

	var request CreateScheduleDTORequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		http.Error(w, "request body must contain one JSON object", http.StatusBadRequest)
		return
	}

	items, err := h.service.CreateSchedule(r.Context(), userID, dayID, request.Dates)
	if err != nil {
		core_http_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toScheduleResponses(items)); err != nil {
		h.log.Error("failed to encode schedule response", zap.Error(err))
	}
}
