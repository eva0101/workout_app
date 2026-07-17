package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (r *ProgramRepository) GetPrograms(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.Program, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT 
			id,
			user_id,
			name,
			started_at
		FROM workoutapp.program
		WHERE user_id = $1;
		`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("select programs: %w", err)
	}
	defer rows.Close()

	var programs []core_domain.Program
	for rows.Next() {
		var program core_domain.Program

		err := rows.Scan(
			&program.ID,
			&program.UserID,
			&program.Name,
			&program.StartedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan program: %w", err)
		}

		programs = append(programs, program)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return programs, nil
}
