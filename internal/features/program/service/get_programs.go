package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *ProgramService) GetPrograms(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.Program, error) {
	programs, err := s.programRepository.GetPrograms(
		ctx,
		userID,
	)
	if err != nil {
		return []core_domain.Program{}, err
	}

	return programs, nil
}
