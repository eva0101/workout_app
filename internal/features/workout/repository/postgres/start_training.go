package workout_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *WorkoutRepository) StartTraining(
	ctx context.Context,
	userID uuid.UUID,
	dayID int,
) (core_domain.StartWorkout, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {

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

		return core_domain.StartWorkout{}, err
	}
	if activeWorkoutExists {

		return core_domain.StartWorkout{}, core_errors.ErrWorkoutAlreadyInProgress
	}

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
			if pgErr.ConstraintName == "unique_training_day" {
				return core_domain.StartWorkout{}, core_errors.ErrTrainingDayAlreadyCompleted
			}
		}

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

			return core_domain.StartWorkout{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {

		return core_domain.StartWorkout{}, err
	}

	return workout, nil
}
