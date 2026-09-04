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

func (r *ProgramRepository) PatchProgram(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
	name string,
) (core_domain.NewProgramName, error) {
	var newProgramName core_domain.NewProgramName

	query := `
	UPDATE workoutapp.program p
	SET name = $3
	FROM workoutapp.users u
	WHERE p.id = $2
  		AND p.user_id = $1
	RETURNING p.name;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		userID,
		programID,
		name,
	).Scan(
		&newProgramName.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.NewProgramName{}, core_errors.ErrProgramNotFound
		}

		r.log.Error(
			"failed to update program",
			zap.Error(err),
		)

		return core_domain.NewProgramName{}, err
	}

	return newProgramName, nil
}
