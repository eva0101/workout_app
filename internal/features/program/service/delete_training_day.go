package program_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ProgramService) DeleteTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
) error {
	if err := s.programRepository.DeleteTrainingDay(
		ctx,
		userID,
		trainingDayID,
	); err != nil {

		return fmt.Errorf("delete training day: %w", err)
	}

	return nil
}
