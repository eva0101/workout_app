package program_repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *ProgramRepository) PatchExerciseSets(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
	sets []core_domain.PlannedSet,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin patch exercise sets transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var lockedExerciseEntryID int
	err = tx.QueryRow(
		ctx,
		`
		SELECT tde.id
		FROM workoutapp.training_day_exercises tde
		JOIN workoutapp.training_days td ON td.id = tde.training_day_id
		JOIN workoutapp.program p ON p.id = td.program_id
		WHERE tde.id = $2 AND p.user_id = $1
		FOR UPDATE OF tde;
		`,
		userID,
		trainingDayExerciseID,
	).Scan(&lockedExerciseEntryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"training day exercise with id='%d': %w",
				trainingDayExerciseID,
				core_errors.ErrExerciseNotFound,
			)
		}
		r.log.Error("failed to lock training day exercise", zap.Error(err))
		return fmt.Errorf("lock training day exercise: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM workoutapp.training_day_exercise_sets WHERE training_day_exercise_id = $1;`,
		lockedExerciseEntryID,
	)
	if err != nil {
		r.log.Error("failed to delete old planned exercise sets", zap.Error(err))
		return fmt.Errorf("delete old planned exercise sets: %w", err)
	}

	for index, set := range sets {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO workoutapp.training_day_exercise_sets (
				training_day_exercise_id,
				set_number,
				repetitions,
				weight
			) VALUES ($1, $2, $3, $4);
			`,
			lockedExerciseEntryID,
			index+1,
			set.Repetitions,
			set.Weight,
		)
		if err != nil {
			r.log.Error("failed to insert updated planned exercise set", zap.Error(err))
			return fmt.Errorf("insert updated planned exercise set: %w", err)
		}
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE workoutapp.training_day_exercises SET sets = $2 WHERE id = $1;`,
		lockedExerciseEntryID,
		len(sets),
	)
	if err != nil {
		r.log.Error("failed to update planned exercise set count", zap.Error(err))
		return fmt.Errorf("update planned exercise set count: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit patch exercise sets transaction", zap.Error(err))
		return fmt.Errorf("commit patch exercise sets transaction: %w", err)
	}

	return nil
}
