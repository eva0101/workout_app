package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (r *ProgramRepository) GetProgram(
	ctx context.Context,
	userID uuid.UUID,
	id int,
) (core_domain.Program, []core_domain.TrainingDay, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var program core_domain.Program
	trainingDays := make([]core_domain.TrainingDay, 0)

	row := r.pool.QueryRow(
		ctx,
		`
		SELECT id, name
		FROM workoutapp.program
		WHERE id = $1
		AND user_id = $2;
		`,
		id,
		userID,
	)

	if err := row.Scan(&program.ID, &program.Name); err != nil {
		return core_domain.Program{}, nil, core_errors.ErrProgramNotFound
	}

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT id, program_id, day_number
		FROM workoutapp.training_days
		WHERE program_id = $1;
		`,
		id,
	)
	if err != nil {
		return core_domain.Program{}, nil, fmt.Errorf("select training days: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trainingDay core_domain.TrainingDay

		err := rows.Scan(
			&trainingDay.ID,
			&trainingDay.ProgramID,
			&trainingDay.DayNumber,
		)
		if err != nil {
			return core_domain.Program{}, nil, fmt.Errorf("scan training day: %w", err)
		}

		trainingDays = append(trainingDays, trainingDay)
	}
	if err := rows.Err(); err != nil {
		return core_domain.Program{}, nil, fmt.Errorf("next rows: %w", err)
	}

	return program, trainingDays, nil
}
