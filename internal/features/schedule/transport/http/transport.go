package schedule_transport_http

import (
	"context"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_logger "workout_app/internal/core/logger"
	core_http_server "workout_app/internal/core/transport/http/server"

	"github.com/google/uuid"
)

type ScheduleService interface {
	CreateSchedule(
		ctx context.Context,
		userID uuid.UUID,
		trainingDayID int,
		dates []string,
	) ([]core_domain.ScheduledWorkout, error)

	GetSchedule(
		ctx context.Context,
		userID uuid.UUID,
		from string,
		to string,
	) ([]core_domain.ScheduledWorkout, error)

	DeleteSchedule(
		ctx context.Context,
		userID uuid.UUID,
		scheduleID int,
	) error
}

type ScheduleHTTPHandler struct {
	service ScheduleService
	log     *core_logger.Logger
}

func NewScheduleHTTPHandler(
	service ScheduleService,
	log *core_logger.Logger,
) *ScheduleHTTPHandler {
	return &ScheduleHTTPHandler{service: service, log: log}
}

func (h *ScheduleHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/days/{id}/schedule",
			Handler: http.HandlerFunc(h.CreateSchedule),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/schedule",
			Handler: http.HandlerFunc(h.GetSchedule),
			Auth:    true,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/workouts/schedule/{id}",
			Handler: http.HandlerFunc(h.DeleteSchedule),
			Auth:    true,
		},
	}
}
