package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *ProgramService) PatchExerciseSets(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
	sets []core_domain.PlannedSet,
) error {
	if err := validatePlannedSets(sets); err != nil {
		return err
	}

	return s.programRepository.PatchExerciseSets(
		ctx,
		userID,
		trainingDayExerciseID,
		sets,
	)
}
