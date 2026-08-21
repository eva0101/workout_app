package workout_repository_postgres

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (r *WorkoutRepository) GetStatistic(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.Statistic, error) {

	statistics := make([]core_domain.Statistic, 0)

	query := `
	SELECT
		exercise_id,
		total_sets,
		total_repeats,
		maximum_repeats,
		total_volume,
		maximum_weight
	FROM workoutapp.statistics
	WHERE user_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var statistic core_domain.Statistic

		err := rows.Scan(
			&statistic.ExerciseID,
			&statistic.TotalSets,
			&statistic.TotalRepeats,
			&statistic.MaximumRepeats,
			&statistic.TotalVolume,
			&statistic.MaximumWeight,
		)
		if err != nil {
			return nil, err
		}

		statistics = append(statistics, statistic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statistics, nil
}
