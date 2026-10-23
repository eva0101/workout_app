package schedule_repository_postgres

import (
	core_logger "workout_app/internal/core/logger"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
)

type ScheduleRepository struct {
	pool *core_pool_pgx.Pool
	log  *core_logger.Logger
}

func NewScheduleRepository(
	pool *core_pool_pgx.Pool,
	log *core_logger.Logger,
) *ScheduleRepository {
	return &ScheduleRepository{
		pool: pool,
		log:  log,
	}
}
