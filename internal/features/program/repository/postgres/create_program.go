package program_repository_postgres

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *ProgramRepository) CreateProgram(
	ctx context.Context,
	userID uuid.UUID,
	name string,
) (core_domain.Program, error) {
	var program core_domain.Program

	err := r.pool.QueryRow(
		ctx,
		`
		INSERT INTO workoutapp.program (
			user_id,
			name,
			started_at
		)
		VALUES ($1, $2, NOW())
		RETURNING id, user_id, name, started_at
		`,
		userID,
		name,
	).Scan(
		&program.ID,
		&program.UserID,
		&program.Name,
		&program.StartedAt,
	)
	if err != nil {
		r.log.Error(
			"failed to create program",
			zap.Error(err),
		)

		return core_domain.Program{}, err
	}

	return program, nil
}
