package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *ProgramService) GetTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	dayID int,
) (core_domain.TrainingDay, error) {
	trainingDay, err := s.programRepository.GetTrainingDay(
		ctx,
		userID,
		dayID,
	)
	if err != nil {
		return core_domain.TrainingDay{}, err
	}

	return trainingDay, nil
}
