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

type StartWorkout struct {
	ID            int
	TrainingDayID int
	Status        string
	BeginAt       time.Time
}

type WorkoutSet struct {
	TrainingDayExercisesID int
	SetNumber              int
	RepsDone               *int
	Weight                 *float64
}

type Workout struct {
	TrainingDayID int
	Status        string
	FatigueScore  int
	TotalTime     time.Duration
}

type Statistic struct {
	ExerciseID     int
	TotalSets      int
	TotalRepeats   int
	MaximumRepeats int
	TotalVolume    float64
	MaximumWeight  float64
}
