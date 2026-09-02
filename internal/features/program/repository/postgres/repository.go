package program_repository_postgres

import (
	core_logger "workout_app/internal/core/logger"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
)

type ProgramRepository struct {
	pool *core_pool_pgx.Pool
	log  *core_logger.Logger
}

func NewProgramRepository(
	pool *core_pool_pgx.Pool,
	log *core_logger.Logger,
) *ProgramRepository {
	return &ProgramRepository{
		pool: pool,
		log:  log,
	}
}
