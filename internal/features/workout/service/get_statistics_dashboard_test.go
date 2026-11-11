package workout_service

import (
	"errors"
	"testing"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"
)

func TestStatisticsPeriodBuckets(t *testing.T) {
	now := time.Date(2026, 3, 31, 13, 15, 0, 0, time.UTC)
	tests := []struct {
		period       string
		wantInterval string
		wantBuckets  int
	}{
		{period: "7d", wantInterval: "day", wantBuckets: 7},
		{period: "30d", wantInterval: "day", wantBuckets: 30},
		{period: "3m", wantInterval: "week", wantBuckets: 14},
		{period: "1y", wantInterval: "month", wantBuckets: 13},
		{period: "all", wantInterval: "month", wantBuckets: 0},
	}

	for _, test := range tests {
		t.Run(test.period, func(t *testing.T) {
			from, interval, err := statisticsPeriod(test.period, now)
			if err != nil {
				t.Fatalf("statisticsPeriod() error = %v", err)
			}
			result := buildWorkoutStatisticsDashboard(nil, nil, test.period, from, interval, now)
			if result.Interval != test.wantInterval {
				t.Fatalf("interval = %q, want %q", result.Interval, test.wantInterval)
			}
			if len(result.Regularity) != test.wantBuckets || len(result.LoadDynamics) != test.wantBuckets {
				t.Fatalf("got %d regularity and %d load buckets, want %d",
					len(result.Regularity), len(result.LoadDynamics), test.wantBuckets)
			}
			if result.Summary != (core_domain.WorkoutStatisticsSummary{}) || len(result.MuscleGroups) != 7 {
				t.Fatalf("empty period has unexpected summary or groups: summary=%+v groups=%+v",
					result.Summary, result.MuscleGroups)
			}
			for index := range result.Regularity {
				if result.Regularity[index].WorkoutsCount != 0 ||
					result.LoadDynamics[index].TotalWeight != 0 ||
					result.LoadDynamics[index].CompletedSets != 0 ||
					result.LoadDynamics[index].ActualReps != 0 ||
					result.LoadDynamics[index].DurationSeconds != 0 {
					t.Fatalf("empty period has non-zero bucket: regularity=%+v load=%+v",
						result.Regularity[index], result.LoadDynamics[index])
				}
			}
		})
	}
}

func TestStatisticsDashboardCalculatesSummaryAndSeries(t *testing.T) {
	now := time.Date(2026, 3, 31, 13, 15, 0, 0, time.UTC)
	from, interval, err := statisticsPeriod("7d", now)
	if err != nil {
		t.Fatalf("statisticsPeriod() error = %v", err)
	}
	workouts := []core_domain.WorkoutStatisticsData{
		{
			CompletedAt:     time.Date(2026, 3, 25, 16, 0, 0, 0, time.UTC),
			DurationSeconds: 1200,
			CompletedSets:   4,
			TotalWeight:     575,
			ActualReps:      23,
		},
		{
			CompletedAt:     time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC),
			DurationSeconds: 2400,
			CompletedSets:   5,
			TotalWeight:     960,
			ActualReps:      32,
		},
	}
	groups := []core_domain.WorkoutStatisticsMuscleGroup{
		{Name: "chest", Title: "Грудь", CompletedSets: 6},
		{Name: "back", Title: "Спина", CompletedSets: 3},
	}

	result := buildWorkoutStatisticsDashboard(workouts, groups, "7d", from, interval, now)
	if result.Summary.WorkoutsCount != 2 || result.Summary.TotalDurationSeconds != 3600 ||
		result.Summary.CompletedSets != 9 || result.Summary.TotalWeight != 1535 {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
	if result.Regularity[0].Date != "2026-03-25" || result.Regularity[0].WorkoutsCount != 1 ||
		result.Regularity[len(result.Regularity)-1].Date != "2026-03-31" ||
		result.Regularity[len(result.Regularity)-1].WorkoutsCount != 1 {
		t.Fatalf("unexpected regularity series: %+v", result.Regularity)
	}
	lastLoad := result.LoadDynamics[len(result.LoadDynamics)-1]
	if lastLoad.TotalWeight != 960 || lastLoad.CompletedSets != 5 ||
		lastLoad.ActualReps != 32 || lastLoad.DurationSeconds != 2400 {
		t.Fatalf("unexpected last load bucket: %+v", lastLoad)
	}
	if len(result.MuscleGroups) != 7 || result.MuscleGroups[0].Name != "chest" ||
		result.MuscleGroups[0].Title != "Грудь" || result.MuscleGroups[0].CompletedSets != 6 ||
		result.MuscleGroups[1].Name != "back" || result.MuscleGroups[1].CompletedSets != 3 {
		t.Fatalf("unexpected muscle groups: %+v", result.MuscleGroups)
	}
}

