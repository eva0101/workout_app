package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *ProgramService) CreateExercise(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
	exerciseID int,
	sets int,
) (core_domain.Exercise, error) {
	exercise, err := s.programRepository.CreateExercise(
		ctx,
		userID,
		trainingDayID,
		exerciseID,
		sets,
	)
	if err != nil {
		return core_domain.Exercise{}, err
	}

	s.log.Info(
		"exercise created",
		zap.String("exercise", exercise.NameExercise),
	)

	return exercise, nil
}
