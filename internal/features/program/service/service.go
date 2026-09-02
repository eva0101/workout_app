package program_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_logger "workout_app/internal/core/logger"

	"github.com/google/uuid"
)

type ProgramService struct {
	programRepository ProgramRepository
	log               *core_logger.Logger
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

	CreateTrainingDay(
		ctx context.Context,
		userID uuid.UUID,
		programID int,
	) (core_domain.TrainingDays, error)

	GetTrainingDay(
		ctx context.Context,
		userID uuid.UUID,
		dayID int,
	) (core_domain.TrainingDay, error)

	CreateExercise(
		ctx context.Context,
		userID uuid.UUID,
		trainingDayID int,
		exerciseID int,
		sets int,
	) (core_domain.Exercise, error)

	DeleteExercise(
		ctx context.Context,
		userID uuid.UUID,
		trainingDayExerciseID int,
	) error

	DeleteTrainingDay(
		ctx context.Context,
		userID uuid.UUID,
		trainingDayID int,
	) error

	DeleteProgram(
		ctx context.Context,
		userID uuid.UUID,
		programID int,
	) error

	PatchProgram(
		ctx context.Context,
		userID uuid.UUID,
		programID int,
		name string,
	) (core_domain.NewProgramName, error)
}

func NewProgramService(
	programRepository ProgramRepository,
	log *core_logger.Logger,
) *ProgramService {
	return &ProgramService{
		programRepository: programRepository,
		log:               log,
	}
}
