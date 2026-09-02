package workout_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *WorkoutRepository) CompleteTraining(
	ctx context.Context,
	userID uuid.UUID,
	workoutID int,
	fatigueScore int,
) (core_domain.Workout, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error(
			"failed to begin completed training transaction",
			zap.Error(err),
		)

		return core_domain.Workout{}, err
	}
	defer tx.Rollback(ctx)

	var workout core_domain.Workout

	completeWorkoutQuery := `
	UPDATE workoutapp.workout
	SET 
		status = 'completed',
		fatigue_score = $1,
		completed_at = NOW(),
		total_time = NOW() - begin_at
	WHERE id = $2
		AND user_id = $3
		AND status = 'in_progress'
	RETURNING
		training_day_id,
		status,
		fatigue_score,
		total_time;
	`

	err = tx.QueryRow(
		ctx,
		completeWorkoutQuery,
		fatigueScore,
		workoutID,
		userID,
	).Scan(
		&workout.TrainingDayID,
		&workout.Status,
		&workout.FatigueScore,
		&workout.TotalTime,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return core_domain.Workout{}, core_errors.ErrWorkoutNotFound
		}

		r.log.Error(
			"failed to complete training",
			zap.Error(err),
		)

		return core_domain.Workout{}, err
	}

	updateStatisticsQuery := `
	INSERT INTO workoutapp.statistics (
		user_id,
		exercise_id,
		total_sets,
		total_repeats,
		maximum_repeats,
		total_volume,
		maximum_weight
	)
	SELECT
		$1,
		tde.exercise_id,
		COUNT(ws.id),
		COALESCE(SUM(ws.reps_done), 0),
		COALESCE(MAX(ws.reps_done), 0),
		COALESCE(SUM(ws.reps_done * ws.weight), 0),
		COALESCE(MAX(ws.weight), 0)
	FROM workoutapp.workout_set ws
	JOIN workoutapp.training_day_exercises tde
		ON tde.id = ws.training_day_exercises_id
	WHERE ws.workout_id = $2
	GROUP BY tde.exercise_id

	ON CONFLICT (user_id, exercise_id)
	DO UPDATE SET
		total_sets =
			workoutapp.statistics.total_sets
			+ EXCLUDED.total_sets,

		total_repeats =
			workoutapp.statistics.total_repeats
			+ EXCLUDED.total_repeats,

		maximum_repeats = GREATEST(
			workoutapp.statistics.maximum_repeats,
			EXCLUDED.maximum_repeats
		),

		total_volume =
			workoutapp.statistics.total_volume
			+ EXCLUDED.total_volume,

		maximum_weight = GREATEST(
			workoutapp.statistics.maximum_weight,
			EXCLUDED.maximum_weight
		);
	`
	_, err = tx.Exec(
		ctx,
		updateStatisticsQuery,
		userID,
		workoutID,
	)
	if err != nil {
		r.log.Error(
			"failed to update training statistics",
			zap.Error(err),
		)

		return core_domain.Workout{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		r.log.Error(
			"failed to commit complete training transaction",
			zap.Error(err),
		)

		return core_domain.Workout{}, err
	}

	return workout, nil
}
