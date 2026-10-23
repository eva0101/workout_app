package schedule_service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

type scheduleRepositoryStub struct {
	createCalls  int
	createdDates []time.Time
}

func (s *scheduleRepositoryStub) CreateSchedule(
	_ context.Context,
	_ uuid.UUID,
	_ int,
	dates []time.Time,
) ([]core_domain.ScheduledWorkout, error) {
	s.createCalls++
	s.createdDates = dates
	return []core_domain.ScheduledWorkout{}, nil
}

func (s *scheduleRepositoryStub) GetSchedule(
	_ context.Context,
	_ uuid.UUID,
	_ time.Time,
	_ time.Time,
) ([]core_domain.ScheduledWorkout, error) {
	return []core_domain.ScheduledWorkout{}, nil
}

func (s *scheduleRepositoryStub) DeleteSchedule(context.Context, uuid.UUID, int) error {
	return nil
}

func TestCreateScheduleValidatesDatesBeforeRepository(t *testing.T) {
	now := time.Date(2026, time.October, 23, 23, 59, 0, 0, time.FixedZone("UTC+5", 5*60*60))

	tests := []struct {
		name       string
		dates      []string
		wantStatus int
	}{
		{name: "invalid format", dates: []string{"2026-10-2"}, wantStatus: http.StatusBadRequest},
		{name: "invalid calendar date", dates: []string{"2026-02-30"}, wantStatus: http.StatusBadRequest},
		{name: "date in past", dates: []string{"2026-10-22"}, wantStatus: http.StatusBadRequest},
		{name: "duplicate dates", dates: []string{"2026-10-23", "2026-10-23"}, wantStatus: http.StatusConflict},
		{name: "empty dates", dates: []string{}, wantStatus: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &scheduleRepositoryStub{}
			service := NewScheduleServiceWithClock(repository, nil, func() time.Time { return now })

			_, err := service.CreateSchedule(context.Background(), uuid.New(), 1, test.dates)
			var appErr *core_errors.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got %v", err)
			}
			if appErr.Status != test.wantStatus {
				t.Fatalf("expected status %d, got %d", test.wantStatus, appErr.Status)
			}
			if repository.createCalls != 0 {
				t.Fatalf("repository called for invalid input %d times", repository.createCalls)
			}
		})
	}
}

func TestCreateScheduleAllowsTodayAndParsesAllDates(t *testing.T) {
	now := time.Date(2026, time.October, 23, 23, 59, 0, 0, time.FixedZone("UTC+5", 5*60*60))
	repository := &scheduleRepositoryStub{}
	service := NewScheduleServiceWithClock(repository, nil, func() time.Time { return now })

	_, err := service.CreateSchedule(
		context.Background(),
		uuid.New(),
		7,
		[]string{"2026-10-23", "2026-11-02"},
	)
	if err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	if repository.createCalls != 1 || len(repository.createdDates) != 2 {
		t.Fatalf("expected one repository call with two dates, got calls=%d dates=%d",
			repository.createCalls, len(repository.createdDates))
	}
	if got := repository.createdDates[0].Format(dateLayout); got != "2026-10-23" {
		t.Fatalf("expected current calendar date, got %s", got)
	}
}

func TestGetScheduleValidatesRange(t *testing.T) {
	tests := []struct {
		name       string
		from       string
		to         string
		wantStatus int
	}{
		{name: "invalid from", from: "2026/10/01", to: "2026-10-31", wantStatus: http.StatusBadRequest},
		{name: "invalid to", from: "2026-10-01", to: "", wantStatus: http.StatusBadRequest},
		{name: "reversed range", from: "2026-11-01", to: "2026-10-31", wantStatus: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &scheduleRepositoryStub{}
			service := NewScheduleService(repository, nil)
			_, err := service.GetSchedule(context.Background(), uuid.New(), test.from, test.to)
			var appErr *core_errors.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got %v", err)
			}
			if appErr.Status != test.wantStatus {
				t.Fatalf("expected status %d, got %d", test.wantStatus, appErr.Status)
			}
		})
	}
}

func TestGetScheduleAllowsSameDayRange(t *testing.T) {
	repository := &scheduleRepositoryStub{}
	service := NewScheduleService(repository, nil)

	if _, err := service.GetSchedule(
		context.Background(),
		uuid.New(),
		"2026-10-23",
		"2026-10-23",
	); err != nil {
		t.Fatalf("same-day range should be valid: %v", err)
	}
}
