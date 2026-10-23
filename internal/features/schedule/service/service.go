package schedule_service

import (
	"context"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_logger "workout_app/internal/core/logger"

	"github.com/google/uuid"
)

type ScheduleRepository interface {
	CreateSchedule(
		ctx context.Context,
		userID uuid.UUID,
		trainingDayID int,
		dates []time.Time,
	) ([]core_domain.ScheduledWorkout, error)

	GetSchedule(
		ctx context.Context,
		userID uuid.UUID,
		from time.Time,
		to time.Time,
	) ([]core_domain.ScheduledWorkout, error)

	DeleteSchedule(
		ctx context.Context,
		userID uuid.UUID,
		scheduleID int,
	) error
}

type ScheduleService struct {
	repository ScheduleRepository
	log        *core_logger.Logger
	now        func() time.Time
}

func NewScheduleService(
	repository ScheduleRepository,
	log *core_logger.Logger,
) *ScheduleService {
	return &ScheduleService{
		repository: repository,
		log:        log,
		now:        time.Now,
	}
}

func NewScheduleServiceWithClock(
	repository ScheduleRepository,
	log *core_logger.Logger,
	now func() time.Time,
) *ScheduleService {
	service := NewScheduleService(repository, log)
	service.now = now
	return service
}
