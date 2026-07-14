package main

import (
	"context"
	"os/signal"
	"syscall"
	core_middleware "workout_app/internal/core/middleware"
	"workout_app/internal/core/pkg/core_pkg_jwt"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
	core_http_server "workout_app/internal/core/transport/http/server"
	authorization_repository_postgres "workout_app/internal/features/authorization/repository/postgres"
	authorization_service "workout_app/internal/features/authorization/service"
	authorization_transport_http "workout_app/internal/features/authorization/transport/http"
	program_repository_postgres "workout_app/internal/features/program/repository/postgres"
	program_service "workout_app/internal/features/program/service"
	program_transport_http "workout_app/internal/features/program/transport/http"
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

	cfg, err := core_pkg_jwt.LoadConfig()
	if err != nil {
		panic(err)
	}
	if cfg.JWTSecret == "" {
		panic("JWT_SECRET is required")
	}

	jwtService := core_pkg_jwt.NewJWTService(cfg.JWTSecret)

	authMiddleware := core_middleware.NewAuthMiddleware(jwtService)

	authorizationRepository := authorization_repository_postgres.NewAuthorizationRepository(pool)
	authorizationService := authorization_service.NewAuthorizationService(authorizationRepository, jwtService)
	authorizationTransportHTTP := authorization_transport_http.NewAuthorizationHTTPHandler(authorizationService)

	programRepository := program_repository_postgres.NewProgramRepository(pool)
	programService := program_service.NewProgramService(programRepository)
	programTransportHTTP := program_transport_http.NewProgramHTTPHandler(programService)

	routes := append(
		authorizationTransportHTTP.Routes(),
		programTransportHTTP.Routes()...,
	)
	router := core_http_server.NewRouter()
	router.RegisterRoutes(
		authMiddleware.Middleware,
		routes...,
	)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		router.ServeMux,
	)

	if err := httpServer.Run(ctx); err != nil {
		panic(err)
	}
}
