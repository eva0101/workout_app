package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
)

func (s *ProgramService) GetMuscleGroups(
	ctx context.Context,
) ([]core_domain.MuscleGroup, error) {
	return s.programRepository.GetMuscleGroups(ctx)
}

func (s *ProgramService) GetExercises(
	ctx context.Context,
	muscleGroupID *int,
) ([]core_domain.CatalogExercise, error) {
	return s.programRepository.GetExercises(ctx, muscleGroupID)
}
