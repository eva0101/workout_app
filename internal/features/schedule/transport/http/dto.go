package schedule_transport_http

import (
	"time"
	core_domain "workout_app/internal/core/domain"
)

type CreateScheduleDTORequest struct {
	Dates []string `json:"dates" example:"2026-10-26,2026-11-02"`
}

type ScheduleDTOResponse struct {
	ID            int       `json:"id" example:"1"`
	ScheduledDate string    `json:"scheduled_date" example:"2026-10-26"`
	TrainingDayID int       `json:"training_day_id" example:"7"`
	DayNumber     int       `json:"day_number" example:"2"`
	ProgramName   string    `json:"program_name" example:"Силовая программа"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

func toScheduleResponses(items []core_domain.ScheduledWorkout) []ScheduleDTOResponse {
	response := make([]ScheduleDTOResponse, 0, len(items))
	for _, item := range items {
		response = append(response, ScheduleDTOResponse{
			ID:            item.ID,
			ScheduledDate: item.ScheduledDate.Format("2006-01-02"),
			TrainingDayID: item.TrainingDayID,
			DayNumber:     item.DayNumber,
			ProgramName:   item.ProgramName,
			CreatedAt:     item.CreatedAt,
		})
	}
	return response
}
