package program_service

import (
	"context"
	"unicode/utf8"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func (s *ProgramService) PatchProgram(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
	name string,
) (core_domain.NewProgramName, error) {
	if utf8.RuneCountInString(name) < 3 || utf8.RuneCountInString(name) > 100 {
		return core_domain.NewProgramName{}, core_errors.ErrInvalidProgramName
	}

	newNameProgram, err := s.programRepository.PatchProgram(
		ctx,
		userID,
		programID,
		name,
	)
	if err != nil {

		return core_domain.NewProgramName{}, err
	}

	return newNameProgram, nil
}
