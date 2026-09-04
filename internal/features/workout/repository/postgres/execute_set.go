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

func (r *WorkoutRepository) ExecuteSet(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
	repsDone *int,
	weight *float64,
) (core_domain.WorkoutSet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return core_domain.WorkoutSet{}, err
	}
	defer tx.Rollback(ctx)

	var workoutID int

	lockWorkoutQuery := `
		SELECT id
		FROM workoutapp.workout
		WHERE user_id = $1
			AND status = 'in_progress'
		FOR UPDATE;
	`

	err = tx.QueryRow(
		ctx,
		lockWorkoutQuery,
		userID,
	).Scan(&workoutID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.WorkoutSet{}, core_errors.ErrActiveWorkoutNotFound
		}

		r.log.Error(
			"failed to lock workout",
			zap.Error(err),
		)

		return core_domain.WorkoutSet{}, err
	}

	var workoutSet core_domain.WorkoutSet

	query := `
		INSERT INTO workoutapp.workout_set (
			workout_id,
			training_day_exercises_id,
			set_number,
			reps_done,
			weight
		)
		SELECT
			$1,
			tde.id,
			COALESCE(MAX(ws.set_number), 0) + 1,
			$3,
			$4
		FROM workoutapp.training_day_exercises tde
		LEFT JOIN workoutapp.workout_set ws
			ON ws.workout_id = $1
			AND ws.training_day_exercises_id = tde.id
		WHERE tde.id = $2
		GROUP BY tde.id
		RETURNING
			training_day_exercises_id,
			set_number,
			reps_done,
			weight;
	`

	err = tx.QueryRow(
		ctx,
		query,
		workoutID,
		trainingDayExerciseID,
		repsDone,
		weight,
	).Scan(
		&workoutSet.TrainingDayExercisesID,
		&workoutSet.SetNumber,
		&workoutSet.RepsDone,
		&workoutSet.Weight,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.WorkoutSet{}, core_errors.ErrActiveWorkoutNotFound
		}

		r.log.Error(
			"failed to execute set",
			zap.Error(err),
		)

		return core_domain.WorkoutSet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error(
			"failed to commit execute set transaction",
			zap.Error(err),
		)

		return core_domain.WorkoutSet{}, err
	}

	return workoutSet, nil
}
