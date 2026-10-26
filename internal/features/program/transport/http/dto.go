package program_transport_http

import (
	"time"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

type CreateProgramDTORequest struct {
	Name string `json:"name" example:"Программа на развитие спины"`
}
type NewProgramNameDTOResponse struct {
	Name string `json:"name" example:"new name"`
}
type CreateProgramDTOResponse struct {
	ID        int       `json:"id"         example:"1"`
	UserID    uuid.UUID `json:"user_id"    example:"8f369546-3071-45b5-85b3-1bf4faae5c8a"`
	Name      string    `json:"name"       example:"Программа на развитие спины"`
	StartedAt time.Time `json:"started_at" example:"2026-09-23T12:52:59.667404+05:00"`
}

type ProgramDTOResponse struct {
	ID        int       `json:"id"         example:"1"`
	Name      string    `json:"name"       example:"Программа на развитие спины"`
	StartedAt time.Time `json:"started_at" example:"2026-09-23T13:46:32.509142+05:00"`
}

type ProgramDetailsDTOResponse struct {
	ID           int                             `json:"id"            example:"1"`
	Name         string                          `json:"name"          example:"Программа на развитие спины"`
	TrainingDays []TrainingDaySummaryDTOResponse `json:"training_days"`
}

type TrainingDaySummaryDTOResponse struct {
	ID        int    `json:"id"         example:"1"`
	ProgramID int    `json:"program_id" example:"1"`
	DayNumber int    `json:"day_number" example:"1"`
	Name      string `json:"name"       minLength:"3" maxLength:"100" example:"День спины"`
}

type CreateTrainingDayDTORequest struct {
	Name string `json:"name" minLength:"3" maxLength:"100" example:"День спины"`
}

type CreateTrainingDayDTOResponse struct {
	ID        int    `json:"id"         example:"1"`
	ProgramID int    `json:"program_id" example:"1"`
	DayNumber int    `json:"day_number" example:"1"`
	Name      string `json:"name"       example:"День спины"`
}

type TrainingDayDTOResponse struct {
	ID        int                        `json:"id"`
	DayNumber int                        `json:"day_number"`
	Name      string                     `json:"name"`
	Exercises []TrainingExerciseResponse `json:"exercises"`
}

type TrainingExerciseResponse struct {
	ID         int `json:"id"`
	ExerciseID int `json:"exercise_id"`
	Sets       int `json:"sets"`
}

type CreateExerciseDTORequest struct {
	ExerciseID int `json:"exercise_id" example:"3"`
	Sets       int `json:"sets"        example:"4"`
}
type CreateExerciseDTOResponse struct {
	ID           int    `json:"id"            example:"1"`
	NameExercise string `json:"name_exercise" example:"Подтягивания на бицепс"`
	Sets         int    `json:"sets"          example:"4"`
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

func ToTrainingDayResponse(day core_domain.TrainingDay) TrainingDayDTOResponse {
	response := TrainingDayDTOResponse{
		ID:        day.ID,
		DayNumber: day.DayNumber,
		Name:      day.Name,
		Exercises: make([]TrainingExerciseResponse, 0, len(day.Exercises)),
	}

	for _, exercise := range day.Exercises {
		response.Exercises = append(response.Exercises, TrainingExerciseResponse{
			ID:         exercise.ID,
			ExerciseID: exercise.ExerciseID,
			Sets:       exercise.Sets,
		})
	}

	return response
}

func ToTrainingDaySummaryResponse(day core_domain.TrainingDays) TrainingDaySummaryDTOResponse {
	return TrainingDaySummaryDTOResponse{
		ID:        day.ID,
		ProgramID: day.ProgramID,
		DayNumber: day.DayNumber,
		Name:      day.Name,
	}
}

func ToTrainingDaySummariesResponse(days []core_domain.TrainingDays) []TrainingDaySummaryDTOResponse {
	result := make([]TrainingDaySummaryDTOResponse, 0, len(days))
	for _, day := range days {
		result = append(result, ToTrainingDaySummaryResponse(day))
	}
	return result
}
