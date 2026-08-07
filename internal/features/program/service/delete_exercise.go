package program_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ProgramService) DeleteExercise(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
) error {
	if err := s.programRepository.DeleteExercise(
		ctx,
		userID,
		trainingDayExerciseID,
	); err != nil {

		return fmt.Errorf("delete exercise: %w", err)
	}

	return nil
}
