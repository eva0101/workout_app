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

	default:
		http.Error(rw, "internal error", http.StatusInternalServerError)
	}
}
