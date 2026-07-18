package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

type ProgramService struct {
	programRepository ProgramRepository
}

type ProgramRepository interface {
	CreateProgram(
		ctx context.Context,
		userID uuid.UUID,
		name string,
	) (core_domain.Program, error)

	GetPrograms(
		ctx context.Context,
		userID uuid.UUID,
	) ([]core_domain.Program, error)

	GetProgram(
		ctx context.Context,
		userID uuid.UUID,
		id int,
	) (core_domain.Program, []core_domain.TrainingDays, error)
}

func NewProgramService(
	programRepository ProgramRepository,
) *ProgramService {
	return &ProgramService{
		programRepository: programRepository,
	}
}
