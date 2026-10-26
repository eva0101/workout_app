package schedule_transport_http

import (
	"context"
	"net/http"
	"testing"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

type scheduleServiceStub struct{}

func (scheduleServiceStub) CreateSchedule(
	_ context.Context,
	_ uuid.UUID,
	_ int,
	_ []string,
) ([]core_domain.ScheduledWorkout, error) {
	return []core_domain.ScheduledWorkout{}, nil
}

func (scheduleServiceStub) GetSchedule(
	context.Context,
	uuid.UUID,
	string,
	string,
) ([]core_domain.ScheduledWorkout, error) {
	return []core_domain.ScheduledWorkout{}, nil
}

func (scheduleServiceStub) DeleteSchedule(context.Context, uuid.UUID, int) error {
	return nil
}

func TestScheduleRoutesRequireAuthentication(t *testing.T) {
	handler := NewScheduleHTTPHandler(scheduleServiceStub{}, nil)
	routes := handler.Routes()
	if len(routes) != 3 {
		t.Fatalf("expected 3 routes, got %d", len(routes))
	}
	for _, route := range routes {
		if !route.Auth {
			t.Errorf("route %s %s is not protected", route.Method, route.Path)
		}
	}
}

func TestScheduleRouteMethodsAndPaths(t *testing.T) {
	handler := NewScheduleHTTPHandler(scheduleServiceStub{}, nil)
	got := make(map[string]bool)
	for _, route := range handler.Routes() {
		got[route.Method+" "+route.Path] = true
	}
	want := []string{
		http.MethodPost + " /days/{id}/schedule",
		http.MethodGet + " /workouts/schedule",
		http.MethodDelete + " /workouts/schedule/{id}",
	}
	for _, route := range want {
		if !got[route] {
			t.Errorf("missing route %s", route)
		}
	}
}
