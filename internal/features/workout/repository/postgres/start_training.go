package workout_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *WorkoutRepository) StartTraining(
	ctx context.Context,
	userID uuid.UUID,
	dayID int,
) (core_domain.StartWorkout, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error(
			"failed to begin start taining transaction",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}
	defer tx.Rollback(ctx)

	var programID int
	var dayNumber int

	query := `
	SELECT 
		td.program_id,
		td.day_number
	FROM workoutapp.training_days td
	JOIN workoutapp.program p
		ON p.id = td.program_id
	WHERE td.id = $1
		AND p.user_id = $2
	FOR UPDATE OF td
	`
	err = tx.QueryRow(
		ctx,
		query,
		dayID,
		userID,
	).Scan(
		&programID,
		&dayNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return core_domain.StartWorkout{}, core_errors.ErrTrainingDayNotFound
		}

		r.log.Error(
			"failed to get training day for training",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}

	var activeWorkoutExists bool

	query = `
	SELECT EXISTS (
		SELECT 1
		FROM workoutapp.workout
		WHERE user_id = $1
			AND status = 'in_progress'
	)
	`
	err = tx.QueryRow(
		ctx,
		query,
		userID,
	).Scan(&activeWorkoutExists)
	if err != nil {
		r.log.Error(
			"failed to check active workout",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}
	if activeWorkoutExists {

		return core_domain.StartWorkout{}, core_errors.ErrWorkoutAlreadyInProgress
	}

	rows, err := tx.Query(
		ctx,
		`
		SELECT id
		FROM workoutapp.training_day_exercises
		WHERE training_day_id = $1
		ORDER BY id
		FOR SHARE;
		`,
		dayID,
	)
	if err != nil {
		r.log.Error("failed to lock training day exercises", zap.Error(err))
		return core_domain.StartWorkout{}, err
	}
	for rows.Next() {
		var trainingDayExerciseID int
		if err := rows.Scan(&trainingDayExerciseID); err != nil {
			rows.Close()
			r.log.Error("failed to scan locked training day exercise", zap.Error(err))
			return core_domain.StartWorkout{}, err
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		r.log.Error("failed to iterate locked training day exercises", zap.Error(err))
		return core_domain.StartWorkout{}, err
	}
	rows.Close()

	var workout core_domain.StartWorkout

	query = `
	INSERT INTO workoutapp.workout (
		user_id,
		training_day_id,
		status,
		begin_at
	)
	VALUES ($1, $2, 'in_progress', NOW())
	RETURNING
		id,
		training_day_id,
		status,
		begin_at
	`
	err = tx.QueryRow(
		ctx,
		query,
		userID,
		dayID,
	).Scan(
		&workout.ID,
		&workout.TrainingDayID,
		&workout.Status,
		&workout.BeginAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "one_active_workout_per_user" {
				return core_domain.StartWorkout{}, core_errors.ErrWorkoutAlreadyInProgress
			}
		}

		r.log.Error(
			"failed to create workout",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO workoutapp.workout_exercise (
			workout_id,
			source_training_day_exercise_id,
			exercise_id,
			exercise_name,
			exercise_title,
			sort_order
		)
		SELECT
			$1,
			tde.id,
			e.id,
			e.name,
			e.title,
			ROW_NUMBER() OVER (ORDER BY tde.id)::INT
		FROM workoutapp.training_day_exercises tde
		JOIN workoutapp.exercise e ON e.id = tde.exercise_id
		WHERE tde.training_day_id = $2;
		`,
		workout.ID,
		dayID,
	)
	if err != nil {
		r.log.Error(
			"failed to snapshot workout exercises",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO workoutapp.workout_exercise_sets (
			workout_exercise_id,
			set_number,
			repetitions,
			weight
		)
		SELECT
			we.id,
			tes.set_number,
			tes.repetitions,
			tes.weight
		FROM workoutapp.workout_exercise we
		JOIN workoutapp.training_day_exercise_sets tes
			ON tes.training_day_exercise_id = we.source_training_day_exercise_id
		WHERE we.workout_id = $1;
		`,
		workout.ID,
	)
	if err != nil {
		r.log.Error(
			"failed to snapshot workout planned sets",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}

	query = `
	UPDATE workoutapp.program
	SET started_at = NOW()
	WHERE id = $1
		AND started_at IS NULL
	`
	if dayNumber == 1 {
		_, err = tx.Exec(
			ctx,
			query,
			programID,
		)
		if err != nil {
			r.log.Error(
				"failed to update program started_at",
				zap.Error(err),
			)

			return core_domain.StartWorkout{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		r.log.Error(
			"failed to commit start training transaction",
			zap.Error(err),
		)

		return core_domain.StartWorkout{}, err
	}

	return workout, nil
}
