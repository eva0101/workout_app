package core_http_errors

import (
	"errors"
	"net/http"
	core_errors "workout_app/internal/core/errors"
)

func WriteError(rw http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core_errors.ErrLoginExists):
		http.Error(rw, err.Error(), http.StatusConflict)

	case errors.Is(err, core_errors.ErrLoginFewSymbols):
		http.Error(rw, err.Error(), http.StatusBadRequest)

	case errors.Is(err, core_errors.ErrLoginManySymbols):
		http.Error(rw, err.Error(), http.StatusBadRequest)

	case errors.Is(err, core_errors.ErrPasswordFewSymbols):
		http.Error(rw, err.Error(), http.StatusBadRequest)

	case errors.Is(err, core_errors.ErrPasswordManySymbols):
		http.Error(rw, err.Error(), http.StatusBadRequest)

	case errors.Is(err, core_errors.ErrInvalidCredentials):
		http.Error(rw, err.Error(), http.StatusConflict)

	case errors.Is(err, core_errors.ErrProgramNotFound):
		http.Error(rw, err.Error(), http.StatusNotFound)

	case errors.Is(err, core_errors.ErrTrainingDayNotFound):
		http.Error(rw, err.Error(), http.StatusNotFound)

	case errors.Is(err, core_errors.ErrExerciseNotFound):
		http.Error(rw, err.Error(), http.StatusNotFound)

	case errors.Is(err, core_errors.ErrInvalidProgramName):
		http.Error(rw, err.Error(), http.StatusBadRequest)

	case errors.Is(err, core_errors.ErrWorkoutAlreadyInProgress):
		http.Error(rw, err.Error(), http.StatusConflict)

	case errors.Is(err, core_errors.ErrActiveWorkoutNotFound):
		http.Error(rw, err.Error(), http.StatusNotFound)

	case errors.Is(err, core_errors.ErrFatigueScoreNotCorrect):
		http.Error(rw, err.Error(), http.StatusBadRequest)

	case errors.Is(err, core_errors.ErrWorkoutNotFound):
		http.Error(rw, err.Error(), http.StatusNotFound)

	case errors.Is(err, core_errors.ErrTrainingDayAlreadyCompleted):
		http.Error(rw, err.Error(), http.StatusConflict)

	default:
		http.Error(rw, "internal error", http.StatusInternalServerError)
	}
}
