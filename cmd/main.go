package main

import (
	"context"
	"os/signal"
	"syscall"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
	core_http_server "workout_app/internal/core/transport/http/server"
	authorization_repository_postgres "workout_app/internal/features/authorization/repository/postgres"
	authorization_service "workout_app/internal/features/authorization/service"
	authorization_transport_http "workout_app/internal/features/authorization/transport/http"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	pool, err := core_pool_pgx.NewPool(
		ctx,
		core_pool_pgx.NewConfigMust(),
	)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	authorizationRepository := authorization_repository_postgres.NewAuthorizationRepository(pool)
	authorizationService := authorization_service.NewAuthorizationService(authorizationRepository)
	authorizationTransportHTTP := authorization_transport_http.NewAuthorizationHTTPHandler(authorizationService)

	router := core_http_server.NewRouter()
	router.RegisterRoutes(
		authorizationTransportHTTP.Routes()...,
	)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		router.ServeMux,
	)

	if err := httpServer.Run(ctx); err != nil {
		panic(err)
	}
}
