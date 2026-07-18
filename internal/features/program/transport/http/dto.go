package program_transport_http

import (
	"time"
	core_domain "workout_app/internal/core/domain"
)

type CreateProgramDTORequest struct {
	Name string `json:"name"`
}

type ProgramDTOResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
}

type ProgramDetailsDTOResponse struct {
	ID           int                       `json:"id"`
	Name         string                    `json:"name"`
	TrainingDays []core_domain.TrainingDay `json:"training_days"`
}

func ToProgramResponse(
	programs []core_domain.Program,
) []ProgramDTOResponse {
	result := make([]ProgramDTOResponse, 0, len(programs))

	for _, program := range programs {
		result = append(result, ProgramDTOResponse{
			ID:        program.ID,
			Name:      program.Name,
			StartedAt: program.StartedAt,
		})
	}

	return result
}
