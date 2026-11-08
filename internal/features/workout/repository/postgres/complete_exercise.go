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

func (r *WorkoutRepository) CompleteExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutExerciseID int,
) (core_domain.WorkoutExerciseCompletion, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error("failed to begin complete exercise transaction", zap.Error(err))
		return core_domain.WorkoutExerciseCompletion{}, err
	}
	defer tx.Rollback(context.Background())

	var workoutStatus string
	err = tx.QueryRow(ctx, `
		SELECT w.status
		FROM workoutapp.workout_exercise we
		JOIN workoutapp.workout w ON w.id = we.workout_id
		WHERE we.id = $1 AND w.user_id = $2
		FOR UPDATE OF w, we;
	`, workoutExerciseID, userID).Scan(&workoutStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.WorkoutExerciseCompletion{}, core_errors.ErrWorkoutExerciseNotFound
		}
		r.log.Error("failed to find workout exercise", zap.Error(err))
		return core_domain.WorkoutExerciseCompletion{}, err
	}
	if workoutStatus != "in_progress" {
		return core_domain.WorkoutExerciseCompletion{}, core_errors.ErrWorkoutExerciseAlreadyCompleted
	}

	var completion core_domain.WorkoutExerciseCompletion
	err = tx.QueryRow(ctx, `
		UPDATE workoutapp.workout_exercise
		SET completed_at = COALESCE(completed_at, NOW())
		WHERE id = $1
		RETURNING id, completed_at;
	`, workoutExerciseID).Scan(&completion.ID, &completion.CompletedAt)
	if err != nil {
		r.log.Error("failed to complete workout exercise", zap.Error(err))
		return core_domain.WorkoutExerciseCompletion{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit complete exercise transaction", zap.Error(err))
		return core_domain.WorkoutExerciseCompletion{}, err
	}
	return completion, nil
}
