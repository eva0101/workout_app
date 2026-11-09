package workout_transport_http

import (
	"time"
	core_domain "workout_app/internal/core/domain"
)

type ExecuteSetDTORequest struct {
	RepsDone *int     `json:"reps_done" example:"4"`
	Weight   *float64 `json:"weight"    example:"8"`
}
type ExecuteSetDTOResponse struct {
	WorkoutExerciseID int      `json:"workout_exercise_id" example:"1"`
	SetNumber         int      `json:"set_number"          example:"5"`
	RepsDone          *int     `json:"reps_done"           example:"4"`
	Weight            *float64 `json:"weight"              example:"8"`
}

type CompleteTrainingDTORequest struct {
	FatigueScore int `json:"fatigue_score" example:"6"`
}
type CompleteTrainingDTOResponse struct {
	TrainingDayID int    `json:"training_day_id" example:"1"`
	Status        string `json:"status"          example:"completed"`
	FatigueScore  int    `json:"fatigue_score"   example:"6"`
	TotalTime     string `json:"total_time"      example:"1m13.573462s"`
}

type DTOCompleteWorkoutResponse struct {
	TrainingDayID int    `json:"training_day_id"`
	Status        string `json:"status"`
	FatigueScore  int    `json:"fatigue_score"`
	TotalTime     string `json:"total_time"`
}

type StartWorkoutDTOResponse struct {
	ID            int       `json:"id"              example:"1"`
	TrainingDayID int       `json:"training_day_id" example:"1"`
	Status        string    `json:"status"          example:"in_progress"`
	BeginAt       time.Time `json:"begin_at"        example:"2026-10-15T14:36:25.742006+05:00"`
}

type WorkoutDTOResponse struct {
	ID              int                       `json:"id" example:"1"`
	TrainingDayName *string                   `json:"training_day_name" example:"День верхней части тела"`
	Status          string                    `json:"status" example:"in_progress"`
	StartedAt       time.Time                 `json:"started_at" example:"2026-10-15T14:36:25.742006+05:00"`
	CompletedAt     *time.Time                `json:"completed_at" example:"2026-10-15T15:36:25.742006+05:00"`
	Exercises       []WorkoutExerciseResponse `json:"exercises"`
}

type WorkoutExerciseResponse struct {
	ID          int                         `json:"id" example:"1"`
	ExerciseID  int                         `json:"exercise_id" example:"4"`
	Name        string                      `json:"name" example:"barbell bench press"`
	Title       string                      `json:"title" example:"Жим штанги лёжа"`
	Order       int                         `json:"order" example:"1"`
	CompletedAt *time.Time                  `json:"completed_at" example:"2026-10-15T15:36:25.742006+05:00"`
	PlannedSets []WorkoutPlannedSetResponse `json:"planned_sets"`
	ActualSets  []WorkoutActualSetResponse  `json:"actual_sets"`
}

type CompleteWorkoutExerciseDTOResponse struct {
	ID          int       `json:"id" example:"1"`
	CompletedAt time.Time `json:"completed_at" example:"2026-10-15T15:36:25.742006+05:00"`
}

type WorkoutPlannedSetResponse struct {
	ID          int      `json:"id" example:"1"`
	SetNumber   int      `json:"set_number" example:"1"`
	Repetitions *int     `json:"repetitions" example:"10"`
	Weight      *float64 `json:"weight" example:"40"`
}

type WorkoutActualSetResponse struct {
	ID        int      `json:"id" example:"1"`
	SetNumber int      `json:"set_number" example:"1"`
	RepsDone  *int     `json:"reps_done" example:"9"`
	Weight    *float64 `json:"weight" example:"37.5"`
}

type StatisticDTOResponse struct {
	ExerciseID     int     `json:"exercise_id"     example:"4"`
	TotalSets      int     `json:"total_sets"      example:"5"`
	TotalRepeats   int     `json:"total_repeats"   example:"20"`
	MaximumRepeats int     `json:"maximum_repeats" example:"4"`
	TotalVolume    float64 `json:"total_volume"    example:"160"`
	MaximumWeight  float64 `json:"maximum_weight"  example:"8"`
}

func ToWorkoutResponse(workout core_domain.WorkoutDetails) WorkoutDTOResponse {
	response := WorkoutDTOResponse{
		ID:              workout.ID,
		TrainingDayName: workout.TrainingDayName,
		Status:          workout.Status,
		StartedAt:       workout.StartedAt,
		CompletedAt:     workout.CompletedAt,
		Exercises:       make([]WorkoutExerciseResponse, 0, len(workout.Exercises)),
	}
	for _, exercise := range workout.Exercises {
		item := WorkoutExerciseResponse{
			ID:          exercise.ID,
			ExerciseID:  exercise.ExerciseID,
			Name:        exercise.Name,
			Title:       exercise.Title,
			Order:       exercise.Order,
			CompletedAt: exercise.CompletedAt,
			PlannedSets: make([]WorkoutPlannedSetResponse, 0, len(exercise.PlannedSets)),
			ActualSets:  make([]WorkoutActualSetResponse, 0, len(exercise.ActualSets)),
		}
		for _, set := range exercise.PlannedSets {
			item.PlannedSets = append(item.PlannedSets, WorkoutPlannedSetResponse{
				ID:          set.ID,
				SetNumber:   set.SetNumber,
				Repetitions: set.Repetitions,
				Weight:      set.Weight,
			})
		}
		for _, set := range exercise.ActualSets {
			item.ActualSets = append(item.ActualSets, WorkoutActualSetResponse{
				ID:        set.ID,
				SetNumber: set.SetNumber,
				RepsDone:  set.RepsDone,
				Weight:    set.Weight,
			})
		}
		response.Exercises = append(response.Exercises, item)
	}
	return response
}
