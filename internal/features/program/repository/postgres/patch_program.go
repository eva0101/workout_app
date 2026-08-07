package program_repository_postgres

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
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

		return core_domain.NewProgramName{}, err
	}

	return newProgramName, nil
}
