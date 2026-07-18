package core_http_server

import (
	"fmt"
	"net/http"
)

type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
	}
}

func (r *Router) RegisterRoutes(
	authMiddleware func(http.Handler) http.Handler,
	routes ...Route,
) {
	for _, route := range routes {

		handler := route.Handler

		if route.Auth {
			handler = authMiddleware(handler)
		}

		r.Handle(
			fmt.Sprintf("%s %s", route.Method, route.Path),
			handler,
		)
	}
}
