package core_errors

import "errors"

var (
	ErrPasswordFewSymbols  = errors.New("The password must be more than 8 characters long")
	ErrPasswordManySymbols = errors.New("The password must be less than 64 characters")

	ErrLoginFewSymbols  = errors.New("The login must be more than 3 characters long")
	ErrLoginManySymbols = errors.New("The login must be less than 32 characters")
	ErrLoginExists      = errors.New("User already exists")

	ErrInvalidCredentials = errors.New("Login or password is incorrect")
)
