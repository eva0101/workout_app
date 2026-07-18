package program_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ProgramRepository) CreateTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
) (core_domain.TrainingDay, error) {

	query := `
	INSERT INTO workoutapp.training_days(program_id, day_number)
	SELECT
		p.id,
		COALESCE(MAX(td.day_number), 0) + 1
	FROM workoutapp.program p
	LEFT JOIN workoutapp.training_days td
		ON td.program_id = p.id
	WHERE p.id = $1
	AND p.user_id = $2
	GROUP BY p.id
	RETURNING id, program_id, day_number;
	`

	var trainingDay core_domain.TrainingDay

	err := r.pool.QueryRow(
		ctx,
		query,
		programID,
		userID,
	).Scan(
		&trainingDay.ID,
		&trainingDay.ProgramID,
		&trainingDay.DayNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.TrainingDay{}, core_errors.ErrProgramNotFound
		}

		return core_domain.TrainingDay{}, err
	}

	return trainingDay, nil
}
