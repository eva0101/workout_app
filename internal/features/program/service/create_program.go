package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *ProgramService) CreateProgram(
	ctx context.Context,
	userID uuid.UUID,
	name string,
) (core_domain.Program, error) {
	program, err := s.programRepository.CreateProgram(
		ctx,
		userID,
		name,
	)
	if err != nil {
		return core_domain.Program{}, err
	}

	return program, nil
}
