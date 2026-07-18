package program_transport_http

import (
	"context"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_http_server "workout_app/internal/core/transport/http/server"

	"github.com/google/uuid"
)

type ProgramHTTPHandler struct {
	programService ProgramService
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
}

func NewProgramHTTPHandler(
	programService ProgramService,
) *ProgramHTTPHandler {
	return &ProgramHTTPHandler{
		programService: programService,
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
	}
}
