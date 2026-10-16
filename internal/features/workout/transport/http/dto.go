package workout_transport_http

import "time"

type ExecuteSetDTORequest struct {
	RepsDone *int     `json:"reps_done" example:"4"`
	Weight   *float64 `json:"weight"    example:"8"`
}
type ExecuteSetDTOResponse struct {
	TrainingDayExercisesID int      `json:"training_day_exercises_id" example:"1"`
	SetNumber              int      `json:"set_number"                example:"5"`
	RepsDone               *int     `json:"reps_done"                 example:"4"`
	Weight                 *float64 `json:"weight"                    example:"8"`
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

type StatisticDTOResponse struct {
	ExerciseID     int     `json:"exercise_id"     example:"4"`
	TotalSets      int     `json:"total_sets"      example:"5"`
	TotalRepeats   int     `json:"total_repeats"   example:"20"`
	MaximumRepeats int     `json:"maximum_repeats" example:"4"`
	TotalVolume    float64 `json:"total_volume"    example:"160"`
	MaximumWeight  float64 `json:"maximum_weight"  example:"8"`
}
