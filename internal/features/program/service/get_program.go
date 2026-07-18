package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *ProgramService) GetProgram(
	ctx context.Context,
	userID uuid.UUID,
	id int,
) (core_domain.Program, []core_domain.TrainingDay, error) {
	program, trainingDays, err := s.programRepository.GetProgram(
		ctx,
		userID,
		id,
	)
	if err != nil {
		return core_domain.Program{}, []core_domain.TrainingDay{}, err
	}

	return program, trainingDays, nil
}
