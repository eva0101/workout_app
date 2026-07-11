package authorization_transport_http

import (
	"encoding/json"
	"net/http"
	core_http_errors "workout_app/internal/core/transport/http/errors"
)

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthorizationHTTPHandler) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(rw, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(rw, "login and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authorizationService.RegisterUser(
		r.Context(),
		req.Login,
		req.Password,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(rw).Encode(user); err != nil {
		return
	}
}
