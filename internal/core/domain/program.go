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
	Name      string
}

type TrainingDay struct {
	ID        int
	DayNumber int
	Name      string
	Exercises []TrainingExercise
}

type TrainingExercise struct {
	ID          int
	ExerciseID  int
	Sets        int
	PlannedSets []PlannedSet
}

type Exercise struct {
	ID           int
	ExerciseID   int
	NameExercise string
	Sets         int
	PlannedSets  []PlannedSet
}

type PlannedSet struct {
	SetNumber   int
	Repetitions *int
	Weight      *float64
}

type MuscleGroup struct {
	ID        int
	Name      string
	Title     string
	ImagePath *string
}

type CatalogExercise struct {
	ID            int
	Name          string
	Title         string
	MuscleGroupID int
	Description   *string
	ImagePath     *string
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
	WorkoutExerciseID int
	SetNumber         int
	RepsDone          *int
	Weight            *float64
}

type Workout struct {
	TrainingDayID int
	Status        string
	FatigueScore  int
	TotalTime     time.Duration
}

type WorkoutDetails struct {
	ID              int
	TrainingDayName *string
	Status          string
	StartedAt       time.Time
	CompletedAt     *time.Time
	Exercises       []WorkoutExerciseSnapshot
}

type WorkoutExerciseSnapshot struct {
	ID          int
	ExerciseID  int
	Name        string
	Title       string
	Order       int
	CompletedAt *time.Time
	PlannedSets []WorkoutPlannedSet
	ActualSets  []WorkoutActualSet
}

type WorkoutExerciseCompletion struct {
	ID          int
	CompletedAt time.Time
}

type WorkoutPlannedSet struct {
	ID          int
	SetNumber   int
	Repetitions *int
	Weight      *float64
}

type WorkoutActualSet struct {
	ID        int
	SetNumber int
	RepsDone  *int
	Weight    *float64
}

type Statistic struct {
	ExerciseID     int
	TotalSets      int
	TotalRepeats   int
	MaximumRepeats int
	TotalVolume    float64
	MaximumWeight  float64
}
