package core_errors

import (
	"errors"
	"net/http"
)

type AppError struct {
	Err    error
	Status int
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(message string, status int) error {
	return &AppError{
		Err:    errors.New(message),
		Status: status,
	}
}

var (
	ErrInvalidArgument = NewAppError(
		"invalid argument",
		http.StatusBadRequest,
	)

	ErrPasswordFewSymbols = NewAppError(
		"the password must be more than 8 characters long",
		http.StatusBadRequest,
	)

	ErrPasswordManySymbols = NewAppError(
		"the password must be less than 64 characters",
		http.StatusBadRequest,
	)

	ErrLoginFewSymbols = NewAppError(
		"the login must be more than 3 characters long",
		http.StatusBadRequest,
	)

	ErrLoginManySymbols = NewAppError(
		"the login must be less than 32 characters",
		http.StatusBadRequest,
	)

	ErrLoginExists = NewAppError(
		"user already exists",
		http.StatusConflict,
	)

	ErrInvalidCredentials = NewAppError(
		"login or password is incorrect",
		http.StatusUnauthorized,
	)

	ErrProgramNotFound = NewAppError(
		"program not found",
		http.StatusNotFound,
	)

	ErrInvalidProgramName = NewAppError(
		"program name have between 3 and 100 symbols",
		http.StatusBadRequest,
	)

	ErrTrainingDayNotFound = NewAppError(
		"training day not found",
		http.StatusNotFound,
	)

	ErrExerciseNotFound = NewAppError(
		"exercise not found",
		http.StatusNotFound,
	)

	ErrTrainingDayAlreadyCompleted = NewAppError(
		"this training day was already completed",
		http.StatusConflict,
	)

	ErrWorkoutAlreadyInProgress = NewAppError(
		"there is already a running workout",
		http.StatusConflict,
	)

	ErrActiveWorkoutNotFound = NewAppError(
		"there is no active training at the moment",
		http.StatusNotFound,
	)

	ErrWorkoutNotFound = NewAppError(
		"there is no training with such ID.",
		http.StatusNotFound,
	)

	ErrFatigueScoreNotCorrect = NewAppError(
		"fatigue score must be in the range from 0 to 10",
		http.StatusBadRequest,
	)
)
