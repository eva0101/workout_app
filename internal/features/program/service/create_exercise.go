package program_service

import (
	"context"
	"math"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *ProgramService) CreateExercise(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
	exerciseID int,
	sets []core_domain.PlannedSet,
) (core_domain.Exercise, error) {
	if len(sets) > 100 {
		return core_domain.Exercise{}, core_errors.ErrInvalidPlannedSets
	}
	for _, set := range sets {
		if (set.Repetitions != nil && *set.Repetitions < 0) ||
			(set.Weight != nil && (math.IsNaN(*set.Weight) || math.IsInf(*set.Weight, 0) || *set.Weight < 0 || *set.Weight > 999.99)) {
			return core_domain.Exercise{}, core_errors.ErrInvalidPlannedSets
		}
	}

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
