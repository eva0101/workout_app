package program_repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *ProgramRepository) DeleteTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin delete training day transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var lockedTrainingDayID int
	err = tx.QueryRow(
		ctx,
		`
		SELECT td.id
		FROM workoutapp.training_days td
		JOIN workoutapp.program p ON p.id = td.program_id
		WHERE td.id = $2 AND p.user_id = $1
		FOR UPDATE OF td;
		`,
		userID,
		trainingDayID,
	).Scan(&lockedTrainingDayID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"training day with id='%d': %w",
				trainingDayID,
				core_errors.ErrTrainingDayNotFound,
			)
		}
		r.log.Error("failed to lock training day before deletion", zap.Error(err))
		return fmt.Errorf("lock training day: %w", err)
	}

	var hasWorkoutHistory bool
	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM workoutapp.workout
			WHERE training_day_id = $1
			UNION ALL
			SELECT 1
			FROM workoutapp.workout_set ws
			JOIN workoutapp.training_day_exercises tde
				ON tde.id = ws.training_day_exercises_id
			WHERE tde.training_day_id = $1
		);`,
		trainingDayID,
	).Scan(&hasWorkoutHistory)
	if err != nil {
		r.log.Error("failed to check training day workout history", zap.Error(err))
		return fmt.Errorf("check workout history: %w", err)
	}
	if hasWorkoutHistory {
		return core_errors.ErrTrainingDayHasWorkoutHistory
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM workoutapp.training_days WHERE id = $1;`,
		trainingDayID,
	)
	if err != nil {
		r.log.Error("failed to delete training day", zap.Error(err))
		return fmt.Errorf("delete training day: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit delete training day transaction", zap.Error(err))
		return fmt.Errorf("commit delete training day transaction: %w", err)
	}

	return nil
}
