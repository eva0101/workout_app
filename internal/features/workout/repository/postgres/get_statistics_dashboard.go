package workout_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *WorkoutRepository) GetStatisticsDashboardData(
	ctx context.Context,
	userID uuid.UUID,
	from *time.Time,
	to time.Time,
) ([]core_domain.WorkoutStatisticsData, []core_domain.WorkoutStatisticsMuscleGroup, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		r.log.Error("failed to begin workout statistics transaction", zap.Error(err))
		return nil, nil, fmt.Errorf("begin workout statistics transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	workouts := make([]core_domain.WorkoutStatisticsData, 0)
	rows, err := tx.Query(ctx, `
		SELECT
			w.completed_at,
			FLOOR(EXTRACT(EPOCH FROM (w.completed_at - w.begin_at)))::BIGINT,
			COALESCE(actual.completed_sets, 0)::BIGINT,
			COALESCE(actual.total_weight, 0)::FLOAT8,
			COALESCE(actual.actual_reps, 0)::BIGINT
		FROM workoutapp.workout w
		LEFT JOIN LATERAL (
			SELECT
				COUNT(*) AS completed_sets,
				SUM(COALESCE(ws.weight * ws.reps_done, 0)) AS total_weight,
				SUM(COALESCE(ws.reps_done, 0)) AS actual_reps
			FROM workoutapp.workout_set ws
			WHERE ws.workout_id = w.id
		) actual ON TRUE
		WHERE w.user_id = $1
			AND w.status = 'completed'
			AND w.completed_at IS NOT NULL
			AND ($2::TIMESTAMPTZ IS NULL OR w.completed_at >= $2)
			AND w.completed_at <= $3
		ORDER BY w.completed_at;
	`, userID, from, to)
	if err != nil {
		r.log.Error("failed to query workout statistics", zap.Error(err))
		return nil, nil, fmt.Errorf("query workout statistics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var workout core_domain.WorkoutStatisticsData
		if err := rows.Scan(
			&workout.CompletedAt,
			&workout.DurationSeconds,
			&workout.CompletedSets,
			&workout.TotalWeight,
			&workout.ActualReps,
		); err != nil {
			r.log.Error("failed to scan workout statistics", zap.Error(err))
			return nil, nil, fmt.Errorf("scan workout statistics: %w", err)
		}
		workouts = append(workouts, workout)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate workout statistics", zap.Error(err))
		return nil, nil, fmt.Errorf("iterate workout statistics: %w", err)
	}
	rows.Close()

	groups := make([]core_domain.WorkoutStatisticsMuscleGroup, 0, 7)
	rows, err = tx.Query(ctx, `
		WITH group_totals AS (
			SELECT we.muscle_group_name, COUNT(*) AS completed_sets
			FROM workoutapp.workout w
			JOIN workoutapp.workout_set ws ON ws.workout_id = w.id
			JOIN workoutapp.workout_exercise we
				ON we.id = ws.workout_exercise_id
				AND we.workout_id = w.id
			WHERE w.user_id = $1
				AND w.status = 'completed'
				AND w.completed_at IS NOT NULL
				AND ($2::TIMESTAMPTZ IS NULL OR w.completed_at >= $2)
				AND w.completed_at <= $3
			GROUP BY we.muscle_group_name
		)
		SELECT mg.name, mg.title, COALESCE(gt.completed_sets, 0)::BIGINT
		FROM workoutapp.muscle_group mg
		LEFT JOIN group_totals gt ON gt.muscle_group_name = mg.name
		WHERE mg.name IN ('chest', 'back', 'legs', 'biceps', 'triceps', 'shoulders', 'abs')
		ORDER BY mg.id;
	`, userID, from, to)
	if err != nil {
		r.log.Error("failed to query workout statistics muscle groups", zap.Error(err))
		return nil, nil, fmt.Errorf("query workout statistics muscle groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group core_domain.WorkoutStatisticsMuscleGroup
		if err := rows.Scan(&group.Name, &group.Title, &group.CompletedSets); err != nil {
			r.log.Error("failed to scan workout statistics muscle group", zap.Error(err))
			return nil, nil, fmt.Errorf("scan workout statistics muscle group: %w", err)
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate workout statistics muscle groups", zap.Error(err))
		return nil, nil, fmt.Errorf("iterate workout statistics muscle groups: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit workout statistics transaction", zap.Error(err))
		return nil, nil, fmt.Errorf("commit workout statistics transaction: %w", err)
	}

	return workouts, groups, nil
}
