package program_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ProgramService) DeleteProgram(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
) error {
	if err := s.programRepository.DeleteProgram(
		ctx,
		userID,
		programID,
	); err != nil {

		return fmt.Errorf("delete program: %w", err)
	}

	return nil
}
