package core_domain

import "time"

const (
	WorkoutHistoryDefaultLimit = 20
	WorkoutHistoryMaxLimit     = 100
)

type WorkoutHistoryItem struct {
	ID                int
	TrainingDayName   *string
	StartedAt         time.Time
	CompletedAt       time.Time
	DurationSeconds   int64
	CompletedSets     int64
	TotalWeight       float64
	ActualReps        int64
	PlannedReps       int64
	CompletionPercent *int
	FatigueScore      *int
}
