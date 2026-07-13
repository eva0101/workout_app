package authorization_transport_http

import (
	"context"
	"net/http"
	core_domain "workout_app/internal/core/domain"
	core_dto "workout_app/internal/core/dto"
	core_http_server "workout_app/internal/core/transport/http/server"
)

type AuthorizationHTTPHandler struct {
	authorizationService AuthorizationService
}

type AuthorizationService interface {
	RegisterUser(
		ctx context.Context,
		login string,
		password string,
	) (core_domain.User, error)

	AuthorizationUser(
		ctx context.Context,
		login string,
		password string,
	) (core_dto.DTOToken, error)
}

func NewAuthorizationHTTPHandler(
	authorizationService AuthorizationService,
) *AuthorizationHTTPHandler {
	return &AuthorizationHTTPHandler{
		authorizationService: authorizationService,
	}
}

func (h *AuthorizationHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.RegisterUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.AuthorizationUser,
		},
	}
}
