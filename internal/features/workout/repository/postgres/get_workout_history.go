package workout_repository_postgres

import (
	"context"
	"fmt"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (r *WorkoutRepository) GetWorkoutHistory(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]core_domain.WorkoutHistoryItem, error) {
	history := make([]core_domain.WorkoutHistoryItem, 0)

	rows, err := r.pool.Query(ctx, `
		WITH completed_workouts AS (
			SELECT id
			FROM workoutapp.workout
			WHERE user_id = $1
				AND status = 'completed'
				AND completed_at IS NOT NULL
		),
		actual_totals AS (
			SELECT
				ws.workout_id,
				COUNT(*) AS completed_sets,
				SUM(ws.weight * ws.reps_done) AS total_weight,
				SUM(ws.reps_done) AS actual_reps
			FROM workoutapp.workout_set ws
			JOIN completed_workouts cw ON cw.id = ws.workout_id
			GROUP BY ws.workout_id
		),
		planned_totals AS (
			SELECT
				we.workout_id,
				SUM(wes.repetitions) AS planned_reps
			FROM workoutapp.workout_exercise we
			JOIN completed_workouts cw ON cw.id = we.workout_id
			JOIN workoutapp.workout_exercise_sets wes
				ON wes.workout_exercise_id = we.id
			GROUP BY we.workout_id
		)
		SELECT
			w.id,
			w.training_day_name,
			w.fatigue_score,
			w.begin_at,
			w.completed_at,
			FLOOR(EXTRACT(EPOCH FROM (w.completed_at - w.begin_at)))::BIGINT,
			COALESCE(actual.completed_sets, 0)::BIGINT,
			COALESCE(actual.total_weight, 0)::FLOAT8,
			COALESCE(actual.actual_reps, 0)::BIGINT,
			COALESCE(planned.planned_reps, 0)::BIGINT,
			CASE
				WHEN COALESCE(planned.planned_reps, 0) > 0
					THEN LEAST(
						100,
						ROUND(COALESCE(actual.actual_reps, 0) * 100.0 / planned.planned_reps)
					)::INT
					ELSE NULL
			END
		FROM completed_workouts cw
		JOIN workoutapp.workout w ON w.id = cw.id
		LEFT JOIN actual_totals actual ON actual.workout_id = w.id
		LEFT JOIN planned_totals planned ON planned.workout_id = w.id
		ORDER BY w.completed_at DESC, w.id DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		r.log.Error("failed to query workout history", zap.Error(err))
		return nil, fmt.Errorf("query workout history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item core_domain.WorkoutHistoryItem
		var completionPercent, fatigueScore pgtype.Int4
		if err := rows.Scan(
			&item.ID,
			&item.TrainingDayName,
			&fatigueScore,
			&item.StartedAt,
			&item.CompletedAt,
			&item.DurationSeconds,
			&item.CompletedSets,
			&item.TotalWeight,
			&item.ActualReps,
			&item.PlannedReps,
			&completionPercent,
		); err != nil {
			r.log.Error("failed to scan workout history item", zap.Error(err))
			return nil, fmt.Errorf("scan workout history item: %w", err)
		}
		if completionPercent.Valid {
			percent := int(completionPercent.Int32)
			item.CompletionPercent = &percent
		}
		if fatigueScore.Valid {
			score := int(fatigueScore.Int32)
			item.FatigueScore = &score
		}
		history = append(history, item)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate workout history", zap.Error(err))
		return nil, fmt.Errorf("iterate workout history: %w", err)
	}

	return history, nil
}
