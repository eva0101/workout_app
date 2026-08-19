package core_errors

import "errors"

var (
	ErrPasswordFewSymbols  = errors.New("the password must be more than 8 characters long")
	ErrPasswordManySymbols = errors.New("the password must be less than 64 characters")

	ErrLoginFewSymbols  = errors.New("the login must be more than 3 characters long")
	ErrLoginManySymbols = errors.New("the login must be less than 32 characters")
	ErrLoginExists      = errors.New("user already exists")

	ErrInvalidCredentials = errors.New("login or password is incorrect")
	ErrInvalidArgument    = errors.New("invalid argument")

	ErrInvalidProgramName = errors.New("program name have between 3 and 100 symbols")

	ErrProgramNotFound     = errors.New("program not found")
	ErrTrainingDayNotFound = errors.New("training day not found")
	ErrExerciseNotFound    = errors.New("exercise not found")

	ErrWorkoutAlreadyInProgress = errors.New("there is already a running workout")
)
