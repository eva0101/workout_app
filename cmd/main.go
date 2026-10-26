package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	core_logger "workout_app/internal/core/logger"
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
	schedule_repository_postgres "workout_app/internal/features/schedule/repository/postgres"
	schedule_service "workout_app/internal/features/schedule/service"
	schedule_transport_http "workout_app/internal/features/schedule/transport/http"
	workout_repository_postgres "workout_app/internal/features/workout/repository/postgres"
	workout_service "workout_app/internal/features/workout/service"
	workout_transport_http "workout_app/internal/features/workout/transport/http"

	"go.uber.org/zap"

	_ "workout_app/docs/swagger"
)

// @title                      Workout App API
// @version                    0.1
// @description                Workout App REST-API scheme
// @host                       127.0.0.1:5050
// @BasePath                   /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	config := core_logger.NewConfig()

	log, err := core_logger.NewLogger(config)
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	pool, err := core_pool_pgx.NewPool(
		ctx,
		core_pool_pgx.NewConfigMust(),
	)
	if err != nil {
		log.Fatal("failed to init postgres connection pool", zap.Error(err))
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

	authorizationRepository := authorization_repository_postgres.NewAuthorizationRepository(pool, log)
	authorizationService := authorization_service.NewAuthorizationService(authorizationRepository, jwtService, log)
	authorizationTransportHTTP := authorization_transport_http.NewAuthorizationHTTPHandler(authorizationService)

	programRepository := program_repository_postgres.NewProgramRepository(pool, log)
	programService := program_service.NewProgramService(programRepository, log)
	programTransportHTTP := program_transport_http.NewProgramHTTPHandler(programService, log)

	workoutRepository := workout_repository_postgres.NewWorkoutRepository(pool, log)
	workoutService := workout_service.NewWorkoutService(workoutRepository, log)
	workoutTransportHTTP := workout_transport_http.NewWorkoutHTTPHandler(workoutService, log)

	scheduleRepository := schedule_repository_postgres.NewScheduleRepository(pool, log)
	scheduleService := schedule_service.NewScheduleService(scheduleRepository, log)
	scheduleTransportHTTP := schedule_transport_http.NewScheduleHTTPHandler(scheduleService, log)

	routes := []core_http_server.Route{}
	routes = append(routes, authorizationTransportHTTP.Routes()...)
	routes = append(routes, programTransportHTTP.Routes()...)
	routes = append(routes, workoutTransportHTTP.Routes()...)
	routes = append(routes, scheduleTransportHTTP.Routes()...)

	router := core_http_server.NewRouter()
	router.RegisterRoutes(
		authMiddleware.Middleware,
		routes...,
	)

	log.Debug("init HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		router.ServeMux,
		core_middleware.Logger(log),
	)

	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
