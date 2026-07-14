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
	}
}
