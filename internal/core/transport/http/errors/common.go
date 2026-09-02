package core_http_errors

import (
	"errors"
	"net/http"
	core_errors "workout_app/internal/core/errors"
)

func WriteError(rw http.ResponseWriter, err error) {
	var appErr *core_errors.AppError

	if errors.As(err, &appErr) {
		http.Error(rw, appErr.Error(), appErr.Status)
		return
	}

	http.Error(rw, "internal error", http.StatusInternalServerError)
}
