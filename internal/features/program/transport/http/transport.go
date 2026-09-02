package program_transport_http

import (
	"context"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_logger "workout_app/internal/core/logger"
	core_http_server "workout_app/internal/core/transport/http/server"

	"github.com/google/uuid"
)

type ProgramHTTPHandler struct {
	programService ProgramService
	log            *core_logger.Logger
}

type ProgramService interface {
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

func NewProgramHTTPHandler(
	programService ProgramService,
	log *core_logger.Logger,
) *ProgramHTTPHandler {
	return &ProgramHTTPHandler{
		programService: programService,
		log:            log,
	}
}

func (h *ProgramHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/programs",
			Handler: http.HandlerFunc(h.CreateProgram),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/programs",
			Handler: http.HandlerFunc(h.GetPrograms),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/programs/{id}",
			Handler: http.HandlerFunc(h.GetProgram),
			Auth:    true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/programs/{id}/days",
			Handler: http.HandlerFunc(h.CreateTrainingDay),
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/days/{id}",
			Handler: http.HandlerFunc(h.GetTrainingDay),
			Auth:    true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/days/{id}/exercises",
			Handler: http.HandlerFunc(h.CreateExercise),
			Auth:    true,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/days/exercises/{id}",
			Handler: http.HandlerFunc(h.DeleteExercise),
			Auth:    true,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/days/{id}",
			Handler: http.HandlerFunc(h.DeleteTrainingDay),
			Auth:    true,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/programs/{id}",
			Handler: http.HandlerFunc(h.DeleteProgram),
			Auth:    true,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/programs/{id}",
			Handler: http.HandlerFunc(h.PatchProgram),
			Auth:    true,
		},
	}
}
