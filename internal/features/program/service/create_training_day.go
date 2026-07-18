package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *ProgramService) CreateTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	programID int,
) (core_domain.TrainingDay, error) {
	trainingDay, err := s.programRepository.CreateTrainingDay(
		ctx,
		userID,
		programID,
	)
	if err != nil {
		return core_domain.TrainingDay{}, err
	}

	return trainingDay, nil
}
