package workout_service

import (
	"context"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (s *WorkoutService) GetStatisticsDashboard(
	ctx context.Context,
	userID uuid.UUID,
	period string,
) (core_domain.WorkoutStatisticsDashboard, error) {
	now := time.Now().UTC()
	from, bucketInterval, err := statisticsPeriod(period, now)
	if err != nil {
		return core_domain.WorkoutStatisticsDashboard{}, err
	}

	workouts, muscleGroups, err := s.workoutRepository.GetStatisticsDashboardData(
		ctx,
		userID,
		from,
		now,
	)
	if err != nil {
		return core_domain.WorkoutStatisticsDashboard{}, err
	}

	return buildWorkoutStatisticsDashboard(workouts, muscleGroups, period, from, bucketInterval, now), nil
}

func statisticsPeriod(period string, now time.Time) (*time.Time, string, error) {
	now = now.UTC()
	var (
		from     *time.Time
		interval string
	)
	switch period {
	case "7d":
		start := utcDate(now).AddDate(0, 0, -6)
		from, interval = &start, "day"
	case "30d":
		start := utcDate(now).AddDate(0, 0, -29)
		from, interval = &start, "day"
	case "3m":
		start := subtractCalendarMonths(now, 3)
		from, interval = &start, "week"
	case "1y":
		start := subtractCalendarMonths(now, 12)
		from, interval = &start, "month"
	case "all":
		interval = "month"
	default:
		return nil, "", core_errors.ErrInvalidStatisticsPeriod
	}
	return from, interval, nil
}

func buildWorkoutStatisticsDashboard(
	workouts []core_domain.WorkoutStatisticsData,
	muscleGroups []core_domain.WorkoutStatisticsMuscleGroup,
	period string,
	from *time.Time,
	interval string,
	now time.Time,
) core_domain.WorkoutStatisticsDashboard {
	result := core_domain.WorkoutStatisticsDashboard{
		Interval:     interval,
		Regularity:   make([]core_domain.WorkoutStatisticsRegularityEntry, 0),
		LoadDynamics: make([]core_domain.WorkoutStatisticsLoadEntry, 0),
		MuscleGroups: buildStatisticsMuscleGroups(muscleGroups),
	}

	bucketDates := make([]time.Time, 0)
	if period == "all" {
		if len(workouts) > 0 {
			earliest := workouts[0].CompletedAt.UTC()
			for _, workout := range workouts[1:] {
				if workout.CompletedAt.Before(earliest) {
					earliest = workout.CompletedAt.UTC()
				}
			}
			bucketDates = statisticsBucketDates(monthStart(earliest), monthStart(now.UTC()), interval)
		}
	} else if from != nil {
		start := *from
		switch interval {
		case "day":
			start = utcDate(start)
		case "week":
			start = weekStart(start)
		case "month":
			start = monthStart(start)
		}
		bucketDates = statisticsBucketDates(start, statisticsBucketStart(now.UTC(), interval), interval)
	}

	bucketIndexes := make(map[string]int, len(bucketDates))
	for _, bucketDate := range bucketDates {
		label := statisticsBucketLabel(bucketDate, interval)
		bucketIndexes[label] = len(result.Regularity)
		result.Regularity = append(result.Regularity, core_domain.WorkoutStatisticsRegularityEntry{
			Date: label,
		})
		result.LoadDynamics = append(result.LoadDynamics, core_domain.WorkoutStatisticsLoadEntry{
			Date: label,
		})
	}

	for _, workout := range workouts {
		result.Summary.WorkoutsCount++
		result.Summary.TotalDurationSeconds += workout.DurationSeconds
		result.Summary.CompletedSets += workout.CompletedSets
		result.Summary.TotalWeight += workout.TotalWeight

		label := statisticsBucketLabel(statisticsBucketStart(workout.CompletedAt.UTC(), interval), interval)
		index, ok := bucketIndexes[label]
		if !ok {
			continue
		}
		result.Regularity[index].WorkoutsCount++
		result.LoadDynamics[index].TotalWeight += workout.TotalWeight
		result.LoadDynamics[index].CompletedSets += workout.CompletedSets
		result.LoadDynamics[index].ActualReps += workout.ActualReps
		result.LoadDynamics[index].DurationSeconds += workout.DurationSeconds
	}

	return result
}

func buildStatisticsMuscleGroups(
	groups []core_domain.WorkoutStatisticsMuscleGroup,
) []core_domain.WorkoutStatisticsMuscleGroup {
	titles := map[string]string{
		"chest":     "Грудь",
		"back":      "Спина",
		"legs":      "Ноги",
		"biceps":    "Бицепс",
		"triceps":   "Трицепс",
		"shoulders": "Плечи",
		"abs":       "Пресс",
	}
	result := make([]core_domain.WorkoutStatisticsMuscleGroup, 0, len(titles))
	byName := make(map[string]core_domain.WorkoutStatisticsMuscleGroup, len(groups))
	for _, group := range groups {
		byName[group.Name] = group
	}
	for _, name := range []string{"chest", "back", "legs", "biceps", "triceps", "shoulders", "abs"} {
		group := core_domain.WorkoutStatisticsMuscleGroup{
			Name:  name,
			Title: titles[name],
		}
		if storedGroup, ok := byName[name]; ok {
			group = storedGroup
		}
		result = append(result, group)
	}
	return result
}

func statisticsBucketDates(start, end time.Time, interval string) []time.Time {
	dates := make([]time.Time, 0)
	for date := start; !date.After(end); date = advanceStatisticsBucket(date, interval) {
		dates = append(dates, date)
	}
	return dates
}

func advanceStatisticsBucket(date time.Time, interval string) time.Time {
	switch interval {
	case "day":
		return date.AddDate(0, 0, 1)
	case "week":
		return date.AddDate(0, 0, 7)
	default:
		return date.AddDate(0, 1, 0)
	}
}

func statisticsBucketStart(date time.Time, interval string) time.Time {
	switch interval {
	case "day":
		return utcDate(date)
	case "week":
		return weekStart(date)
	default:
		return monthStart(date)
	}
}

func statisticsBucketLabel(date time.Time, interval string) string {
	if interval == "month" {
		return date.Format("2006-01")
	}
	return date.Format("2006-01-02")
}

func utcDate(date time.Time) time.Time {
	date = date.UTC()
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

func weekStart(date time.Time) time.Time {
	date = utcDate(date)
	daysSinceMonday := (int(date.Weekday()) + 6) % 7
	return date.AddDate(0, 0, -daysSinceMonday)
}

func monthStart(date time.Time) time.Time {
	date = date.UTC()
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func subtractCalendarMonths(date time.Time, months int) time.Time {
	date = date.UTC()
	target := time.Date(date.Year(), date.Month()-time.Month(months), 1, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.UTC)
	lastDay := time.Date(target.Year(), target.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	day := date.Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(target.Year(), target.Month(), day, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.UTC)
}