func TestStatisticsDashboardUsesCalendarPeriodBoundaries(t *testing.T) {
	now := time.Date(2026, 3, 31, 13, 15, 0, 0, time.UTC)
	threeMonthsAgo, interval, err := statisticsPeriod("3m", now)
	if err != nil {
		t.Fatalf("statisticsPeriod(3m) error = %v", err)
	}
	if got, want := threeMonthsAgo.Format(time.RFC3339), "2025-12-31T13:15:00Z"; got != want {
		t.Fatalf("3m start = %s, want %s", got, want)
	}
	workouts := []core_domain.WorkoutStatisticsData{
		{CompletedAt: time.Date(2025, 12, 31, 23, 0, 0, 0, time.UTC), CompletedSets: 2},
		{CompletedAt: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC), CompletedSets: 3},
		{CompletedAt: time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC), CompletedSets: 4},
	}
	result := buildWorkoutStatisticsDashboard(workouts, nil, "3m", threeMonthsAgo, interval, now)
	if result.Regularity[0].Date != "2025-12-29" || result.Regularity[0].WorkoutsCount != 1 ||
		result.Regularity[1].Date != "2026-01-05" || result.Regularity[1].WorkoutsCount != 1 {
		t.Fatalf("unexpected weekly regularity: %+v", result.Regularity[:2])
	}

	yearStart, interval, err := statisticsPeriod("1y", now)
	if err != nil {
		t.Fatalf("statisticsPeriod(1y) error = %v", err)
	}
	yearResult := buildWorkoutStatisticsDashboard(nil, nil, "1y", yearStart, interval, now)
	if len(yearResult.LoadDynamics) != 13 || yearResult.LoadDynamics[0].Date != "2025-03" ||
		yearResult.LoadDynamics[len(yearResult.LoadDynamics)-1].Date != "2026-03" {
		t.Fatalf("unexpected monthly series: %+v", yearResult.LoadDynamics)
	}

	allResult := buildWorkoutStatisticsDashboard(
		[]core_domain.WorkoutStatisticsData{
			{CompletedAt: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC), CompletedSets: 2},
			{CompletedAt: time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC), CompletedSets: 3},
		},
		nil,
		"all",
		nil,
		"month",
		now,
	)
	if len(allResult.Regularity) != 3 || allResult.Regularity[0].Date != "2026-01" ||
		allResult.Regularity[1].Date != "2026-02" || allResult.Regularity[1].WorkoutsCount != 0 ||
		allResult.Regularity[2].Date != "2026-03" || allResult.Regularity[2].WorkoutsCount != 1 {
		t.Fatalf("unexpected all-time monthly series: %+v", allResult.Regularity)
	}
}

func TestStatisticsPeriodRejectsUnknownValue(t *testing.T) {
	if _, _, err := statisticsPeriod("90d", time.Now()); !errors.Is(err, core_errors.ErrInvalidStatisticsPeriod) {
		t.Fatalf("statisticsPeriod() error = %v, want invalid period", err)
	}
}
