package authorization_transport_http

import (
	"encoding/json"
	"net/http"
	core_dto "workout_app/internal/core/dto"
	core_http_errors "workout_app/internal/core/transport/http/errors"
)

// AuthorizationUser godoc
// @Summary          Авторизация пользователя
// @Description      Авторизация, генерация и сохранение jwt токена под текущего пользователя
// @Tags             authorization
// @Accept           json
// @Produce          json
// @Param            request body core_dto.DTOAuthorizationUserRequest true "AuthorizationUser тело запроса"
// @Success          200 {object} core_dto.DTOToken "Успешно авторизованный пользователь"
// @Failure          400 {string} string "Bad request"
// @Failure          409 {string} string "Conflict"
// @Failure          500 {string} string "Internal server error"
// @Router           /auth/login [post]
func (h *AuthorizationHTTPHandler) AuthorizationUser(rw http.ResponseWriter, r *http.Request) {
	var req core_dto.DTOAuthorizationUserRequest

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

	token, err := h.authorizationService.AuthorizationUser(
		r.Context(),
		req.Login,
		req.Password,
	)
	if err != nil {
		core_http_errors.WriteError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(rw).Encode(token); err != nil {
		return
	}
}
