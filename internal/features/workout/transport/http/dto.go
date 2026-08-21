package workout_transport_http

type DTOExecuteSetRequest struct {
	RepsDone *int     `json:"reps_done"`
	Weight   *float64 `json:"weight"`
}

type DTOExecuteSetResponse struct {
	TrainingDayExercisesID int  `json:"training_day_exercises_id"`
	SetNumber              *int `json:"set_number"`
	RepsDone               *int `json:"reps_done"`
	Weight                 *int `json:"weight"`
}

type DTOFatigueScoreRequest struct {
	FatigueScore int `json:"fatigue_score"`
}

type DTOCompleteWorkoutResponse struct {
	TrainingDayID int    `json:"training_day_id"`
	Status        string `json:"status"`
	FatigueScore  int    `json:"fatigue_score"`
	TotalTime     string `json:"total_time"`
}
