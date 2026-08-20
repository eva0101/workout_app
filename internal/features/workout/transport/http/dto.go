package workout_transport_http

type ExecuteSetDTORequest struct {
	RepsDone *int     `json:"reps_done"`
	Weight   *float64 `json:"weight"`
}
