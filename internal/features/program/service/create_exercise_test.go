package program_service

import (
	"context"
	"testing"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

func TestCreateExerciseRejectsInvalidPlannedSets(t *testing.T) {
	negative := -1
	negativeWeight := -0.5
	tests := map[string][]core_domain.PlannedSet{
		"negative repetitions": {{Repetitions: &negative}},
		"negative weight":      {{Weight: &negativeWeight}},
		"too many sets":        make([]core_domain.PlannedSet, 101),
	}

	service := &ProgramService{}
	for name, sets := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := service.CreateExercise(context.Background(), uuid.Nil, 1, 1, sets)
			if err != core_errors.ErrInvalidPlannedSets {
				t.Fatalf("expected ErrInvalidPlannedSets, got %v", err)
			}
		})
	}
}

func TestValidatePlannedSetsAllowsEmptyList(t *testing.T) {
	if err := validatePlannedSets([]core_domain.PlannedSet{}); err != nil {
		t.Fatalf("expected empty planned sets to be valid, got %v", err)
	}
}

func TestPatchExerciseSetsRejectsInvalidPlannedSets(t *testing.T) {
	negative := -1
	negativeWeight := -0.5
	tests := map[string][]core_domain.PlannedSet{
		"negative repetitions": {{Repetitions: &negative}},
		"negative weight":      {{Weight: &negativeWeight}},
		"too many sets":        make([]core_domain.PlannedSet, 101),
	}

	service := &ProgramService{}
	for name, sets := range tests {
		t.Run(name, func(t *testing.T) {
			err := service.PatchExerciseSets(context.Background(), uuid.Nil, 1, sets)
			if err != core_errors.ErrInvalidPlannedSets {
				t.Fatalf("expected ErrInvalidPlannedSets, got %v", err)
			}
		})
	}
}
