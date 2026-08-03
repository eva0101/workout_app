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
	ID           int                        `json:"id"`
	Name         string                     `json:"name"`
	TrainingDays []core_domain.TrainingDays `json:"training_days"`
}

type TrainingDayResponse struct {
	ID        int                        `json:"id"`
	DayNumber int                        `json:"day_number"`
	Exercises []TrainingExerciseResponse `json:"exercises"`
}

type TrainingExerciseResponse struct {
	ID         int `json:"id"`
	ExerciseID int `json:"exercise_id"`
	Sets       int `json:"sets"`
	Reps       int `json:"reps"`
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

func ToTrainingDayResponse(day core_domain.TrainingDay) TrainingDayResponse {
	response := TrainingDayResponse{
		ID:        day.ID,
		DayNumber: day.DayNumber,
		Exercises: make([]TrainingExerciseResponse, 0, len(day.Exercises)),
	}

	for _, exercise := range day.Exercises {
		response.Exercises = append(response.Exercises, TrainingExerciseResponse{
			ID:         exercise.ID,
			ExerciseID: exercise.ExerciseID,
			Sets:       exercise.Sets,
			Reps:       exercise.Reps,
		})
	}

	return response
}
