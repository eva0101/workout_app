package schedule_service

import (
	"context"
	"fmt"
	"net/http"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

func parseDate(value string) (time.Time, error) {
	date, err := time.Parse(dateLayout, value)
	if err != nil || date.Format(dateLayout) != value {
		return time.Time{}, core_errors.NewAppError(
			fmt.Sprintf("invalid date %q; expected YYYY-MM-DD", value),
			http.StatusBadRequest,
		)
	}
	return date, nil
}

func (s *ScheduleService) CreateSchedule(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
	dateValues []string,
) ([]core_domain.ScheduledWorkout, error) {
	if trainingDayID <= 0 {
		return nil, core_errors.NewAppError("training day id must be positive", http.StatusBadRequest)
	}
	if len(dateValues) == 0 {
		return nil, core_errors.NewAppError("dates must contain at least one date", http.StatusBadRequest)
	}

	today, err := parseDate(s.now().Format(dateLayout))
	if err != nil {
		return nil, fmt.Errorf("get current calendar date: %w", err)
	}

	dates := make([]time.Time, 0, len(dateValues))
	seen := make(map[string]struct{}, len(dateValues))
	for _, value := range dateValues {
		date, err := parseDate(value)
		if err != nil {
			return nil, err
		}
		if date.Before(today) {
			return nil, core_errors.NewAppError(
				fmt.Sprintf("scheduled date %s cannot be in the past", value),
				http.StatusBadRequest,
			)
		}
		if _, exists := seen[value]; exists {
			return nil, core_errors.NewAppError(
				fmt.Sprintf("scheduled date %s is duplicated", value),
				http.StatusConflict,
			)
		}
		seen[value] = struct{}{}
		dates = append(dates, date)
	}

	return s.repository.CreateSchedule(ctx, userID, trainingDayID, dates)
}

func (s *ScheduleService) GetSchedule(
	ctx context.Context,
	userID uuid.UUID,
	fromValue string,
	toValue string,
) ([]core_domain.ScheduledWorkout, error) {
	from, err := parseDate(fromValue)
	if err != nil {
		return nil, err
	}
	to, err := parseDate(toValue)
	if err != nil {
		return nil, err
	}
	if from.After(to) {
		return nil, core_errors.NewAppError("from date must be on or before to date", http.StatusBadRequest)
	}
	return s.repository.GetSchedule(ctx, userID, from, to)
}

func (s *ScheduleService) DeleteSchedule(
	ctx context.Context,
	userID uuid.UUID,
	scheduleID int,
) error {
	if scheduleID <= 0 {
		return core_errors.NewAppError("schedule id must be positive", http.StatusBadRequest)
	}
	return s.repository.DeleteSchedule(ctx, userID, scheduleID)
}
