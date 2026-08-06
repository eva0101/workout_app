package program_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ProgramRepository) CreateExercise(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
	exerciseID int,
	sets int,
) (core_domain.Exercise, error) {
	var exists int

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT 1
		FROM workoutapp.exercise
		WHERE id = $1;
		`,
		exerciseID,
	).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_domain.Exercise{}, core_errors.ErrExerciseNotFound
	}
	if err != nil {
		return core_domain.Exercise{}, err
	}

	query := `
	WITH inserted AS (
    INSERT INTO workoutapp.training_day_exercises (
        training_day_id,
        exercise_id,
        sets
    )
    SELECT
        td.id,
        $3,
        $4
    FROM workoutapp.training_days td
    JOIN workoutapp.program p
        ON p.id = td.program_id
    WHERE td.id = $2
        AND p.user_id = $1
    RETURNING *
)
	SELECT
		i.id,
		e.name,
		i.sets
	FROM inserted i
	JOIN workoutapp.exercise e
		ON e.id = i.exercise_id;
	`

	var exercise core_domain.Exercise

	err = r.pool.QueryRow(
		ctx,
		query,
		userID,
		trainingDayID,
		exerciseID,
		sets,
	).Scan(
		&exercise.ID,
		&exercise.NameExercise,
		&exercise.Sets,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.Exercise{}, core_errors.ErrTrainingDayNotFound
		}

		return core_domain.Exercise{}, err
	}

	return exercise, nil
}
