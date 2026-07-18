package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type Program struct {
	ID        int
	UserID    uuid.UUID
	Name      string
	StartedAt time.Time
}

type TrainingDay struct {
	ID        int
	ProgramID int
	DayNumber int
}
