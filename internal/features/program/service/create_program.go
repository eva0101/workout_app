package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

	s.log.Info(
		"program created",
		zap.String("name_program", program.Name),
	)

	return program, nil
}
