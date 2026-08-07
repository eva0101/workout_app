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

type TrainingDays struct {
	ID        int
	ProgramID int
	DayNumber int
}

type TrainingDay struct {
	ID        int
	DayNumber int
	Exercises []TrainingExercise
}

type TrainingExercise struct {
	ID         int
	ExerciseID int
	Sets       int
}

type Exercise struct {
	ID           int
	NameExercise string
	Sets         int
}

type NewProgramName struct {
	Name string
}
