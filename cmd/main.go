package main

import (
	"context"
	"os/signal"
	"syscall"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
	core_http_server "workout_app/internal/core/transport/http/server"
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

	httpServer := core_http_server.NewHTTPServer(core_http_server.NewConfigMust())

	if err := httpServer.Run(ctx); err != nil {
		panic(err)
	}
}
