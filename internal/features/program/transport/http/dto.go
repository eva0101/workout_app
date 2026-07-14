package program_transport_http

import (
	"time"

	"github.com/google/uuid"
)

type DTOCreateProgramRequest struct {
	Name string `json:"name"`
}

type DTOProgram struct {
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
}
