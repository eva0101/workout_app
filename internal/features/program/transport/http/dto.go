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
	ID          int                  `json:"id"`
	ExerciseID  int                  `json:"exercise_id"`
	Sets        int                  `json:"sets"`
	PlannedSets []PlannedSetResponse `json:"planned_sets"`
}

type CreateExerciseDTORequest struct {
	ExerciseID int                          `json:"exercise_id" example:"3"`
	Sets       []CreatePlannedSetDTORequest `json:"sets"`
}

type CreatePlannedSetDTORequest struct {
	Repetitions *int     `json:"repetitions" example:"10"`
	Weight      *float64 `json:"weight"      example:"40"`
}

type PatchExerciseSetsDTORequest struct {
	Sets []PatchPlannedSetDTORequest `json:"sets"`
}

type PatchPlannedSetDTORequest struct {
	Reps   *int     `json:"reps"   example:"10"`
	Weight *float64 `json:"weight" example:"40"`
}

type PlannedSetResponse struct {
	SetNumber   int      `json:"set_number"   example:"1"`
	Repetitions *int     `json:"repetitions"  example:"10"`
	Weight      *float64 `json:"weight"       example:"40"`
}

type CreateExerciseDTOResponse struct {
	ID           int                  `json:"id"            example:"1"`
	ExerciseID   int                  `json:"exercise_id"   example:"3"`
	NameExercise string               `json:"name_exercise" example:"Подтягивания на бицепс"`
	Sets         int                  `json:"sets"          example:"4"`
	PlannedSets  []PlannedSetResponse `json:"planned_sets"`
}

type MuscleGroupDTOResponse struct {
	ID        int     `json:"id"         example:"1"`
	Name      string  `json:"name"       example:"chest"`
	Title     string  `json:"title"      example:"Грудь"`
	ImagePath *string `json:"image_path" example:"/exercises/groups/chest.webp"`
}

type CatalogExerciseDTOResponse struct {
	ID            int     `json:"id"              example:"1"`
	Name          string  `json:"name"            example:"barbell bench press"`
	Title         string  `json:"title"           example:"Жим штанги лёжа"`
	MuscleGroupID int     `json:"muscle_group_id" example:"1"`
	Description   *string `json:"description"`
	ImagePath     *string `json:"image_path"      example:"/exercises/chest/barbell-bench-press"`
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
			ID:          exercise.ID,
			ExerciseID:  exercise.ExerciseID,
			Sets:        exercise.Sets,
			PlannedSets: ToPlannedSetResponse(exercise.PlannedSets),
		})
	}

	return response
}

func ToPlannedSetResponse(sets []core_domain.PlannedSet) []PlannedSetResponse {
	result := make([]PlannedSetResponse, 0, len(sets))
	for _, set := range sets {
		result = append(result, PlannedSetResponse{
			SetNumber:   set.SetNumber,
			Repetitions: set.Repetitions,
			Weight:      set.Weight,
		})
	}
	return result
}

func ToMuscleGroupResponse(groups []core_domain.MuscleGroup) []MuscleGroupDTOResponse {
	result := make([]MuscleGroupDTOResponse, 0, len(groups))
	for _, group := range groups {
		result = append(result, MuscleGroupDTOResponse{
			ID:        group.ID,
			Name:      group.Name,
			Title:     group.Title,
			ImagePath: group.ImagePath,
		})
	}
	return result
}

func ToCatalogExerciseResponse(exercises []core_domain.CatalogExercise) []CatalogExerciseDTOResponse {
	result := make([]CatalogExerciseDTOResponse, 0, len(exercises))
	for _, exercise := range exercises {
		result = append(result, CatalogExerciseDTOResponse{
			ID:            exercise.ID,
			Name:          exercise.Name,
			Title:         exercise.Title,
			MuscleGroupID: exercise.MuscleGroupID,
			Description:   exercise.Description,
			ImagePath:     exercise.ImagePath,
		})
	}
	return result
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
