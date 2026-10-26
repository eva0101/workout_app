package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type ScheduledWorkout struct {
	ID            int       `json:"id"`
	UserID        uuid.UUID `json:"-"`
	TrainingDayID int       `json:"training_day_id"`
	ScheduledDate time.Time `json:"scheduled_date"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	DayNumber     int       `json:"day_number,omitempty"`
	ProgramName   string    `json:"program_name,omitempty"`
}
