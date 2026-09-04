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

func (r *ProgramRepository) CreateTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
) (core_domain.TrainingDays, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return core_domain.TrainingDays{}, err
	}
	defer tx.Rollback(ctx)

	var programIDLocked int

	lockProgramQuery := `
		SELECT id
		FROM workoutapp.program
		WHERE id = $1
			AND user_id = $2
		FOR UPDATE;
	`

	err = tx.QueryRow(
		ctx,
		lockProgramQuery,
		programID,
		userID,
	).Scan(&programIDLocked)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.TrainingDays{}, core_errors.ErrProgramNotFound
		}

		r.log.Error(
			"failed to lock program",
			zap.Error(err),
		)

		return core_domain.TrainingDays{}, err
	}

	var trainingDay core_domain.TrainingDays

	query := `
		INSERT INTO workoutapp.training_days (
			program_id,
			day_number
		)
		SELECT
			$1,
			COALESCE(MAX(day_number), 0) + 1
		FROM workoutapp.training_days
		WHERE program_id = $1
		RETURNING id, program_id, day_number;
	`

	err = tx.QueryRow(
		ctx,
		query,
		programIDLocked,
	).Scan(
		&trainingDay.ID,
		&trainingDay.ProgramID,
		&trainingDay.DayNumber,
	)

	if err != nil {
		r.log.Error(
			"failed to create training day",
			zap.Error(err),
		)

		return core_domain.TrainingDays{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error(
			"failed to commit create training day transaction",
			zap.Error(err),
		)

		return core_domain.TrainingDays{}, err
	}

	return trainingDay, nil
}
