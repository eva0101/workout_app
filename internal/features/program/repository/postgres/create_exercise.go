package program_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *ProgramRepository) CreateExercise(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
	exerciseID int,
	sets []core_domain.PlannedSet,
) (core_domain.Exercise, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return core_domain.Exercise{}, err
	}
	defer tx.Rollback(ctx)

	var exerciseName string
	err = tx.QueryRow(
		ctx,
		`
		SELECT name
		FROM workoutapp.exercise
		WHERE id = $1;
		`,
		exerciseID,
	).Scan(&exerciseName)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_domain.Exercise{}, core_errors.ErrExerciseNotFound
	}
	if err != nil {
		r.log.Error(
			"failed to check exercise",
			zap.Error(err),
		)

		return core_domain.Exercise{}, err
	}

	var exercise core_domain.Exercise
	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO workoutapp.training_day_exercises (
			training_day_id,
			exercise_id,
			sets
		)
		SELECT td.id, $3, $4
		FROM workoutapp.training_days td
		JOIN workoutapp.program p ON p.id = td.program_id
		WHERE td.id = $2 AND p.user_id = $1
		RETURNING id, sets;
		`,
		userID,
		trainingDayID,
		exerciseID,
		len(sets),
	).Scan(
		&exercise.ID,
		&exercise.Sets,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.Exercise{}, core_errors.ErrTrainingDayNotFound
		}

		r.log.Error(
			"failed to create exercise in training day",
			zap.Error(err),
		)

		return core_domain.Exercise{}, err
	}

	exercise.ExerciseID = exerciseID
	exercise.NameExercise = exerciseName
	exercise.PlannedSets = make([]core_domain.PlannedSet, 0, len(sets))
	for index, set := range sets {
		set.SetNumber = index + 1
		_, err := tx.Exec(
			ctx,
			`
			INSERT INTO workoutapp.training_day_exercise_sets (
				training_day_exercise_id,
				set_number,
				repetitions,
				weight
			) VALUES ($1, $2, $3, $4);
			`,
			exercise.ID,
			set.SetNumber,
			set.Repetitions,
			set.Weight,
		)
		if err != nil {
			r.log.Error("failed to create planned exercise set", zap.Error(err))
			return core_domain.Exercise{}, err
		}
		exercise.PlannedSets = append(exercise.PlannedSets, set)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit create exercise transaction", zap.Error(err))
		return core_domain.Exercise{}, err
	}

	return exercise, nil
}
