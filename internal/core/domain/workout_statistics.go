package core_domain

import "time"

type WorkoutStatisticsData struct {
	CompletedAt     time.Time
	DurationSeconds int64
	CompletedSets   int64
	TotalWeight     float64
	ActualReps      int64
}

type WorkoutStatisticsMuscleGroup struct {
	Name          string
	Title         string
	CompletedSets int64
}

type WorkoutStatisticsSummary struct {
	WorkoutsCount        int64
	TotalDurationSeconds int64
	CompletedSets        int64
	TotalWeight          float64
}

type WorkoutStatisticsRegularityEntry struct {
	Date          string
	WorkoutsCount int64
}

type WorkoutStatisticsLoadEntry struct {
	Date            string
	TotalWeight     float64
	CompletedSets   int64
	ActualReps      int64
	DurationSeconds int64
}

type WorkoutStatisticsDashboard struct {
	Summary      WorkoutStatisticsSummary
	Interval     string
	Regularity   []WorkoutStatisticsRegularityEntry
	LoadDynamics []WorkoutStatisticsLoadEntry
	MuscleGroups []WorkoutStatisticsMuscleGroup
}
